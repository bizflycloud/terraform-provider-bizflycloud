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
		RootDisk:         &rootDiskPayload,
		NetworkPlan:      d.Get("network_plan").(string),
		BillingPlan:      d.Get("billing_plan").(string),
		UserData:         d.Get("user_data").(string),
	}
	var (
		isCreatedWan         bool
		usingV6Wan           bool
		freeWanV4FirewallIDs []string
		freeWanV6FirewallIDs []string
	)
	defaultPublicIPv4List := d.Get("default_public_ipv4").([]interface{})
	defaultPublicIPv6List := d.Get("default_public_ipv6").([]interface{})
	if len(defaultPublicIPv4List) > 1 {
		return fmt.Errorf("only one default public ipv4 is allowed")
	}
	if len(defaultPublicIPv6List) > 1 {
		return fmt.Errorf("only one default public ipv6 is allowed")
	}
	for _, v := range defaultPublicIPv4List {
		freeWan := v.(map[string]interface{})
		isCreatedWan = true
		freeWanV4FirewallIDs = readStringArray(freeWan["firewall_ids"].(*schema.Set).List())
	}
	for _, v := range defaultPublicIPv6List {
		freeWan := v.(map[string]interface{})
		usingV6Wan = true
		freeWanV6FirewallIDs = readStringArray(freeWan["firewall_ids"].(*schema.Set).List())
	}
	scr.IsCreatedWan = &isCreatedWan
	scr.IPv6 = usingV6Wan
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
	if err != nil {
		return fmt.Errorf("error listing ports: %v", err)
	}
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
						errChan <- fmt.Errorf("error attaching firewall for port %s: %v", portID, err)
					}
				}(port.ID)
			} else {
				wg.Add(1)
				go func(portID string) {
					defer wg.Done()
					if err := attachFirewallsForPort(client, portID, freeWanV4FirewallIDs); err != nil {
						errChan <- fmt.Errorf("error attaching firewall for port %s: %v", portID, err)
					}
				}(port.ID)
			}
		}
	}
	wg.Wait()
	if len(errChan) > 0 {
		return <-errChan
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
	rootDisk, err := getServerRootDisk(client, server.ID)
	if err != nil {
		log.Printf("[WARN] get rootdisk of server %s failed: %+v", server.ID, err)
		return err
	}
	networkInterfaces, _ := client.NetworkInterface.List(context.Background(), &gobizfly.ListNetworkInterfaceOptions{
		Type: "LAN_WAN",
	})
	firewalls, err := client.Firewall.List(context.Background(), &gobizfly.ListOptions{})
	if err != nil {
		return fmt.Errorf("error listing firewalls: %v", err)
	}
	userFirewallIDs := make([]string, len(firewalls))
	for i, firewall := range firewalls {
		userFirewallIDs[i] = firewall.ID
	}
	vpcNetworkIDs := make([]string, 0)
	networkInterfaceIds := make([]string, 0)
	_ = d.Set("default_public_ipv6", make([]map[string]interface{}, 0))
	_ = d.Set("default_public_ipv4", make([]map[string]interface{}, 0))
	for _, networkInterface := range networkInterfaces {
		if networkInterface.DeviceID != d.Id() {
			continue
		}
		serverNetworkInterface := make(map[string]interface{})
		serverNetworkInterface["id"] = networkInterface.ID
		serverNetworkInterface["firewall_ids"] = networkInterface.SecurityGroups
		serverNetworkInterface["enabled"] = networkInterface.Status == "ACTIVE"
		serverNetworkInterface["ip_address"] = networkInterface.IPAddress
		if networkInterface.Type == "WAN" && networkInterface.BillingType == "free" {
			if networkInterface.IPVersion == 6 {
				_ = d.Set("default_public_ipv6", []map[string]interface{}{serverNetworkInterface})
			} else {
				_ = d.Set("default_public_ipv4", []map[string]interface{}{serverNetworkInterface})
			}
		} else {
			networkInterfaceIds = append(networkInterfaceIds, networkInterface.ID)
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
	_ = d.Set("is_available", server.IsAvailable)
	_ = d.Set("locked", server.Locked)
	_ = d.Set("network_plan", server.NetworkPlan)
	_ = d.Set("ssh_key", server.KeyName)
	_ = d.Set("vpc_network_ids", vpcNetworkIDs)
	_ = d.Set("network_interface_ids", networkInterfaceIds)
	if err = d.Set("volume_ids", flatternBizflyCloudVolumeIDs(server.AttachedVolumes)); err != nil {
		return fmt.Errorf("Error setting `volume_ids`: %+v", err)
	}
	_ = d.Set("os_type", rootDisk.ImageMetadata.ImageType)
	_ = d.Set("os_id", rootDisk.ImageMetadata.ImageID)
	_ = d.Set("root_disk_volume_type", rootDisk.VolumeType)
	_ = d.Set("root_disk_size", rootDisk.Size)
	var state string
	if server.Status == "ACTIVE" {
		state = "running"
	} else if server.Status == "SHUTOFF" {
		state = "stopped"
	} else {
		state = server.Status
	}
	_ = d.Set("state", state)
	return nil
}

func resourceBizflyCloudServerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	id := d.Id()
	if d.HasChange("name") {
		// Rename server
		newName := d.Get("name").(string)
		if err := client.Server.Rename(context.Background(), id, newName); err != nil {
			return fmt.Errorf("Error when rename server [%s]: %v", id, err)
		}
	}
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
	if d.HasChange("category") {
		// Change category of the server
		task, err := client.Server.ChangeCategory(context.Background(), id, d.Get("category").(string))
		if err != nil {
			return fmt.Errorf("error when change category of server [%s]: %v", id, err)
		}
		// wait for server is active again
		_, err = waitForServerUpdate(d, meta, task.TaskID)
		if err != nil {
			return fmt.Errorf("Error updating cloud server with task id (%s): %s", d.Id(), err)
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
	if d.HasChange("default_public_ipv4") {
		if err := updateFreeWantPort(d, client, "default_public_ipv4"); err != nil {
			log.Printf("[ERROR] Error updating free wan port: %v", err)
			return err
		}
	}
	if d.HasChange("default_public_ipv6") {
		if err := updateFreeWantPort(d, client, "default_public_ipv6"); err != nil {
			log.Printf("[ERROR] Error updating free wan port: %v", err)
			return err
		}
	}
	if d.HasChange("state") {
		oldState, newState := d.GetChange("state")
		if oldState.(string) == "ERROR" {
			return fmt.Errorf("cannot change server state because server state is %s", oldState.(string))
		}
		if newState.(string) == "running" {
			_, err := client.Server.Start(context.Background(), d.Id())
			if err != nil {
				return fmt.Errorf("error changing state of server: %v", err)
			}
		} else if newState.(string) == "stopped" {
			_, err := client.Server.Stop(context.Background(), d.Id())
			if err != nil {
				return fmt.Errorf("error changing state of server: %v", err)
			}
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
		Refresh:        newServerStateRefreshFunc(d, "status", meta),
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
		Refresh:    updateServerStateRefreshFunc(d, "status", meta, taskID),
		Timeout:    600 * time.Second,
		Delay:      20 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	return stateConf.WaitForState()
}

func newServerStateRefreshFunc(d *schema.ResourceData, attribute string, meta interface{}) resource.StateRefreshFunc {
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

func updateServerStateRefreshFunc(d *schema.ResourceData, attribute string, meta interface{}, taskID string) resource.StateRefreshFunc {
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
	if len(firewallIDs) == 0 {
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
	if len(firewallIDs) == 0 {
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

func enableIpv6ForServer(client *gobizfly.Client, serverID string) (wanIpv6 *gobizfly.WanIP, err error) {
	enableErr := client.Server.EnableIpv6(context.Background(), serverID)
	if enableErr != nil {
		err = fmt.Errorf("Error enable ipv6 for server %s failed: %v", serverID, enableErr)
		return
	}
	wanIps, wanIpsErr := client.WanIP.List(context.Background())
	if wanIpsErr != nil {
		err = fmt.Errorf("Error list wan ip failed: %v", wanIpsErr)
		return
	}
	for _, wanIp := range wanIps {
		if wanIp.DeviceID != serverID {
			continue
		}
		if wanIp.BillingType == "free" && wanIp.IpVersion == 6 {
			wanIpv6 = wanIp
			break
		}
	}
	if wanIpv6 == nil {
		err = fmt.Errorf("Error enable ipv6 for server %s failed", serverID)
	}
	return
}

func updateFreeWantPort(d *schema.ResourceData, client *gobizfly.Client, field string) error {
	oldPublicIP, newPublicIP := d.GetChange(field)
	newPublicIPList := newPublicIP.([]interface{})
	oldPublicIPList := oldPublicIP.([]interface{})
	if len(newPublicIPList) > 1 {
		return fmt.Errorf("only one %s is allowed", field)
	}
	removeFreeWan := make([]string, 0)
	addFirewallIDs := make([]string, 0)
	removeFirewallIDs := make([]string, 0)
	enablePorts := make([]string, 0)
	disablePorts := make([]string, 0)
	var (
		freeWanID string
	)
	if len(oldPublicIPList) == 1 && len(newPublicIPList) == 0 {
		removeFreeWan = append(removeFreeWan, oldPublicIPList[0].(map[string]interface{})["id"].(string))
	}
	if len(oldPublicIPList) == 0 && len(newPublicIPList) == 1 {
		if field == "default_public_ipv6" {
			serverID := d.Id()
			wanIpv6, err := enableIpv6ForServer(client, serverID)
			if err != nil {
				return err
			}
			freeWanID = wanIpv6.ID
			newFreeWan := newPublicIPList[0].(map[string]interface{})
			if newFreeWan["enabled"].(bool) {
				enablePorts = append(enablePorts, freeWanID)
			} else {
				disablePorts = append(disablePorts, freeWanID)
			}
			newFirewallIDs := readStringArray(newFreeWan["firewall_ids"].(*schema.Set).List())
			oldFirewallIDs := wanIpv6.SecurityGroups
			for _, firewallID := range oldFirewallIDs {
				if !checkIDInList(firewallID, newFirewallIDs) {
					removeFirewallIDs = append(removeFirewallIDs, firewallID)
				}
			}
			for _, firewallID := range newFirewallIDs {
				if !checkIDInList(firewallID, oldFirewallIDs) {
					addFirewallIDs = append(addFirewallIDs, firewallID)
				}
			}
		} else {
			return errors.New("cannot add free wan ipv4 after creating server")
		}
	}
	// update firewall
	if len(oldPublicIPList) == 1 && len(newPublicIPList) == 1 {
		oldFreeWan := oldPublicIPList[0].(map[string]interface{})
		newFreeWan := newPublicIPList[0].(map[string]interface{})
		oldFirewallIDs := readStringArray(oldFreeWan["firewall_ids"].(*schema.Set).List())
		newFirewallIDs := readStringArray(newFreeWan["firewall_ids"].(*schema.Set).List())
		freeWanID = oldFreeWan["id"].(string)
		if oldFreeWan["enabled"].(bool) != newFreeWan["enabled"].(bool) {
			if newFreeWan["enabled"].(bool) {
				enablePorts = append(enablePorts, freeWanID)
			} else {
				disablePorts = append(disablePorts, freeWanID)
			}
		}
		for _, firewallID := range oldFirewallIDs {
			if !checkIDInList(firewallID, newFirewallIDs) {
				removeFirewallIDs = append(removeFirewallIDs, firewallID)
			}
		}
		for _, firewallID := range newFirewallIDs {
			if !checkIDInList(firewallID, oldFirewallIDs) {
				addFirewallIDs = append(addFirewallIDs, firewallID)
			}
		}
	}
	log.Printf("[DEBUG] removeFreeWan %s: %#v", field, removeFreeWan)
	log.Printf("[DEBUG] addFirewallIDs %s: %#v", field, addFirewallIDs)
	log.Printf("[DEBUG] removeFirewallIDs %s: %#v", field, removeFirewallIDs)
	log.Printf("[DEBUG] enablePorts %s: %#v", field, enablePorts)
	log.Printf("[DEBUG] disablePorts %s: %#v", field, disablePorts)
	var wg sync.WaitGroup
	errChan := make(chan error, len(removeFreeWan)+len(addFirewallIDs)+
		len(removeFirewallIDs)+len(enablePorts)+len(disablePorts))
	for _, portID := range removeFreeWan {
		wg.Add(1)
		go func(portID string) {
			defer wg.Done()
			if err := client.NetworkInterface.Delete(context.Background(), portID); err != nil {
				errChan <- fmt.Errorf("error detaching server for port %s: %v", portID, err)
			}
		}(portID)
	}
	for _, firewallID := range addFirewallIDs {
		wg.Add(1)
		go func(firewallID string) {
			defer wg.Done()
			if err := attachFirewallsForPort(client, freeWanID, []string{firewallID}); err != nil {
				errChan <- fmt.Errorf("error attaching firewall for port %s: %v", freeWanID, err)
			}
		}(firewallID)
	}
	for _, firewallID := range removeFirewallIDs {
		wg.Add(1)
		go func(firewallID string) {
			defer wg.Done()
			if err := detachFirewallsForPort(client, freeWanID, []string{firewallID}); err != nil {
				errChan <- fmt.Errorf("error detaching firewall for port %s: %v", freeWanID, err)
			}
		}(firewallID)
	}
	for _, portID := range enablePorts {
		wg.Add(1)
		go func(portID string) {
			defer wg.Done()
			if _, err := client.NetworkInterface.Action(context.Background(), portID,
				&gobizfly.ActionNetworkInterfacePayload{
					Action: "enable",
				}); err != nil {
				errChan <- fmt.Errorf("error enabling port %s: %v", portID, err)
			}
		}(portID)
	}
	for _, portID := range disablePorts {
		wg.Add(1)
		go func(portID string) {
			defer wg.Done()
			if _, err := client.NetworkInterface.Action(context.Background(), portID,
				&gobizfly.ActionNetworkInterfacePayload{
					Action: "disable",
				}); err != nil {
				errChan <- fmt.Errorf("error disabling port %s: %v", portID, err)
			}
		}(portID)
	}
	wg.Wait()
	if len(errChan) > 0 {
		return <-errChan
	}
	return nil
}

func getServerRootDisk(client *gobizfly.Client, serverId string) (*gobizfly.Volume, error) {
	volumes, err := client.Volume.List(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("list volume failed: %+v", err)
	}

	for _, vol := range volumes {
		attachments := vol.Attachments
		if len(attachments) == 0 {
			continue
		}
		attachedServerId := attachments[0].ServerID
		isRootDisk := vol.AttachedType == "rootdisk"
		isServerRootDisk := (attachedServerId == serverId) && isRootDisk
		if isServerRootDisk {
			return vol, nil
		}
	}
	return nil, fmt.Errorf("rootdisk of server %s not found.", serverId)
}
