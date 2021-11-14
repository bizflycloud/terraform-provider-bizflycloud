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
	"fmt"
	"log"

	"time"

	"github.com/bizflycloud/gobizfly"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	retryCount = 10
	waitTime   = 1 * time.Minute
)

func resourceBizFlyCloudAutoscalingLaunchConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceBizFlyCloudAutoscalingLaunchConfigurationCreate,
		Read:   resourceBizFlyCloudAutoscalingLaunchConfigurationRead,
		Delete: resourceBizFlyCloudAutoscalingLaunchConfigurationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: resourceLaunchConfigurationSchema(),
	}
}

func resourceBizFlyCloudAutoscalingLaunchConfigurationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()

	lcr := &gobizfly.LaunchConfiguration{
		AvailabilityZone: d.Get("availability_zone").(string),
		Flavor:           d.Get("flavor").(string),
		Name:             d.Get("name").(string),
		NetworkPlan:      d.Get("network_plan").(string),
		ProfileType:      d.Get("instance_type").(string),
		Metadata: map[string]interface{}{
			"category": d.Get("instance_type").(string),
		},
		SSHKey:   d.Get("ssh_key").(string),
		UserData: d.Get("user_data").(string),
	}

	if v, ok := d.GetOk("data_disks"); ok {
		var dataDisks []*gobizfly.AutoScalingDataDisk
		bdms := v.([]interface{})

		for _, bdm := range bdms {
			if bdm == nil {
				continue
			}

			blockDeviceMapping := readBlockDeviceMappingFromConfig(bdm.(map[string]interface{}))
			dataDisks = append(dataDisks, blockDeviceMapping)
		}

		lcr.DataDisks = dataDisks
	}

	if v, ok := d.GetOk("rootdisk"); ok {
		bdms := v.([]interface{})

		for _, bdm := range bdms {
			if bdm == nil {
				continue
			}

			lcr.RootDisk = readBlockDeviceMappingFromConfig(bdm.(map[string]interface{}))
		}

	}

	lcr.OperatingSystem = readOperatingSystemFromConfig(d)

	if v, ok := d.GetOk("networks"); ok {
		var networks []*gobizfly.AutoScalingNetworks
		netList := v.([]interface{})

		for _, ni := range netList {
			if ni == nil {
				continue
			}
			netData := ni.(map[string]interface{})
			networks = append(networks, readNetworksFromConfig(netData))
		}
		lcr.Networks = networks
	}

	profile, err := client.AutoScaling.LaunchConfigurations().Create(context.Background(), lcr)

	if err != nil {
		return fmt.Errorf("[ERROR] Launch Configuration create failed: %v", err)
	}

	d.SetId(profile.ID)
	return resourceBizFlyCloudAutoscalingLaunchConfigurationRead(d, meta)
}

func resourceBizFlyCloudAutoscalingLaunchConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	return dataBizFlyCloudAutoScalingLaunchConfigurationRead(d, meta)
}

func resourceBizFlyCloudAutoscalingLaunchConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	log.Printf("[DEBUG] Launch Configuration destroy: %v", d.Id())

	for retry := retryCount; retry > 0; retry-- {
		if err := client.AutoScaling.LaunchConfigurations().Delete(context.Background(), d.Id()); err != nil {
			log.Printf("[ERROR] Launch Configuration destroy %v was failed: %v", d.Id(), err)
			time.Sleep(waitTime)
			continue
		}

		log.Printf("[DEBUG] Launch Template deleted: %v", d.Id())
		return nil
	}

	return fmt.Errorf("[ERROR] Launch Configuration destroy %v was failed", d.Id())
}

func readBlockDeviceMappingFromConfig(bdm map[string]interface{}) *gobizfly.AutoScalingDataDisk {
	blockDeviceMapping := &gobizfly.AutoScalingDataDisk{
		DeleteOnTermination: bdm["delete_on_termination"].(bool),
		Size:                bdm["volume_size"].(int),
		Type:                bdm["volume_type"].(string),
	}

	return blockDeviceMapping
}

func readNetworksFromConfig(net map[string]interface{}) *gobizfly.AutoScalingNetworks {
	securityGroupSet := net["security_groups"].(*schema.Set)
	securityGroups := make([]*string, securityGroupSet.Len())
	for _, i := range securityGroupSet.List() {
		securityGroups = append(securityGroups, i.(*string))
	}

	network := &gobizfly.AutoScalingNetworks{
		ID:             net["network_id"].(string),
		SecurityGroups: securityGroups,
	}

	return network
}

func readOperatingSystemFromConfig(os *schema.ResourceData) gobizfly.AutoScalingOperatingSystem {

	return gobizfly.AutoScalingOperatingSystem{
		CreateFrom: os.Get("os.0.create_from").(string),
		ID:         os.Get("os.0.uuid").(string),
	}
}
