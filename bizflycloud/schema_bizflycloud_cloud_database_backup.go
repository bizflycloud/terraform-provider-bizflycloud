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

func dataCloudDatabaseBackupSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"parent_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"node_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"instance_id": {
			Type:     schema.TypeString,
			Computed: true,
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
			Required: true,
			ForceNew: true,
		},
		"instance_id": {
			Type:     schema.TypeString,
			Required: true,
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
