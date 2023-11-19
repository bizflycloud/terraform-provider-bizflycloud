// This file is part of terraform-provider-bizflycloud
//
// Copyright (C) 2021  Bizfly Cloud
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>

package bizflycloud

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bizflycloud/gobizfly"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	attachTypeDataDisk = "datadisk"
	attachTypeRootDisk = "rootdisk"
)

func resourceBizflyCloudServer() *schema.Resource {
	return &schema.Resource{
		Create:        resourceBizflyCloudServerCreate,
		Read:          resourceBizflyCloudServerRead,
		Update:        resourceBizflyCloudServerUpdate,
		Delete:        resourceBizflyCloudServerDelete,
		SchemaVersion: 1,
		Schema:        resourceServerSchema(),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
		},
	}
}

func resourceBizflyCloudServerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()

	// Build up creation options
	rootDiskPayload := gobizfly.ServerDisk{Size: d.Get("root_disk_size").(int)}
	rootDiskVolumeType := d.Get("root_disk_volume_type").(string)
	rootDiskPayload.VolumeType = &rootDiskVolumeType
	vpcNetworkIDs := readStringArray(d.Get("vpc_network_ids").(*schema.Set).List())
	scr := &gobizfly.ServerCreateRequest{
		Name:       d.Get("name").(string),
		FlavorName: d.Get("flavor_name").(string),
		SSHKey:     d.Get("ssh_key").(string),
		Type:       d.Get("category").(string),
		OS: &gobizfly.ServerOS{
			Type: d.Get("os_type").(string),
			ID:   d.Get("os_id").(string),
		},
		AvailabilityZone: d.Get("availability_zone").(string),
		Password:         d.Get("password").(bool),
		RootDisk:         &rootDiskPayload,
		NetworkPlan:      d.Get("network_plan").(string),
		BillingPlan:      d.Get("billing_plan").(string),
		UserData:         d.Get("user_data").(string),
	}
	var (
		isCreatedWan         bool
		usingV6Wan           bool
		freeWanV4BlockCount  int
		freeWanV6BlockCount  int
		networkInterfaceIDs  []string
		freeWanV4FirewallIDs []string
		freeWanV6FirewallIDs []string
	)
	attachPortIPFirewallIDs := make(map[string][]string)
	for _, v := range d.Get("free_wan").(*schema.Set).List() {
		freeWan := v.(map[string]interface{})
		if freeWan["ip_version"].(int) == 4 {
			freeWanV4BlockCount++
			if freeWanV4BlockCount > 1 {
				return fmt.Errorf("only one free WAN V4 is allowed")
			}
			isCreatedWan = true
			freeWanV4FirewallIDs = readStringArray(freeWan["firewall_ids"].(*schema.Set).List())
		} else {
			freeWanV6BlockCount++
			usingV6Wan = true
			if freeWanV6BlockCount > 1 {
				return fmt.Errorf("only one free WAN V6 is allowed")
			}
			freeWanV6FirewallIDs = readStringArray(freeWan["firewall_ids"].(*schema.Set).List())
		}
	}
	for _, v := range d.Get("network_interface").(*schema.Set).List() {
		networkInterface := v.(map[string]interface{})
		portIDStr := networkInterface["id"].(string)
		firewallIDs := readStringArray(networkInterface["firewall_ids"].(*schema.Set).List())
		attachPortIPFirewallIDs[portIDStr] = firewallIDs
		networkInterfaceIDs = append(networkInterfaceIDs, portIDStr)
	}
	scr.IsCreatedWan = &isCreatedWan
	scr.IPv6 = usingV6Wan
	scr.VPCNetworkIds = vpcNetworkIDs
	scr.NetworkInterfaces = networkInterfaceIDs
	log.Printf("[DEBUG] Create Cloud Server configuration: %#v", scr)

	tasks, err := client.Server.Create(context.Background(), scr)
	if err != nil {
		return fmt.Errorf("Error creating server: %s", err)
	}
	// Set ID of server with task ID, we need to change to the real ID after server is created
	d.SetId(tasks.Task[0])
	log.Printf("[INFO] Server is creating with task ID: %s", d.Id())
	// wait for cloud server to become active
	_, err = waitForServerCreate(d, meta)
	if err != nil {
		return fmt.Errorf("Error creating cloud server with task id (%s): %s", d.Id(), err)
	}

	ports, err := client.NetworkInterface.List(context.Background(), &gobizfly.ListNetworkInterfaceOptions{
		Type:   "LAN_WAN",
		Status: "ACTIVE",
	})
	var wg sync.WaitGroup
	errChan := make(chan error, len(ports))
	for _, port := range ports {
		if port.DeviceID != d.Id() {
			continue
		}
		if port.BillingType == "free" && port.Type == "WAN" {
			if port.IPVersion == 6 {
				wg.Add(1)
				go func(portID string) {
					defer wg.Done()
					if err := attachFirewallsForPort(client, portID, freeWanV6FirewallIDs); err != nil {
						errChan <- fmt.Errorf("error attaching firewall for port %s: %v", port.ID, err)
					}
				}(port.ID)
			} else {
				wg.Add(1)
				go func(portID string) {
					defer wg.Done()
					if err := attachFirewallsForPort(client, portID, freeWanV4FirewallIDs); err != nil {
						errChan <- fmt.Errorf("error attaching firewall for port %s: %v", port.ID, err)
					}
				}(port.ID)
			}
		} else if checkIDInList(port.ID, networkInterfaceIDs) {
			wg.Add(1)
			go func(portID string) {
				defer wg.Done()
				if err := attachFirewallsForPort(client, portID, attachPortIPFirewallIDs[port.ID]); err != nil {
					errChan <- fmt.Errorf("error attaching firewall for port %s: %v", port.ID, err)
				}
			}(port.ID)
		}
	}
	wg.Wait()
	if len(errChan) > 0 {
		return <-errChan
	}
	// attach volume in volume_ids list after server is created
	if attr, ok := d.GetOk("volume_ids"); ok {
		var volumes []string
		for _, id := range attr.(*schema.Set).List() {
			if id == nil {
				continue
			}
			volumeId := id.(string)
			if volumeId == "" {
				continue
			}
			volumes = append(volumes, volumeId)
		}
		err = attachVolumes(d.Id(), volumes, client)
		if err != nil {
			return err
		}
	}
	return resourceBizflyCloudServerRead(d, meta)
}

func resourceBizflyCloudServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	server, err := client.Server.Get(context.Background(), d.Id())
	if err != nil {
		if errors.Is(err, gobizfly.ErrNotFound) {
			log.Printf("[WARN] Bizfly Cloud Server (%s) is not found", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving server: %v", err)
	}
	networkInterfaces, _ := client.NetworkInterface.List(context.Background(), &gobizfly.ListNetworkInterfaceOptions{
		Type: "LAN_WAN",
	})
	serverFreeWan := make([]map[string]interface{}, 0)
	serverNetworkInterfaces := make([]map[string]interface{}, 0)
	vpcNetworkIDs := make([]string, 0)
	for _, networkInterface := range networkInterfaces {
		if networkInterface.DeviceID != d.Id() {
			continue
		}
		serverNetworkInterface := make(map[string]interface{})
		serverNetworkInterface["id"] = networkInterface.ID
		serverNetworkInterface["firewall_ids"] = networkInterface.SecurityGroups
		if networkInterface.Type == "WAN" && networkInterface.BillingType == "free" {
			if networkInterface.IPVersion == 6 {
				serverNetworkInterface["ip_version"] = 6
			} else {
				serverNetworkInterface["ip_version"] = 4
			}
			serverFreeWan = append(serverFreeWan, serverNetworkInterface)
		} else {
			serverNetworkInterfaces = append(serverNetworkInterfaces, serverNetworkInterface)
			if networkInterface.Type == "LAN" {
				vpcNetworkIDs = append(vpcNetworkIDs, networkInterface.NetworkID)
			}
		}
	}
	_ = d.Set("name", server.Name)
	_ = d.Set("key_name", server.KeyName)
	_ = d.Set("status", server.Status)
	_ = d.Set("flavor_name", formatFlavor(server.Flavor.Name))
	_ = d.Set("category", server.Category)
	_ = d.Set("user_id", server.UserID)
	_ = d.Set("project_id", server.ProjectID)
	_ = d.Set("availability_zone", server.AvailabilityZone)
	_ = d.Set("created_at", server.CreatedAt)
	_ = d.Set("updated_at", server.UpdatedAt)
	_ = d.Set("billing_plan", server.BillingPlan)
	_ = d.Set("zone_name", server.ZoneName)
	_ = d.Set("is_available", server.IsAvailable)
	_ = d.Set("locked", server.Locked)
	_ = d.Set("network_plan", server.NetworkPlan)
	if d.Get("network_interface") != nil {
		_ = d.Set("network_interface", serverNetworkInterfaces)
	} else {
		_ = d.Set("vpc_network_ids", vpcNetworkIDs)
	}
	_ = d.Set("free_wan", serverFreeWan)

	if err = d.Set("volume_ids", flatternBizflyCloudVolumeIDs(server.AttachedVolumes)); err != nil {
		return fmt.Errorf("Error setting `volume_ids`: %+v", err)
	}
	return nil
}

func resourceBizflyCloudServerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	id := d.Id()
	if d.HasChange("flavor_name") {
		// Resize server to new flavor
		task, err := client.Server.Resize(context.Background(), id, d.Get("flavor_name").(string))
		if err != nil {
			return fmt.Errorf("Error when resize server [%s]: %v", id, err)
		}
		// wait for server is active again
		_, err = waitForServerUpdate(d, meta, task.TaskID)
		if err != nil {
			return fmt.Errorf("Error updating cloud server with task id (%s): %s", d.Id(), err)
		}
	}
	if d.HasChange("vpc_network_ids") {
		oldVPCIds, newVPCIds := d.GetChange("vpc_network_ids")
		oldIDSet := newSet(oldVPCIds.(*schema.Set).List())
		newIDSet := newSet(newVPCIds.(*schema.Set).List())
		var (
			attachVPCs []string
			detachVPCs []string
		)
		for vpcId := range leftDiff(oldIDSet, newIDSet) {
			detachVPCs = append(detachVPCs, vpcId)
		}
		for vpcId := range leftDiff(newIDSet, oldIDSet) {
			attachVPCs = append(attachVPCs, vpcId)
		}
		log.Printf("[DEBUG] attachVPCs: %#v", attachVPCs)
		log.Printf("[DEBUG] detachVPCs: %#v", detachVPCs)
		if len(detachVPCs) > 0 {
			_, err := client.Server.RemoveVPC(context.Background(), d.Id(), detachVPCs)
			if err != nil {
				return fmt.Errorf("Error removing VPCs from server [%s]: %v", d.Id(), err)
			}
		}
		if len(attachVPCs) > 0 {
			_, err := client.Server.AddVPC(context.Background(), d.Id(), attachVPCs)
			if err != nil {
				return fmt.Errorf("Error adding VPCs to server [%s]: %v", d.Id(), err)
			}
		}
	}

	if d.HasChange("category") {
		// Change category of the server
		task, err := client.Server.ChangeCategory(context.Background(), id, d.Get("category").(string))
		if err != nil {
			return fmt.Errorf("Error when change category of server [%s]: %v", id, err)
		}
		// wait for server is active again
		_, err = waitForServerUpdate(d, meta, task.TaskID)
		if err != nil {
			return fmt.Errorf("Error updating cloud server with task id (%s): %s", d.Id(), err)
		}
	}

	// if volume_ids is changed, update the attached volumes
	if d.HasChange("volume_ids") {
		oldIDs, newIDs := d.GetChange("volume_ids")

		oldIDSet := newSet(oldIDs.(*schema.Set).List())
		newIDSet := newSet(newIDs.(*schema.Set).List())
		for volumeID := range leftDiff(newIDSet, oldIDSet) {
			_, err := client.Volume.Attach(context.Background(), volumeID, id)
			if err != nil {
				return fmt.Errorf("Error attaching volume %q to server (%s): %v", volumeID, id, err)
			}
		}
		for volumeID := range leftDiff(oldIDSet, newIDSet) {
			_, err := client.Volume.Detach(context.Background(), volumeID, id)
			if err != nil {
				return fmt.Errorf("Error detaching volume %q from server (%s): %v", volumeID, id, err)
			}
		}
	}
	if d.HasChanges("network_plan") {
		_, newNetworkPlan := d.GetChange("network_plan")
		err := client.Server.ChangeNetworkPlan(context.Background(), d.Id(), newNetworkPlan.(string))
		if err != nil {
			return fmt.Errorf("error changing network plan of server [%s]: %v", d.Id(), err)
		}
	}
	if d.HasChange("billing_plan") {
		_, newBillingPlan := d.GetChange("billing_plan")
		err := client.Server.SwitchBillingPlan(context.Background(), d.Id(), newBillingPlan.(string))
		if err != nil {
			return fmt.Errorf("error changing billing plan of server [%s]: %v", d.Id(), err)
		}
	}
	if d.HasChange("free_wan") {
		oldFreeWan, newFreeWan := d.GetChange("free_wan")
		var (
			freeWanV4ID, freeWanV6ID string
		)
		oldFreeWanList := oldFreeWan.(*schema.Set).List()
		newFreeWanList := newFreeWan.(*schema.Set).List()
		oldFreeWanMap := make(map[int]bool)
		newFreeWanMap := make(map[int]bool)
		oldFreeWanFirewallIDs := make(map[int][]string)
		newFreeWanFirewallIDs := make(map[int][]string)
		removeFreeWan := make([]string, 0)
		addFirewallIDs := make(map[string][]string)
		log.Printf("[DEBUG] oldFreeWanMap: %#v", oldFreeWanMap)
		log.Printf("[DEBUG] newFreeWanMap: %#v", newFreeWanMap)

		removeFirewallIDs := make(map[string][]string)
		for _, v := range oldFreeWanList {
			freeWan := v.(map[string]interface{})
			portID := freeWan["id"].(string)
			ipVersion := freeWan["ip_version"].(int)
			oldFreeWanMap[ipVersion] = true
			firewallIDs := readStringArray(freeWan["firewall_ids"].(*schema.Set).List())
			oldFreeWanFirewallIDs[ipVersion] = firewallIDs
			if ipVersion == 4 {
				freeWanV4ID = portID
			} else {
				freeWanV6ID = portID
			}
		}
		log.Printf("[DEBUG] oldFreeWanFirewallIDs: %#v", oldFreeWanFirewallIDs)
		var portID string
		for _, v := range newFreeWanList {
			freeWan := v.(map[string]interface{})
			ipVersion := freeWan["ip_version"].(int)
			if ipVersion == 4 {
				portID = freeWanV4ID
			} else {
				portID = freeWanV6ID
			}

			firewallIDs := readStringArray(freeWan["firewall_ids"].(*schema.Set).List())
			newFreeWanFirewallIDs[ipVersion] = firewallIDs
			newFreeWanMap[ipVersion] = true
			if !oldFreeWanMap[ipVersion] {
				// add port
				return errors.New("cannot add free wan port after creating server")
			} else {
				// check diff in firewall
				for _, firewallID := range firewallIDs {
					if !checkIDInList(firewallID, oldFreeWanFirewallIDs[ipVersion]) {
						addFirewallIDs[portID] = append(addFirewallIDs[portID], firewallID)
					}
				}
			}
		}
		log.Printf("[DEBUG] newFreeWanFirewallIDs: %#v", newFreeWanFirewallIDs)
		for _, v := range oldFreeWanList {
			freeWan := v.(map[string]interface{})
			portID = freeWan["id"].(string)
			ipVersion := freeWan["ip_version"].(int)
			if !newFreeWanMap[ipVersion] {
				// remove port
				removeFreeWan = append(removeFreeWan, portID)
			} else {
				// check diff in firewall
				for _, firewallID := range oldFreeWanFirewallIDs[ipVersion] {
					if !checkIDInList(firewallID, newFreeWanFirewallIDs[ipVersion]) {
						removeFirewallIDs[portID] = append(removeFirewallIDs[portID], firewallID)
					}
				}
			}
		}
		log.Printf("[DEBUG] removeFreeWan: %#v", removeFreeWan)
		log.Printf("[DEBUG] addFirewallIDs: %#v", addFirewallIDs)
		log.Printf("[DEBUG] removeFirewallIDs: %#v", removeFirewallIDs)
		var wg sync.WaitGroup
		errChan := make(chan error, len(removeFreeWan)+len(addFirewallIDs)+len(removeFirewallIDs))
		for _, portID = range removeFreeWan {
			wg.Add(1)
			go func(portID string) {
				defer wg.Done()
				if err := detachServerForPort(client, portID); err != nil {
					errChan <- fmt.Errorf("error detaching server for port %s: %v", portID, err)
				}
			}(portID)
		}
		for portID, firewallIDs := range addFirewallIDs {
			wg.Add(1)
			go func(portID string, firewallIDs []string) {
				defer wg.Done()
				if err := attachFirewallsForPort(client, portID, firewallIDs); err != nil {
					errChan <- fmt.Errorf("error attaching firewall for port %s: %v", portID, err)
				}
			}(portID, firewallIDs)
		}
		for portID, firewallIDs := range removeFirewallIDs {
			wg.Add(1)
			go func(portID string, firewallIDs []string) {
				defer wg.Done()
				if err := detachFirewallsForPort(client, portID, firewallIDs); err != nil {
					errChan <- fmt.Errorf("error detaching firewall for port %s: %v", portID, err)
				}
			}(portID, firewallIDs)
		}
		wg.Wait()
		if len(errChan) > 0 {
			return <-errChan
		}
	}
	if d.HasChange("network_interface") {
		oldInterface, newInterface := d.GetChange("network_interface")
		oldNetworkInterfaces := readNetworkInterface(oldInterface)
		newNetworkInterfaces := readNetworkInterface(newInterface)
		log.Printf("[DEBUG] oldNetworkInterfaces: %#v", oldNetworkInterfaces)
		log.Printf("[DEBUG] newNetworkInterfaces: %#v", newNetworkInterfaces)
		isOldNetInterface := make(map[string]bool)
		isNewNetInterface := make(map[string]bool)
		addFirewallsMap := make(map[string][]string)
		removeFirewallsMap := make(map[string][]string)
		var (
			addNetworkInterfaces    []ServerNetworkInterface
			removeNetworkInterfaces []ServerNetworkInterface
		)
		for _, networkInterface := range oldNetworkInterfaces {
			isOldNetInterface[networkInterface.ID] = true
		}
		for _, networkInterface := range newNetworkInterfaces {
			isNewNetInterface[networkInterface.ID] = true
			if !isOldNetInterface[networkInterface.ID] {
				addNetworkInterfaces = append(addNetworkInterfaces, networkInterface)
			} else {
				for _, oldNetworkInterface := range oldNetworkInterfaces {
					if oldNetworkInterface.ID == networkInterface.ID {
						if !cmp.Equal(oldNetworkInterface.FirewallIDs, networkInterface.FirewallIDs) {
							// find the add and remove firewalls
							for _, firewallID := range networkInterface.FirewallIDs {
								if !checkIDInList(firewallID, oldNetworkInterface.FirewallIDs) {
									addFirewallsMap[networkInterface.ID] = append(addFirewallsMap[networkInterface.ID], firewallID)
								} else {
									removeFirewallsMap[networkInterface.ID] = append(removeFirewallsMap[networkInterface.ID], firewallID)
								}
							}
						}
					}
				}
			}
		}
		log.Printf("[DEBUG] isOldNetInterface: %#v", isOldNetInterface)
		log.Printf("[DEBUG] isNewNetInterface: %#v", isNewNetInterface)
		for _, networkInterface := range oldNetworkInterfaces {
			if !isNewNetInterface[networkInterface.ID] {
				removeNetworkInterfaces = append(removeNetworkInterfaces, networkInterface)
			}
		}
		log.Printf("[DEBUG] addNetworkInterfaces: %#v", addNetworkInterfaces)
		log.Printf("[DEBUG] removeNetworkInterfaces: %#v", removeNetworkInterfaces)
		log.Printf("[DEBUG] addFirewallsMap: %#v", addFirewallsMap)
		log.Printf("[DEBUG] removeFirewallsMap: %#v", removeFirewallsMap)
		var wg sync.WaitGroup
		errChan := make(chan error, len(addNetworkInterfaces)+len(removeNetworkInterfaces)+
			len(addFirewallsMap)+len(removeFirewallsMap))
		for _, addNetworkInterface := range addNetworkInterfaces {
			wg.Add(1)
			go func(networkInterface *ServerNetworkInterface) {
				defer wg.Done()
				if err := attachServerForPort(client, d.Id(), networkInterface.ID); err != nil {
					errChan <- fmt.Errorf("error attaching server for port %s: %v", networkInterface.ID, err)
				}
			}(&addNetworkInterface)
		}
		for _, removeNetworkInterface := range removeNetworkInterfaces {
			wg.Add(1)
			go func(networkInterface *ServerNetworkInterface) {
				defer wg.Done()
				if err := detachServerForPort(client, networkInterface.ID); err != nil {
					errChan <- fmt.Errorf("error detaching server for port %s: %v", networkInterface.ID, err)
				}
			}(&removeNetworkInterface)
		}
		for portID, firewallIDs := range addFirewallsMap {
			wg.Add(1)
			go func(portID string, firewallIDs []string) {
				defer wg.Done()
				if err := attachFirewallsForPort(client, portID, firewallIDs); err != nil {
					errChan <- fmt.Errorf("error attaching firewall for port %s: %v", portID, err)
				}
			}(portID, firewallIDs)
		}
		for portID, firewallIDs := range removeFirewallsMap {
			wg.Add(1)
			go func(portID string, firewallIDs []string) {
				defer wg.Done()
				if err := detachFirewallsForPort(client, portID, firewallIDs); err != nil {
					errChan <- fmt.Errorf("error detaching firewall for port %s: %v", portID, err)
				}
			}(portID, firewallIDs)
		}
		wg.Wait()
		if len(errChan) > 0 {
			return <-errChan
		}
	}
	return resourceBizflyCloudServerRead(d, meta)
}

func resourceBizflyCloudServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	// destroy the cloud server
	server, err := client.Server.Get(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving cloud server: %v", err)
	}
	var rootDiskID string
	for _, v := range server.AttachedVolumes {
		if v.AttachedType == attachTypeRootDisk {
			rootDiskID = v.ID
		}
	}
	task, err := client.Server.Delete(context.Background(), d.Id(), []string{rootDiskID})
	if err != nil {
		return fmt.Errorf("Error delete cloud server %v", err)
	}

	_, err = waitforServerDelete(d, meta, task.TaskID)
	if err != nil && !errors.Is(err, gobizfly.ErrNotFound) {
		return fmt.Errorf("Error delete cloud server with task id (%s): %s", d.Id(), err)
	}
	return nil
}

func waitForServerCreate(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	log.Printf("[INFO] Waiting for server with task id (%s) to be created", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending:        []string{"BUILD"},
		Target:         []string{"ACTIVE"},
		Refresh:        newServerStateRefreshfunc(d, "status", meta),
		Timeout:        1200 * time.Second,
		Delay:          20 * time.Second,
		MinTimeout:     3 * time.Second,
		NotFoundChecks: 100,
	}
	return stateConf.WaitForState()
}

func waitforServerDelete(d *schema.ResourceData, meta interface{}, taskID string) (interface{}, error) {
	log.Printf("[INFO] Waiting for server with task id (%s) to be deleted", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"false"},
		Target:     []string{"true"},
		Refresh:    waitToDeleteServerRefreshFunc(d, meta, taskID),
		Timeout:    600 * time.Second,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	return stateConf.WaitForState()
}

