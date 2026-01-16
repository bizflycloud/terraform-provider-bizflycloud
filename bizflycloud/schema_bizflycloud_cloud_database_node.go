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
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataCloudDatabaseNodeSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"availability_zone": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"created_at": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"datastore": {
			Computed: true,
			Type:     schema.TypeMap,
			Elem:     &schema.Resource{Schema: dataElemCloudDatabaseDataStore()},
		},
		"description": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"dns": {
			Computed: true,
			Type:     schema.TypeMap,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"flavor": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"id": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
		"instance_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"message": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
		"node_type": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
		"operating_status": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},

		"port_access": {
			Computed: true,
			Type:     schema.TypeInt,
		},
		"private_addresses": {
			Computed: true,
			Type:     schema.TypeSet,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"public_addresses": {
			Computed: true,
			Type:     schema.TypeSet,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"region_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"replica_of": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"replicas": {
			Computed: true,
			Type:     schema.TypeList,
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
					},
					"instance_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"role": {
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
		"volume": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeFloat},
		},
	}
}

func resourceCloudDatabaseNodeSchema() map[string]*schema.Schema {
	s := dataCloudDatabaseNodeSchema()
	// Add updatable fields
	s["volume_size"] = &schema.Schema{
		Type:     schema.TypeInt,
		Optional: true,
	}
	s["flavor_name"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	// Make id required for resource
	s["id"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
	return s
}
