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
	"github.com/YakDriver/regexache"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataCloudDatabaseBackupScheduleSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"first_execution_time": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"instance_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"instance_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"limit_backup": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"next_execution_time": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"node_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"node_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"cron_expression": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func resourceCloudDatabaseBackupScheduleSchema() map[string]*schema.Schema {
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
		"cron_expression": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringMatch(regexache.MustCompile(`^[0-9A-Za-z_\s #()*+,/?^|-]*$`), "see https://crontab.guru"),
		},
	}
}
