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
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceCloudDatabaseInstanceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"instance_type": {
			Type:     schema.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				"basic",
				"premium",
				"enterprise",
				"dedicated",
			}, false),
		},
		"flavor_name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"volume_size": {
			Type:     schema.TypeInt,
			Required: true,
			ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
				switch v := val.(int); true {
				case v%10 > 0 || v > 2000:
					errs = append(errs, fmt.Errorf("%q must be greater than 2000 and divisible by 10, got: %d", key, v))
				}
				return
			},
		},
		"datastore_type": {
			Type:     schema.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				"MongoDB",
				"MariaDB",
				"Redis",
			}, false),
		},
		"datastore_version_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"network_ids": {
			Type:     schema.TypeSet,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Required: true,
		},
		"public_access": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"availability_zone": {
			Type:     schema.TypeString,
			Required: true,
		},
		"autoscaling_enable": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"autoscaling_volume_threshold": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  90,
		},
		"autoscaling_volume_limited": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  90,
		},
		"backup_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"task_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"configuration": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"dns": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"nodes": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"id": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
	}
}

func resourceCloudDatabaseNodeSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"replica_of": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"instance_id": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"name": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"configuration": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"role": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				"primary",
				"secondary",
				"replica",
			}, false),
		},
		"private_addresses": {
			Type:     schema.TypeSet,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Computed: true,
		},
		"public_addresses": {
			Type:     schema.TypeSet,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Computed: true,
		},
		"availability_zone": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"created_at": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"datastore": {
			Type:     schema.TypeMap,
			Computed: true,
		},
		"description": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"dns": {
			Type:     schema.TypeMap,
			Computed: true,
		},
		"flavor_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"region_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"task_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"volume_size": {
			Type:     schema.TypeInt,
			Optional: true,
		},
	}
}

func resourceCloudDatabaseBackupSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"parent_id": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"node_id": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"instance_id": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"created": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"datastore": {
			Type:     schema.TypeMap,
			Computed: true,
		},
		"description": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"project_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"size": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"updated": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func resourceCloudDatabaseScheduleSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"node_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"limit_backup": {
			Type:     schema.TypeInt,
			Required: true,
			ForceNew: true,
		},
		"schedule_type": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				"hourly",
				"daily",
				"weekly",
				"monthly",
			}, false),
		},
		"minute": {
			Type:     schema.TypeSet,
			Elem:     &schema.Schema{Type: schema.TypeInt, ValidateFunc: validation.IntBetween(0, 59)},
			Optional: true,
			ForceNew: true,
		},
		"hour": {
			Type:     schema.TypeSet,
			Elem:     &schema.Schema{Type: schema.TypeInt, ValidateFunc: validation.IntBetween(0, 23)},
			Optional: true,
			ForceNew: true,
		},
		"day_of_month": {
			Type:     schema.TypeSet,
			Elem:     &schema.Schema{Type: schema.TypeInt, ValidateFunc: validation.IntBetween(1, 31)},
			Optional: true,
			ForceNew: true,
		},
		"day_of_week": {
			Type:     schema.TypeSet,
			Elem:     &schema.Schema{Type: schema.TypeInt, ValidateFunc: validation.IntBetween(1, 12)},
			Optional: true,
			ForceNew: true,
		},
	}
}

func resourceCloudDatabaseConfigurationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"datastore_type": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				"MongoDB",
				"MariaDB",
				"Redis",
			}, false),
		},
		"datastore_version_id": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"datastore_version_name": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"parameters": {
			Type:     schema.TypeMap,
			Required: true,
		},
	}
}
