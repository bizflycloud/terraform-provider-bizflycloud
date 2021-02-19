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

	"time"

	"github.com/bizflycloud/gobizfly"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	requestPerSecond = "request_per_second"
)

func resourceBizFlyCloudAutoscalingScaleInPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceBizFlyCloudAutoscalingScaleInPolicyCreate,
		Read:   resourceBizFlyCloudAutoscalingScaleInPolicyRead,
		Update: resourceBizFlyCloudAutoscalingScaleInPolicyUpdate,
		Delete: resourceBizFlyCloudAutoscalingScalePolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: resourceScaleInPolicySchema(),
	}
}

func resourceBizFlyCloudAutoscalingScaleOutPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceBizFlyCloudAutoscalingScaleOutPolicyCreate,
		Read:   resourceBizFlyCloudAutoscalingScaleOutPolicyRead,
		Update: resourceBizFlyCloudAutoscalingScaleOutPolicyUpdate,
		Delete: resourceBizFlyCloudAutoscalingScalePolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: resourceScaleOutPolicySchema(),
	}
}

func resourceBizFlyCloudAutoscalingScaleInPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	clusterID := d.Get("cluster_id").(string)

	if _, err := waitForAutoScalingGroupPolicyAvailableInteractive(d, meta); err != nil {
		return fmt.Errorf("[ERROR] error when wait scaling policy available to interactive (%s): %s", d.Get("cluster_id"), err)
	}

	if d.Get("metric_type").(string) == requestPerSecond {
		policies, err := client.AutoScaling.Policies().List(context.Background(), clusterID)
		if err != nil {
			return fmt.Errorf("[ERROR] errors when create scale in policy for cluster: %s, error: %s", clusterID, err)
		}

		lb := policies.LoadBalancerPolicyInformations

		lbpcr := &gobizfly.LoadBalancersPolicyCreateRequest{
			CoolDown:   d.Get("cooldown").(int),
			Event:      "CLUSTER_SCALE_IN",
			MetricType: d.Get("metric_type").(string),
			RangeTime:  d.Get("range_time").(int),
			LoadBalancers: gobizfly.LoadBalancerScalingPolicy{
				ID:         lb.LoadBalancerID,
				TargetID:   lb.ServerGroupID,
				TargetType: "backend",
			},
			ScaleSize: d.Get("scale_size").(int),
			Threshold: d.Get("threshold").(int),
		}
		task, err := client.AutoScaling.Policies().CreateLoadBalancers(context.Background(), clusterID, lbpcr)
		if err != nil {
			return fmt.Errorf("[ERROR] errors when create scale in policy for cluster: %s, error: %s", clusterID, err)
		}

		_ = d.Set("task_id", task.TaskID)
	} else {

		pcr := &gobizfly.PolicyCreateRequest{
			CoolDown:   d.Get("cooldown").(int),
			Event:      "CLUSTER_SCALE_IN",
			MetricType: d.Get("metric_type").(string),
			RangeTime:  d.Get("range_time").(int),
			ScaleSize:  d.Get("scale_size").(int),
			Threshold:  d.Get("threshold").(int),
		}

		task, err := client.AutoScaling.Policies().Create(context.Background(), clusterID, pcr)
		if err != nil {
			return fmt.Errorf("[ERROR] errors when create scale in policy for cluster: %s, error: %s", clusterID, err)
		}

		_ = d.Set("task_id", task.TaskID)
	}

	_, err := waitForAutoScalingGroupPolicyReady(d, meta)
	if err != nil {
		return fmt.Errorf("[ERROR] errors when create scale in policy for cluster: %s, error: %s", clusterID, err)
	}

	return resourceBizFlyCloudAutoscalingScaleInPolicyRead(d, meta)
}

