// This file is part of terraform-provider-bizflycloud
//
// Copyright (C) 2026  Bizfly Cloud
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
	"errors"
	"fmt"
	"log"

	"github.com/bizflycloud/gobizfly"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceBizflyCloudKafka() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBizflyCloudKafkaRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"version_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nodes": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"flavor": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"volume_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_network_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_access": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"obs_dashboard_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceBizflyCloudKafkaRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()

	// If ID provided, prefer Get by ID
	if id, ok := d.GetOk("id"); ok && id.(string) != "" {
		cluster, err := client.Kafka.Get(context.Background(), id.(string))
		if err != nil {
			if errors.Is(err, gobizfly.ErrNotFound) {
				d.SetId("")
				return nil
			}
			return fmt.Errorf("error getting kafka cluster %s: %w", id.(string), err)
		}
		d.SetId(cluster.ID)
		_ = d.Set("name", cluster.Name)
		_ = d.Set("version_id", cluster.KafkaVersion)
		_ = d.Set("nodes", len(cluster.Nodes))
		_ = d.Set("flavor", cluster.Flavor)
		_ = d.Set("volume_size", cluster.VolumeSize)
		_ = d.Set("status", cluster.Status)
		_ = d.Set("availability_zone", cluster.AvailabilityZone)
		// ClusterResponse doesn't expose VPCNetworkID via Get; set empty
		_ = d.Set("vpc_network_id", "")
		_ = d.Set("public_access", cluster.PublicAccess)
		_ = d.Set("obs_dashboard_url", cluster.OBS.DashboardURL)
		return nil
	}

	// Otherwise try to lookup by name
	name := d.Get("name").(string)
	clusters, err := client.Kafka.List(context.Background(), &gobizfly.KafkaClusterListOptions{Name: name})
	if err != nil {
		return fmt.Errorf("error listing kafka clusters: %w", err)
	}
	if len(clusters) == 0 {
		return fmt.Errorf("no kafka cluster found with name %s", name)
	}
	// pick first match
	match := clusters[0]
	d.SetId(match.ID)
	_ = d.Set("name", match.Name)
	_ = d.Set("version_id", match.KafkaVersion)
	_ = d.Set("nodes", match.Nodes)
	_ = d.Set("flavor", match.Flavor)
	_ = d.Set("volume_size", match.VolumeSize)
	_ = d.Set("status", match.Status)
	_ = d.Set("availability_zone", match.AvailabilityZone)
	_ = d.Set("vpc_network_id", match.VPCNetworkID)
	_ = d.Set("public_access", match.PublicAccess)
	_ = d.Set("obs_dashboard_url", match.OBS.DashboardURL)
	log.Printf("[DEBUG] Found kafka cluster %s (id=%s)", match.Name, match.ID)
	return nil
}
