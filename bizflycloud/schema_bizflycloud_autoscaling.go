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
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataAutoScalingGroupSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"desired_capacity": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"launch_configuration_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"launch_configuration_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"launch_configuration_only": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"load_balancers": {
			Type:     schema.TypeList,
			Computed: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: dataAutoScalingLoadBalancerInfoSchema(),
			},
		},
		"max_size": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"min_size": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"node_ids": {
			Type:     schema.TypeSet,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"scale_in_info": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: dataScaleInPolicySchema(),
			},
		},
		"scale_out_info": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: dataScaleOutPolicySchema(),
			},
		},
	}
}
func resourceAutoScalingGroupSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"desired_capacity": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"launch_configuration_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"launch_configuration_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"launch_configuration_only": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"load_balancers": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: resourceAutoScalingLoadBalancerInfoSchema(),
			},
		},
		"max_size": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"min_size": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"node_ids": {
			Type:     schema.TypeSet,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"scale_in_info": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: dataScaleInPolicySchema(),
			},
		},
		"scale_out_info": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: dataScaleOutPolicySchema(),
			},
		},
		"status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"task_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func dataAutoScalingLoadBalancerInfoSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"load_balancer_id": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"server_group_id": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"server_group_port": {
			Computed: true,
			Type:     schema.TypeInt,
		},
	}
}

func resourceAutoScalingLoadBalancerInfoSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"load_balancer_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"server_group_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"server_group_port": {
			Type:     schema.TypeInt,
			Required: true,
		},
	}
}

func dataScalePolicySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cooldown": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"metric_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"range_time": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"threshold": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"scale_size": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	}
}

func dataScaleInPolicySchema() map[string]*schema.Schema {
	commonSchema := dataScalePolicySchema()

	return commonSchema
}

func dataScaleOutPolicySchema() map[string]*schema.Schema {
	commonSchema := dataScalePolicySchema()

	return commonSchema
}

func resourceScalePolicySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cluster_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"cooldown": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"metric_type": {
			Type:     schema.TypeString,
			Required: true,
		},
		"range_time": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"threshold": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"scale_size": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  1,
		},
		"task_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func resourceScaleInPolicySchema() map[string]*schema.Schema {
	commonSchema := resourceScalePolicySchema()
	return commonSchema
}

func resourceScaleOutPolicySchema() map[string]*schema.Schema {
	commonSchema := resourceScalePolicySchema()
	return commonSchema
}

// Profiles
func dataLaunchConfigurationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"availability_zone": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"data_disks": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"delete_on_termination": {
						Type:     schema.TypeBool,
						Computed: true,
					},
					"volume_size": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"volume_type": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"flavor": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"instance_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"network_plan": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"networks": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"network_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"security_groups": {
						Type:     schema.TypeSet,
						Computed: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
				},
			},
		},
		"os": {
			Type:     schema.TypeList,
			Computed: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"create_from": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"error": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"uuid": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"os_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"rootdisk": {
			Type:     schema.TypeList,
			Computed: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"delete_on_termination": {
						Type:     schema.TypeBool,
						Computed: true,
					},
					"volume_size": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"volume_type": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"ssh_key": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"status": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"user_data": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func resourceLaunchConfigurationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"availability_zone": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"data_disks": {
			Type:     schema.TypeList,
			Optional: true,
			ForceNew: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"delete_on_termination": {
						Type:     schema.TypeBool,
						Optional: true,
						Default:  true,
						ForceNew: true,
					},
					"volume_size": {
						Type:     schema.TypeInt,
						Required: true,
						ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
							switch v := val.(int); true {
							case v < 20 || v%10 > 0:
								errs = append(errs, fmt.Errorf("%q must be greater than and divisible by 10, got: %d", key, v))
							case v > 2000:
								errs = append(errs, fmt.Errorf("%q must be greater than and divisible by 10, got: %d", key, v))
							}
							return
						},
					},
					"volume_type": {
						Type:     schema.TypeString,
						Required: true,
						ForceNew: true,
					},
				},
			},
		},
		"flavor": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"instance_type": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"network_plan": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  "free_datatransfer",
			ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
				switch v := val.(string); true {
				case v != "free_datatransfer" && v != "free_bandwidth":
					errs = append(errs, fmt.Errorf("%q must be free_datatransfer or free_bandwidth , got: %s", key, v))
				}
				return
			},
		},
		"networks": {
			Type:     schema.TypeList,
			Optional: true,
			ForceNew: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"network_id": {
						Type:     schema.TypeString,
						Required: true,
					},
					"security_groups": {
						Type:     schema.TypeSet,
						Optional: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
				},
			},
		},
		"os": {
			Type:     schema.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"create_from": {
						Type:     schema.TypeString,
						Required: true,
						ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
							switch v := val.(string); true {
							case v != "image" && v != "snapshot":
								errs = append(errs, fmt.Errorf("%q must be image or snapshot , got: %s", key, v))
							}
							return
						},
					},
					"error": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"uuid": {
						Type:     schema.TypeString,
						Required: true,
					},
					"os_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"rootdisk": {
			Type:     schema.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"delete_on_termination": {
						Type:     schema.TypeBool,
						Optional: true,
						Default:  true,
					},
					"volume_size": {
						Type:     schema.TypeInt,
						Required: true,
					},
					"volume_type": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
		"ssh_key": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"status": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"user_data": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
	}
}
