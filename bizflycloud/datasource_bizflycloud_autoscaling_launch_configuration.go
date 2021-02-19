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
	"fmt"
	"log"

	"github.com/bizflycloud/gobizfly"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceBizFlyCloudLaunchConfiguration() *schema.Resource {
	return &schema.Resource{
		Read:   dataBizFlyCloudAutoScalingLaunchConfigurationRead,
		Schema: dataLaunchConfigurationSchema(),
	}
}

func dataBizFlyCloudAutoScalingLaunchConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()

	if v, ok := d.GetOk("id"); ok {
		d.SetId(v.(string))
	}
	profileID := d.Id()

	log.Printf("[DEBUG] Reading Launch Configuration: %s", profileID)

	profile, err := client.AutoScaling.LaunchConfigurations().Get(context.Background(), profileID)

	log.Printf("[DEBUG] Checking for error: %s", err)
	if err != nil {
		return fmt.Errorf("error describing AutoScaling Groups: %w", err)
	}

	log.Printf("[DEBUG] Found Launch Configuration: %s", profileID)
	log.Printf("[DEBUG] bizflycloud_autoscaling_launch_configuration - Single Launch Configuration found: %s", profile.Name)

	d.SetId(profile.ID)
	_ = d.Set("availability_zone", profile.AvailabilityZone)
	_ = d.Set("flavor", profile.Flavor)
	_ = d.Set("name", profile.Name)
	_ = d.Set("network_plan", profile.NetworkPlan)
	_ = d.Set("instance_type", profile.ProfileType)
	_ = d.Set("key_name", profile.SSHKey)
	_ = d.Set("status", profile.Status)
	_ = d.Set("user_data", profile.UserData)

	if err := d.Set("data_disks", getDataDisks(profile.DataDisks)); err != nil {
		return fmt.Errorf("error setting data_disks: %w", err)
	}

	if err := d.Set("networks", getNetworks(profile.Networks)); err != nil {
		return fmt.Errorf("error setting networks: %w", err)
	}

	if err := d.Set("os", getOperatingSystem(profile.OperatingSystem)); err != nil {
		return fmt.Errorf("error setting os: %w", err)
	}

	if err := d.Set("rootdisk", getRootDisk(profile.RootDisk)); err != nil {
		return fmt.Errorf("error setting rootdisk: %w", err)
	}

	return nil
}

func getDataDisks(m []*gobizfly.AutoScalingDataDisk) []interface{} {
	r := []interface{}{}

	for _, v := range m {
		r = append(r, map[string]interface{}{
			"delete_on_termination": v.DeleteOnTermination,
			"volume_size":           v.Size,
			"volume_type":           v.Type,
		})
	}

	return r
}

func getNetworks(networks []*gobizfly.AutoScalingNetworks) []interface{} {
	r := []interface{}{}

	for _, network := range networks {
		r = append(r, map[string]interface{}{
			"network_id":      network.ID,
			"security_groups": network.SecurityGroups,
		})
	}

	return r
}

func getOperatingSystem(os gobizfly.AutoScalingOperatingSystem) []interface{} {
	r := []interface{}{}

	r = append(r, map[string]interface{}{
		"create_from": os.CreateFrom,
		"error":       os.Error,
		"uuid":        os.ID,
		"os_name":     os.OSName,
	})

	return r
}

func getRootDisk(d *gobizfly.AutoScalingDataDisk) []interface{} {
	r := []interface{}{}

	r = append(r, map[string]interface{}{
		"delete_on_termination": d.DeleteOnTermination,
		"volume_size":           d.Size,
		"volume_type":           d.Type,
	})

	return r
}
