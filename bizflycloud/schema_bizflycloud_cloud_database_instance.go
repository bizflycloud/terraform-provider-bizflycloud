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

func dataCloudDatabaseInstanceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"autoscaling": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeInt},
		},
		"created_at": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"datastore": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem:     &schema.Resource{Schema: dataElemCloudDatabaseDataStore()},
		},
		"dns": {
			Computed: true,
			Type:     schema.TypeMap,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"enable_failover": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"instance_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"nodes": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"availability_zone": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"created_at": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"id": {
						Type:     schema.TypeString,
						Computed: true,
						Optional: true,
					},
					"name": {
						Type:     schema.TypeString,
						Computed: true,
						Optional: true,
					},
					"operating_status": {
						Type:     schema.TypeString,
						Computed: true,
						Optional: true,
					},

					"region_name": {
						Type:     schema.TypeString,
						Computed: true,
					},

					"role": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"status": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"public_access": {
			Type:     schema.TypeBool,
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
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func resourceCloudDatabaseInstanceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"autoscaling": {
			Type:     schema.TypeMap,
			Required: true,
			Elem:     &schema.Schema{Type: schema.TypeInt},
		},
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
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

		"datastore": {
			Type:     schema.TypeMap,
			Required: true,
			ForceNew: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"type": {
						Type:     schema.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							"MariaDB",
							"MongoDB",
							"MySQL",
							"Postgres",
							"Redis",
						}, false),
					},
					"name": {
						Type:     schema.TypeString,
						Required: true,
					},
					"version_id": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
		"network_ids": {
			Type:     schema.TypeSet,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Required: true,
			ForceNew: true,
		},
		"secondaries": {
			Type:     schema.TypeSet,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"availability_zone": {
						Type:     schema.TypeString,
						Required: true,
					},
					"quantity": {
						Type:     schema.TypeInt,
						Optional: true,
						Default:  1,
						ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
							switch v := val.(int); true {
							case v < 0:
								errs = append(errs, fmt.Errorf("%q must be greater than or equal to 0, got: %d", key, v))
							}
							return
						},
					},
				},
			},
		},
		"public_access": {
			Type:     schema.TypeBool,
			Optional: true,
			ForceNew: true,
		},
		"availability_zone": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"backup_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"configuration_group": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"task_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"status": {
			Type:     schema.TypeString,
			Computed: true,
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
					"availability_zone": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"created_at": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"id": {
						Type:     schema.TypeString,
						Computed: true,
						Optional: true,
					},
					"name": {
						Type:     schema.TypeString,
						Computed: true,
						Optional: true,
					},
					"operating_status": {
						Type:     schema.TypeString,
						Computed: true,
						Optional: true,
					},

					"region_name": {
						Type:     schema.TypeString,
						Computed: true,
					},

					"role": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"status": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"init_databases": {
			Type:     schema.TypeList,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
		},
		"users": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"username": {
						Type:     schema.TypeString,
						Required: true,
					},
					"host": {
						Type:     schema.TypeString,
						Optional: true,
						Default:  "%",
					},
					"password": {
						Type:      schema.TypeString,
						Required:  true,
						Sensitive: true,
					},
					"databases": {
						Type:     schema.TypeList,
						Elem:     &schema.Schema{Type: schema.TypeString},
						Required: true,
					},
				},
			},
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
	}
}