func waitForServerUpdate(d *schema.ResourceData, meta interface{}, taskID string) (interface{}, error) {
	log.Printf("[INFO] Waiting for server with task id (%s) to be created", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"HARD_REBOOT", "MIGRATING", "REBUILD", "RESIZE"},
		Target:     []string{"ACTIVE"},
		Refresh:    updateServerStateRefreshfunc(d, "status", meta, taskID),
		Timeout:    600 * time.Second,
		Delay:      20 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	return stateConf.WaitForState()
}

func newServerStateRefreshfunc(d *schema.ResourceData, attribute string, meta interface{}) resource.StateRefreshFunc {
	client := meta.(*CombinedConfig).gobizflyClient()
	return func() (interface{}, string, error) {
		// Get task result from cloud server API
		resp, err := client.Server.GetTask(context.Background(), d.Id())
		if err != nil {
			return nil, "", err
		}
		// if the task is not ready, we need to wait for a moment
		if !resp.Ready {
			log.Println("[DEBUG] Cloud Server is not ready")
			return nil, "", nil
		}
		// server is ready now, set ID for resourceData
		d.SetId(resp.Result.Server.ID)
		err = resourceBizflyCloudServerRead(d, meta)
		if err != nil {
			return nil, "", err
		}

		if attr, ok := d.GetOk(attribute); ok { // nolint
			server, err := client.Server.Get(context.Background(), resp.Result.Server.ID)
			if err != nil {
				return nil, "", fmt.Errorf("Error retrieving cloud server: %v", err)
			}
			switch attr := attr.(type) {
			case bool:
				return &server, strconv.FormatBool(attr), nil
			default:
				return &server, attr.(string), nil
			}
		}
		return nil, "", nil
	}
}

