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
				Required: false,
			},
			"flavor_name": {
				Type:     schema.TypeString,
				Required: false,
			},
			"ssh_key": {
				Type:     schema.TypeString,
				Optional: false,
			},
			"category": {
				Type:     schema.TypeString,
				Required: false,
			},
			"password": {
				Type:     schema.TypeBool,
				Optional: false,
			},
			"os_type": {
				Type:     schema.TypeString,
				Required: false,
			},
			"os_id": {
				Type:     schema.TypeString,
				Required: false,
			},
			"root_disk_type": {
				Type:     schema.TypeString,
				Required: false,
			},
			"root_disk_size": {
				Type:     schema.TypeInt,
				Required: false,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Required: false,
			},
			"user_id": {
				Type:     schema.TypeString,
				Computed: false,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: false,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: false,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: false,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: false,
			},
			"volume_ids": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: false,
				Computed: false,
			},
			"lan_ip": {
				Type:     schema.TypeString,
				Computed: false,
			},
			"wan_ipv4": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: false,
			},
			"wan_ipv6": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: false,
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
