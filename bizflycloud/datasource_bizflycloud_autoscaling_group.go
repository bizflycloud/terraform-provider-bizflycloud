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

	"github.com/bizflycloud/gobizfly"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceBizFlyCloudAutoScalingGroup() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceBizFlyCloudAutoScalingGroupRead,
		Schema: dataAutoScalingGroupSchema(),
	}
}

func dataSourceBizFlyCloudAutoScalingGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()

	if v, ok := d.GetOk("id"); ok {
		d.SetId(v.(string))
	}

	groupID := d.Id()

	log.Printf("[DEBUG] Reading Autoscaling Group: %s", groupID)
	group, err := client.AutoScaling.AutoScalingGroups().Get(context.Background(), groupID)

	log.Printf("[DEBUG] Checking for error: %s", err)
	if err != nil {
		return fmt.Errorf("error describing AutoScaling Groups: %w", err)
	}

	log.Printf("[DEBUG] Found Autoscaling Group: %s", groupID)
	log.Printf("[DEBUG] bizflycloud_autoscaling_group - Single Auto Scaling Group found: %s", group.Name)

	d.SetId(group.ID)
	_ = d.Set("desired_capacity", group.DesiredCapacity)
	_ = d.Set("launch_configuration_id", group.ProfileID)
	_ = d.Set("launch_configuration_name", group.ProfileName)
	_ = d.Set("max_size", group.MaxSize)
	_ = d.Set("min_size", group.MinSize)
	_ = d.Set("name", group.Name)
	_ = d.Set("node_ids", group.NodeIDs)
	_ = d.Set("status", group.Status)

	if err := d.Set("load_balancers", readLoadBalancerInformation(group.LoadBalancerPolicies)); err != nil {
		return fmt.Errorf("error setting load_balancers: %w", err)
	}

	return nil
}

func readLoadBalancerInformation(l gobizfly.LoadBalancerPolicy) []map[string]interface{} {
	var results []map[string]interface{}
	if l.LoadBalancerID != "" {
		results = append(results, map[string]interface{}{
			"load_balancer_id":  l.LoadBalancerID,
			"server_group_id":   l.ServerGroupID,
			"server_group_port": l.ServerGroupPort,
		})
	}

	return results
}