func updateServerStateRefreshfunc(d *schema.ResourceData, attribute string, meta interface{}, taskID string) resource.StateRefreshFunc {
	client := meta.(*CombinedConfig).gobizflyClient()
	return func() (interface{}, string, error) {
		// Get task result from cloud server API
		resp, err := client.Server.GetTask(context.Background(), taskID)
		if err != nil {
			return nil, "", err
		}
		// if the task is not ready, we need to wait for a moment
		if !resp.Ready {
			log.Println("[DEBUG] Cloud Server is not ready")
			return nil, "", nil
		}
		err = resourceBizflyCloudServerRead(d, meta)
		if err != nil {
			return nil, "", err
		}
		if attr, ok := d.GetOk(attribute); ok { // nolint
			server, err := client.Server.Get(context.Background(), d.Id())
			if err != nil {
				return nil, "", fmt.Errorf("Error retrieving cloud server: %v", err)
			}
			switch attr := attr.(type) {
			case bool:
				return &server, strconv.FormatBool(attr), nil
			default:
				return &server, attr.(string), nil
			}
		}
		return nil, "", nil
	}
}

func waitToDeleteServerRefreshFunc(d *schema.ResourceData, meta interface{}, taskID string) resource.StateRefreshFunc {
	client := meta.(*CombinedConfig).gobizflyClient()
	return func() (interface{}, string, error) {

		resp, err := client.Server.GetTask(context.Background(), taskID)
		if err != nil {
			return nil, "false", err
		}
		server, err := client.Server.Get(context.Background(), d.Id())
		if errors.Is(err, gobizfly.ErrNotFound) {
			return server, "true", nil
		} else if err != nil {
			return nil, "false", err
		}
		return server, strconv.FormatBool(resp.Ready), nil
	}
}

