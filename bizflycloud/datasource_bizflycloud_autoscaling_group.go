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
	d.Set("desired_capacity", group.DesiredCapacity)
	d.Set("launch_configuration_id", group.ProfileID)
	d.Set("launch_configuration_name", group.ProfileName)
	d.Set("max_size", group.MaxSize)
	d.Set("min_size", group.MinSize)
	d.Set("name", group.Name)
	d.Set("node_ids", group.NodeIDs)
	d.Set("status", group.Status)

	if err := d.Set("load_balancers", readLoadBalancerInfo(group.LoadBalancerPolicyInformations)); err != nil {
		return fmt.Errorf("error setting load_balancers: %w", err)
	}

	if err := d.Set("scale_in_info", readScaleInPolicyInformation(group.ScaleInPolicyInformations)); err != nil {
		return fmt.Errorf("error setting scale_in_info: %w", err)
	}

	if err := d.Set("scale_out_info", readScaleOutPolicyInformation(group.ScaleOutPolicyInformations)); err != nil {
		return fmt.Errorf("error setting scale_out_info: %w", err)
	}
	return nil
}

func readLoadBalancerInfo(l gobizfly.LoadBalancerPolicyInformation) []map[string]interface{} {
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

func readScaleInPolicyInformation(policies []gobizfly.ScalePolicyInformation) []map[string]interface{} {
	var results []map[string]interface{}
	for _, p := range policies {
		results = append(results, map[string]interface{}{
			"cooldown":    p.CoolDown,
			"metric_type": p.MetricType,
			"range_time":  p.RangeTime,
			"scale_size":  p.ScaleSize,
			"threshold":   p.Threshold,
		})
	}
	return results
}

func readScaleOutPolicyInformation(policies []gobizfly.ScalePolicyInformation) []map[string]interface{} {
	var results []map[string]interface{}
	for _, p := range policies {
		results = append(results, map[string]interface{}{
			"cooldown":    p.CoolDown,
			"metric_type": p.MetricType,
			"range_time":  p.RangeTime,
			"scale_size":  p.ScaleSize,
			"threshold":   p.Threshold,
		})
	}
	return results
}
