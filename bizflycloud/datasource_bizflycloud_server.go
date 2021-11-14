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
	"github.com/bizflycloud/gobizfly"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceBizFlyCloudServers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBizFlyCloudServerRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"flavor_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ssh_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"category": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"password": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"os_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"os_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"root_disk_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"root_disk_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"volume_ids": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"lan_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"wan_ipv4": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"wan_ipv6": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
		},
	}
}

func dataSourceBizFlyCloudServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	osServers, err := client.Server.List(context.Background(), &gobizfly.ListOptions{
		Page:  0,
		Limit: 1000,
	})
	if err != nil {
		return err
	}
	id, okId := d.GetOk("id")
	if okId {
		for _, server := range osServers {
			if !strings.EqualFold(strings.ToLower(server.ID), strings.ToLower(id.(string))) {
				continue
			}
			d.SetId(server.ID)
			err := resourceBizFlyCloudServerRead(d, meta)
			if err != nil {
				return fmt.Errorf("Couldn't set data")
			}
			break
		}
	} else {
		return fmt.Errorf("Server ID must be set")
	}
	return nil
}