func formatFlavor(s string) string {
	// This function will be removed in the near future when the API format for us
	if strings.Contains(s, ".") {
		return strings.Split(s, ".")[1]
	}
	return strings.Join(strings.Split(s, "_")[:2], "_")
}

func flatternBizflyCloudVolumeIDs(volumeids []gobizfly.AttachedVolume) *schema.Set {
	flattenedVolumes := schema.NewSet(schema.HashString, []interface{}{})
	for _, v := range volumeids {
		if v.AttachedType == attachTypeDataDisk {
			flattenedVolumes.Add(v.ID)
		}
	}
	return flattenedVolumes
}

func attachVolumes(id string, volumeids []string, client *gobizfly.Client) error {
	for _, vid := range volumeids {
		_, err := client.Volume.Attach(context.Background(), id, vid)
		if err != nil {
			return err
		}
	}
	return nil
}

func flatternBizflyCloudIPs(ips []gobizfly.IP) *schema.Set {
	flatternIPs := schema.NewSet(schema.HashString, []interface{}{})
	for _, ip := range ips {
		flatternIPs.Add(ip.Address)
	}
	return flatternIPs
}

func readStringArray(items []interface{}) []string {
	stringArray := make([]string, 0)
	for i := 0; i < len(items); i++ {
		networkInterface := items[i].(string)
		stringArray = append(stringArray, networkInterface)
	}
	return stringArray
}

