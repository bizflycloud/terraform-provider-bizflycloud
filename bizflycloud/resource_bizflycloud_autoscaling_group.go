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
	"strconv"
	"time"

	"github.com/bizflycloud/gobizfly"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBizFlyCloudAutoscalingGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceBizFlyCloudAutoscalingGroupCreate,
		Read:   resourceBizFlyCloudAutoscalingGroupRead,
		Update: resourceBizFlyCloudAutoscalingGroupUpdate,
		Delete: resourceBizFlyCloudAutoscalingGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: resourceAutoScalingGroupSchema(),
	}
}

func resourceBizFlyCloudAutoscalingGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()

	ascr := &gobizfly.AutoScalingGroupCreateRequest{
		DesiredCapacity: d.Get("desired_capacity").(int),
		MaxSize:         d.Get("max_size").(int),
		MinSize:         d.Get("min_size").(int),
		Name:            d.Get("name").(string),
		ProfileID:       d.Get("launch_configuration_id").(string),
	}

	if v, ok := d.GetOk("load_balancers"); ok && len(v.([]interface{})) > 0 {
		ascr.LoadBalancerPolicies = readLoadBalancersFromConfig(d)
	}

	if _, ok := d.GetOk("scale_in_info"); ok {
		ascr.ScaleInPolicies = &[]gobizfly.ScalePolicy{}
	}
	if _, ok := d.GetOk("scale_out_info"); ok {
		ascr.ScaleInPolicies = &[]gobizfly.ScalePolicy{}
	}

	task, err := client.AutoScaling.AutoScalingGroups().Create(context.Background(), ascr)
	if err != nil {
		return fmt.Errorf("[ERROR] create auto scaling group %s failed: %s", d.Get("name"), err)
	}

	log.Printf("[DEBUG] creating auto scaling group with task %s", task.TaskID)

	d.SetId(task.ID)
	_ = d.Set("task_id", task.TaskID)

	// wait for auto scaling group to become active
	_, err = waitForAutoScalingGroupReady(d, meta)
	if err != nil {
		return fmt.Errorf("[ERROR] create auto scaling group (%s) failed: %s", d.Get("name").(string), err)
	}

	return resourceBizFlyCloudAutoscalingGroupRead(d, meta)
}

func resourceBizFlyCloudAutoscalingGroupRead(d *schema.ResourceData, meta interface{}) error {
	if err := dataSourceBizFlyCloudAutoScalingGroupRead(d, meta); err != nil {
		return err
	}

	return nil
}

func resourceBizFlyCloudAutoscalingGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()

	asur := &gobizfly.AutoScalingGroupUpdateRequest{
		DesiredCapacity: d.Get("desired_capacity").(int),
		MaxSize:         d.Get("max_size").(int),
		MinSize:         d.Get("min_size").(int),
		Name:            d.Get("name").(string),
		ProfileID:       d.Get("launch_configuration_id").(string),
		ProfileOnly:     d.Get("launch_configuration_only").(bool),
	}

	if v, ok := d.GetOk("load_balancers"); ok && len(v.([]interface{})) > 0 {
		asur.LoadBalancerPolicies = readLoadBalancersFromConfig(d)
	}

	// wait for auto scaling group to become active
	_, err := waitForAutoScalingGroupReady(d, meta)
	if err != nil {
		return fmt.Errorf("[ERROR] updating auto scaling group (%s) failed: %s", d.Get("name").(string), err)
	}

	task, err := client.AutoScaling.AutoScalingGroups().Update(context.Background(), d.Id(), asur)
	if err != nil {
		return fmt.Errorf("[ERROR] update auto scaling group %s failed: %s", d.Get("name"), err)
	}

	log.Printf("[DEBUG] updating auto scaling group with task %s", task.TaskID)
	d.SetId(task.ID)
	_ = d.Set("task_id", task.TaskID)

	// wait for auto scaling group to become active
	_, err = waitForAutoScalingGroupReady(d, meta)
	if err != nil {
		return fmt.Errorf("[ERROR] updating auto scaling group (%s) failed: %s", d.Get("name").(string), err)
	}

	return resourceBizFlyCloudAutoscalingGroupRead(d, meta)
}

func resourceBizFlyCloudAutoscalingGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()

	err := client.AutoScaling.AutoScalingGroups().Delete(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error delete auto scaling group %v", err)
	}

	return nil
}

func waitForAutoScalingGroupReady(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	log.Printf("[INFO] Waiting for auto scaling group (%s) to be ready", d.Get("name").(string))
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"CREATING", "RESIZING", "UPDATING"},
		Target:     []string{"ACTIVE", "ERROR"},
		Refresh:    newStateRefreshfunc(d, "status", meta),
		Timeout:    3600 * time.Second,
		Delay:      20 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	return stateConf.WaitForState()
}

func newStateRefreshfunc(d *schema.ResourceData, attribute string, meta interface{}) resource.StateRefreshFunc {
	client := meta.(*CombinedConfig).gobizflyClient()
	return func() (interface{}, string, error) {
		resp, err := client.AutoScaling.Tasks().Get(context.Background(), d.Get("task_id").(string))
		if err != nil {
			return nil, "", err
		}
		// if the task is not ready, we need to wait for a moment
		if !resp.Ready && len(resp.Result.Action) > 0 {
			log.Println("[DEBUG] auto scaling is not ready")
			return nil, "", nil
		}

		err = resourceBizFlyCloudAutoscalingGroupRead(d, meta)
		if err != nil {
			return nil, "", err
		}

		if attr, ok := d.GetOk(attribute); ok {
			asg, err := client.AutoScaling.AutoScalingGroups().Get(context.Background(), d.Id())
			if err != nil {
				return nil, "", fmt.Errorf("Error retrieving auto scaling group: %v", err)
			}
			switch attr := attr.(type) {
			case bool:
				return &asg, strconv.FormatBool(attr), nil
			default:
				return &asg, attr.(string), nil
			}
		}
		return nil, "", nil
	}
}

func readLoadBalancersFromConfig(l *schema.ResourceData) *gobizfly.LoadBalancerPolicy {
	return &gobizfly.LoadBalancerPolicy{
		LoadBalancerID:  l.Get("load_balancers.0.load_balancer_id").(string),
		ServerGroupID:   l.Get("load_balancers.0.server_group_id").(string),
		ServerGroupPort: l.Get("load_balancers.0.server_group_port").(int),
	}
}
