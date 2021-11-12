// This file is part of gobizfly
//
// Copyright (C) 2020  BizFly Cloud
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
	"time"

	"github.com/bizflycloud/gobizfly"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	attachTypeDataDisk = "datadisk"
	attachTypeRootDisk = "rootdisk"
)

func resourceBizFlyCloudServer() *schema.Resource {
	return &schema.Resource{
		Create:        resourceBizFlyCloudServerCreate,
		Read:          resourceBizFlyCloudServerRead,
		Update:        resourceBizFlyCloudServerUpdate,
		Delete:        resourceBizFlyCloudServerDelete,
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"flavor_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ssh_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"category": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"os_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"os_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"root_disk_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"root_disk_size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"volume_ids": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
			},
			"lan_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"wan_ipv4": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"wan_ipv6": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
		},
	}
}

func resourceBizFlyCloudServerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()

	// Build up creation options
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
		RootDisk: &gobizfly.ServerDisk{
			Type: d.Get("root_disk_type").(string),
			Size: d.Get("root_disk_size").(int),
		},
	}
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
	return resourceBizFlyCloudServerRead(d, meta)
}

func resourceBizFlyCloudServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	server, err := client.Server.Get(context.Background(), d.Id())
	if err != nil {
		if errors.Is(err, gobizfly.ErrNotFound) {
			log.Printf("[WARN] BizFly Cloud Server (%s) is not found", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving server: %v", err)
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

	if err := d.Set("volume_ids", flatternBizFlyCloudVolumeIDs(server.AttachedVolumes)); err != nil {
		return fmt.Errorf("Error setting `volume_ids`: %+v", err)
	}

	_ = d.Set("lan_ip", server.IPAddresses.LanAddresses[0].Address)
	if err := d.Set("wan_ipv4", flatternBizFlyCloudIPs(server.IPAddresses.WanV4Addresses)); err != nil {
		return fmt.Errorf("Error setting `wan_ipv4`: #{err}")
	}
	if err := d.Set("wan_ipv6", flatternBizFlyCloudIPs(server.IPAddresses.WanV6Addresses)); err != nil {
		return fmt.Errorf("Error setting `wan_ipv6`: ${err}")
	}
	return nil
}

func resourceBizFlyCloudServerUpdate(d *schema.ResourceData, meta interface{}) error {
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
		newSet := func(ids []interface{}) map[string]struct{} {
			out := make(map[string]struct{}, len(ids))
			for _, id := range ids {
				out[id.(string)] = struct{}{}
			}
			return out
		}

		// leftDiff returns all elements in Left that are not in Right
		leftDiff := func(left, right map[string]struct{}) map[string]struct{} {
			out := make(map[string]struct{})
			for l := range left {
				if _, ok := right[l]; !ok {
					out[l] = struct{}{}
				}
			}
			return out
		}

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
	return resourceBizFlyCloudServerRead(d, meta)
}

func resourceBizFlyCloudServerDelete(d *schema.ResourceData, meta interface{}) error {
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
	err = client.Server.Delete(context.Background(), d.Id(), []string{rootDiskID})
	if err != nil {
		return fmt.Errorf("Error delete cloud server %v", err)
	}
	// TODO check server is deleted
	// remove rootdisk of the server
	return nil
}

func waitForServerCreate(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	log.Printf("[INFO] Waiting for server with task id (%s) to be created", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"BUILD"},
		Target:     []string{"ACTIVE"},
		Refresh:    newServerStateRefreshfunc(d, "status", meta),
		Timeout:    600 * time.Second,
		Delay:      20 * time.Second,
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
		err = resourceBizFlyCloudServerRead(d, meta)
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
		err = resourceBizFlyCloudServerRead(d, meta)
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
func formatFlavor(s string) string {
	// This function will be removed in the near future when the API format for us
	if strings.Contains(s, ".") {
		return strings.Split(s, ".")[1]
	}
	return strings.Join(strings.Split(s, "_")[:2], "_")
}

func flatternBizFlyCloudVolumeIDs(volumeids []gobizfly.AttachedVolume) *schema.Set {
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

func flatternBizFlyCloudIPs(ips []gobizfly.IP) *schema.Set {
	flatternIPs := schema.NewSet(schema.HashString, []interface{}{})
	for _, ip := range ips {
		flatternIPs.Add(ip.Address)
	}
	return flatternIPs
}