func newSet(ids []interface{}) map[string]struct{} {
	out := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		out[id.(string)] = struct{}{}
	}
	return out
}

func leftDiff(left, right map[string]struct{}) map[string]struct{} {
	out := make(map[string]struct{})
	for l := range left {
		if _, ok := right[l]; !ok {
			out[l] = struct{}{}
		}
	}
	return out
}

func uniqueList(list []string) []string {
	unique := make(map[string]struct{})
	for _, v := range list {
		unique[v] = struct{}{}
	}
	out := make([]string, 0, len(unique))
	for v := range unique {
		out = append(out, v)
	}
	return out
}

func attachFirewallsForPort(client *gobizfly.Client, portID string, firewallIDs []string) error {
	if firewallIDs == nil || len(firewallIDs) == 0 {
		return nil
	}
	_, err := client.NetworkInterface.Action(context.Background(), portID,
		&gobizfly.ActionNetworkInterfacePayload{
			Action:         "add_firewall",
			SecurityGroups: firewallIDs,
		})
	return err

}

func detachFirewallsForPort(client *gobizfly.Client, portID string, firewallIDs []string) error {
	if firewallIDs == nil || len(firewallIDs) == 0 {
		return nil
	}
	_, err := client.NetworkInterface.Action(context.Background(), portID,
		&gobizfly.ActionNetworkInterfacePayload{
			Action:         "remove_firewall",
			SecurityGroups: firewallIDs,
		})
	return err
}

