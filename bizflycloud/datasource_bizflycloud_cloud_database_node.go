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
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceBizFlyCloudDatabaseNode() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceBizFlyCloudDatabaseNodeRead,
		Schema: resourceCloudDatabaseNodeSchema(),
	}
}

func dataSourceBizFlyCloudDatabaseNodeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()

	if v, ok := d.GetOk("id"); ok {
		d.SetId(v.(string))
	}

	nodeID := d.Id()

	log.Printf("[DEBUG] Reading Database node: %s", nodeID)
	node, err := client.CloudDatabase.Nodes().Get(context.Background(), nodeID)

	log.Printf("[DEBUG] Checking for error: %s", err)
	if err != nil {
		return fmt.Errorf("error describing Database node: %w", err)
	}

	log.Printf("[DEBUG] Found Database node: %s", nodeID)
	log.Printf("[DEBUG] bizflycloud_cloud_database_node - Single Database node found: %s", node.Name)

	d.SetId(node.ID)
	_ = d.Set("status", node.Status)
	_ = d.Set("availability_zone", node.AvailabilityZone)
	_ = d.Set("created_at", node.CreatedAt)
	_ = d.Set("description", node.Description)
	_ = d.Set("region_name", node.RegionName)
	_ = d.Set("volume_size", node.Volume.Size)

	if err := d.Set("dns", ConvertStruct(node.DNS)); err != nil {
		return fmt.Errorf("error setting dns for node %s: %s", d.Id(), err)
	}
	if err := d.Set("datastore", ConvertStruct(node.Datastore)); err != nil {
		return fmt.Errorf("error setting datastore for node %s: %s", d.Id(), err)
	}
	if err := d.Set("private_addresses", flatternCloudDatabaseAddresses(node.Addresses.Private)); err != nil {
		return fmt.Errorf("error setting private_addresses for node %s: %s", d.Id(), err)
	}
	if err := d.Set("public_addresses", flatternCloudDatabaseAddresses(node.Addresses.Public)); err != nil {
		return fmt.Errorf("error setting public_addresses for node %s: %s", d.Id(), err)
	}

	return nil
}

func flatternCloudDatabaseAddresses(addresses []gobizfly.CloudDatabaseAddressesDetail) *schema.Set {
	flatternAddresses := schema.NewSet(schema.HashString, []interface{}{})
	for _, address := range addresses {
		flatternAddresses.Add(address.IPAddress)
	}
	return flatternAddresses
}