func resourceBizFlyCloudAutoscalingScaleInPolicyRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceBizFlyCloudAutoscalingScaleInPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	clusterID := d.Get("cluster_id").(string)
	policyID := d.Id()

	if d.HasChange("metric_type") {
		return fmt.Errorf("[ERROR] value of metric_type is not allow change")
	}

	if _, err := waitForAutoScalingGroupPolicyAvailableInteractive(d, meta); err != nil {
		return fmt.Errorf("[ERROR] error when wait scaling policy available to interactive (%s): %s", d.Get("cluster_id"), err)
	}

	if d.Get("metric_type").(string) == requestPerSecond {
		policies, err := client.AutoScaling.Policies().List(context.Background(), clusterID)
		if err != nil {
			return fmt.Errorf("[ERROR] errors when update scale in policy for cluster: %s, error: %s", clusterID, err)
		}

		lb := policies.LoadBalancerPolicyInformations

		lbpur := &gobizfly.LoadBalancersPolicyUpdateRequest{
			CoolDown:   d.Get("cooldown").(int),
			Event:      "CLUSTER_SCALE_IN",
			MetricType: d.Get("metric_type").(string),
			RangeTime:  d.Get("range_time").(int),
			LoadBalancers: gobizfly.LoadBalancerScalingPolicy{
				ID:         lb.LoadBalancerID,
				TargetID:   lb.ServerGroupID,
				TargetType: "backend",
			},
			ScaleSize: d.Get("scale_size").(int),
			Threshold: d.Get("threshold").(int),
		}
		task, err := client.AutoScaling.Policies().UpdateLoadBalancers(context.Background(), clusterID, policyID, lbpur)
		if err != nil {
			return fmt.Errorf("[ERROR] errors when update scale in policy for cluster: %s, error: %s", clusterID, err)
		}

		_ = d.Set("task_id", task.TaskID)
	} else {

		pur := &gobizfly.PolicyUpdateRequest{
			CoolDown:   d.Get("cooldown").(int),
			Event:      "CLUSTER_SCALE_IN",
			MetricType: d.Get("metric_type").(string),
			RangeTime:  d.Get("range_time").(int),
			ScaleSize:  d.Get("scale_size").(int),
			Threshold:  d.Get("threshold").(int),
		}

		task, err := client.AutoScaling.Policies().Update(context.Background(), clusterID, policyID, pur)
		if err != nil {
			return fmt.Errorf("[ERROR] errors when update scale in policy for cluster: %s, error: %s", clusterID, err)
		}

		_ = d.Set("task_id", task.TaskID)
	}

	_, err := waitForAutoScalingGroupPolicyReady(d, meta)
	if err != nil {
		return fmt.Errorf("[ERROR] errors when update scale in policy for cluster: %s, error: %s", clusterID, err)
	}

	return resourceBizFlyCloudAutoscalingScaleInPolicyRead(d, meta)
}

func resourceBizFlyCloudAutoscalingScaleOutPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	clusterID := d.Get("cluster_id").(string)

	if _, err := waitForAutoScalingGroupPolicyAvailableInteractive(d, meta); err != nil {
		return fmt.Errorf("[ERROR] error when wait scaling policy available to interactive (%s): %s", d.Get("cluster_id"), err)
	}

	if d.Get("metric_type").(string) == requestPerSecond {
		policies, err := client.AutoScaling.Policies().List(context.Background(), clusterID)
		if err != nil {
			return fmt.Errorf("[ERROR] errors when create scale out policy for cluster: %s, error: %s", clusterID, err)
		}

		lb := policies.LoadBalancerPolicyInformations

		lbpcr := &gobizfly.LoadBalancersPolicyCreateRequest{
			CoolDown:   d.Get("cooldown").(int),
			Event:      "CLUSTER_SCALE_OUT",
			MetricType: d.Get("metric_type").(string),
			RangeTime:  d.Get("range_time").(int),
			LoadBalancers: gobizfly.LoadBalancerScalingPolicy{
				ID:         lb.LoadBalancerID,
				TargetID:   lb.ServerGroupID,
				TargetType: "backend",
			},
			ScaleSize: d.Get("scale_size").(int),
			Threshold: d.Get("threshold").(int),
		}
		task, err := client.AutoScaling.Policies().CreateLoadBalancers(context.Background(), clusterID, lbpcr)
		if err != nil {
			return fmt.Errorf("[ERROR] errors when create scale out policy for cluster: %s, error: %s", clusterID, err)
		}

		_ = d.Set("task_id", task.TaskID)
	} else {

		pcr := &gobizfly.PolicyCreateRequest{
			CoolDown:   d.Get("cooldown").(int),
			Event:      "CLUSTER_SCALE_OUT",
			MetricType: d.Get("metric_type").(string),
			RangeTime:  d.Get("range_time").(int),
			ScaleSize:  d.Get("scale_size").(int),
			Threshold:  d.Get("threshold").(int),
		}

		task, err := client.AutoScaling.Policies().Create(context.Background(), clusterID, pcr)
		if err != nil {
			return fmt.Errorf("[ERROR] errors when create scale out policy for cluster: %s, error: %s", clusterID, err)
		}

		_ = d.Set("task_id", task.TaskID)
	}

	_, err := waitForAutoScalingGroupPolicyReady(d, meta)
	if err != nil {
		return fmt.Errorf("[ERROR] errors when create scale out policy for cluster: %s, error: %s", clusterID, err)
	}

	return resourceBizFlyCloudAutoscalingScaleOutPolicyRead(d, meta)
}