func attachServerForPort(client *gobizfly.Client, serverID, portID string) error {
	_, err := client.NetworkInterface.Action(context.Background(), portID,
		&gobizfly.ActionNetworkInterfacePayload{
			Action:   "attach_server",
			ServerID: serverID,
		})
	return err
}

func detachServerForPort(client *gobizfly.Client, portID string) error {
	_, err := client.NetworkInterface.Action(context.Background(), portID,
		&gobizfly.ActionNetworkInterfacePayload{
			Action: "detach_server",
		})
	return err
}

func checkIDInList(id string, IDs []string) bool {
	for _, v := range IDs {
		if v == id {
			return true
		}
	}
	return false
}

func readNetworkInterface(inputInterface interface{}) (networkInterfaces []ServerNetworkInterface) {
	for _, v := range inputInterface.(*schema.Set).List() {
		networkInterface := v.(map[string]interface{})
		var portIDStr string
		if portID := networkInterface["id"]; portID != nil {
			portIDStr = networkInterface["id"].(string)
		}
		firewallIDs := readStringArray(networkInterface["firewall_ids"].(*schema.Set).List())
		networkInterfaces = append(networkInterfaces, ServerNetworkInterface{
			ID:          portIDStr,
			FirewallIDs: firewallIDs,
		})
	}
	return
}