func resourceBizFlyCloudAutoscalingScaleOutPolicyRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceBizFlyCloudAutoscalingScaleOutPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	clusterID := d.Get("cluster_id").(string)
	policyID := d.Id()

	if d.HasChange("metric_type") {
		return fmt.Errorf("[ERROR] value of metric_type is not allow change")
	}

	if _, err := waitForAutoScalingGroupPolicyAvailableInteractive(d, meta); err != nil {
		return fmt.Errorf("[ERROR] error when wait scaling policy available to interactive (%s): %s", d.Get("cluster_id"), err)
	}

	if d.Get("metric_type").(string) == requestPerSecond {
		policies, err := client.AutoScaling.Policies().List(context.Background(), clusterID)
		if err != nil {
			return fmt.Errorf("[ERROR] errors when update scale out policy for cluster: %s, error: %s", clusterID, err)
		}

		lb := policies.LoadBalancerPolicyInformations

		lbpur := &gobizfly.LoadBalancersPolicyUpdateRequest{
			CoolDown:   d.Get("cooldown").(int),
			Event:      "CLUSTER_SCALE_OUT",
			MetricType: d.Get("metric_type").(string),
			RangeTime:  d.Get("range_time").(int),
			LoadBalancers: gobizfly.LoadBalancerScalingPolicy{
				ID:         lb.LoadBalancerID,
				TargetID:   lb.ServerGroupID,
				TargetType: "backend",
			},
			ScaleSize: d.Get("scale_size").(int),
			Threshold: d.Get("threshold").(int),
		}
		task, err := client.AutoScaling.Policies().UpdateLoadBalancers(context.Background(), clusterID, policyID, lbpur)
		if err != nil {
			return fmt.Errorf("[ERROR] errors when update scale out policy for cluster: %s, error: %s", clusterID, err)
		}

		_ = d.Set("task_id", task.TaskID)
	} else {

		pur := &gobizfly.PolicyUpdateRequest{
			CoolDown:   d.Get("cooldown").(int),
			Event:      "CLUSTER_SCALE_OUT",
			MetricType: d.Get("metric_type").(string),
			RangeTime:  d.Get("range_time").(int),
			ScaleSize:  d.Get("scale_size").(int),
			Threshold:  d.Get("threshold").(int),
		}

		task, err := client.AutoScaling.Policies().Update(context.Background(), clusterID, policyID, pur)
		if err != nil {
			return fmt.Errorf("[ERROR] errors when update scale out policy for cluster: %s, error: %s", clusterID, err)
		}

		_ = d.Set("task_id", task.TaskID)
	}

	_, err := waitForAutoScalingGroupPolicyReady(d, meta)
	if err != nil {
		return fmt.Errorf("[ERROR] errors when update scale out policy for cluster: %s, error: %s", clusterID, err)
	}

	return resourceBizFlyCloudAutoscalingScaleInPolicyRead(d, meta)
}

func resourceBizFlyCloudAutoscalingScalePolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	clusterID := d.Get("cluster_id").(string)
	policyID := d.Id()

	if _, err := waitForAutoScalingGroupPolicyAvailableInteractive(d, meta); err != nil {
		return fmt.Errorf("[ERROR] error when wait scaling policy available to interactive (%s): %s", d.Get("cluster_id"), err)
	}

	if err := client.AutoScaling.Policies().Delete(context.Background(), clusterID, policyID); err != nil {
		log.Printf("[WARNING] errors when delete scale policy %s for cluster: %s, error: %s", policyID, clusterID, err)
	}

	return nil
}

// Wait other tasks done
func waitForAutoScalingGroupPolicyAvailableInteractive(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	client := meta.(*CombinedConfig).gobizflyClient()

	log.Printf("[INFO] Waiting for scaling policy for (%s) to be available", d.Get("cluster_id").(string))
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"DONE"},
		Refresh: func() (interface{}, string, error) {
			policies, err := client.AutoScaling.Policies().List(context.Background(), d.Get("cluster_id").(string))
			if err != nil {
				return nil, "", fmt.Errorf("[ERROR] error when wait scaling policy available to interactive (%s): %s", d.Get("cluster_id"), err)
			}

			if len(policies.DoingTasks) > 0 {
				return &policies, "PENDING", nil
			}

			return &policies, "DONE", nil
		},
		Timeout:    3600 * time.Second,
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
	}
	return stateConf.WaitForState()
}

// Wait to create new done
func waitForAutoScalingGroupPolicyReady(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	log.Printf("[INFO] Waiting for create scaling policy for (%s) to be ready", d.Get("cluster_id").(string))
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING"},
		Target:     []string{"DONE"},
		Refresh:    newStateRefreshPolicyfunc(d, "ready", meta),
		Timeout:    3600 * time.Second,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}
	return stateConf.WaitForState()
}

func newStateRefreshPolicyfunc(d *schema.ResourceData, attribute string, meta interface{}) resource.StateRefreshFunc {
	client := meta.(*CombinedConfig).gobizflyClient()
	return func() (interface{}, string, error) {
		resp, err := client.AutoScaling.Tasks().Get(context.Background(), d.Get("task_id").(string))
		if err != nil {
			return nil, "", fmt.Errorf("[ERROR] error when wait task %s done: %s", d.Get("task_id"), err)
		}

		if !resp.Ready {
			return &resp, "PENDING", nil
		}

		// Set policy_id to d
		d.SetId(resp.Result.Data.(map[string]interface{})["policy_id"].(string))

		return &resp, "DONE", nil
	}
}
