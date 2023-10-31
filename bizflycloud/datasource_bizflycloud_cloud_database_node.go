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

func datasourceBizflyCloudDatabaseNode() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceBizflyCloudDatabaseNodeRead,
		Schema: dataCloudDatabaseNodeSchema(),
	}
}

func dataSourceBizflyCloudDatabaseNodeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()

	if v, ok := d.GetOk("id"); ok {
		d.SetId(v.(string))
	}

	nodeID := d.Id()

	log.Printf("[DEBUG] Reading database node: %s", nodeID)
	node, err := client.CloudDatabase.Nodes().Get(context.Background(), nodeID)

	log.Printf("[DEBUG] Checking for error: %s", err)
	if err != nil {
		return fmt.Errorf("error describing database node: %w", err)
	}

	log.Printf("[DEBUG] Found database node: %s", nodeID)
	log.Printf("[DEBUG] bizflycloud_cloud_database_node - Single database node found: %s", node.Name)

	d.SetId(node.ID)
	_ = d.Set("availability_zone", node.AvailabilityZone)
	_ = d.Set("created_at", node.CreatedAt)
	_ = d.Set("description", node.Description)
	_ = d.Set("enable_failover", node.EnableFailover)
	_ = d.Set("flavor", node.Flavor)
	_ = d.Set("instance_id", node.InstanceID)
	_ = d.Set("name", node.Name)
	_ = d.Set("node_type", node.NodeType)
	_ = d.Set("operating_status", node.OperatingStatus)
	_ = d.Set("region_name", node.RegionName)
	_ = d.Set("role", node.Role)
	_ = d.Set("status", node.Status)
	_ = d.Set("volume", map[string]interface{}{
		"size": node.Volume.Size,
		"used": node.Volume.Used,
	})

	if node.ReplicaOf != "" {
		_ = d.Set("replica_of", node.ReplicaOf)
	}

	if len(node.Replicas) > 0 {
		_ = d.Set("replicas", readNodeReplicas(node.Replicas))
	}

	if err := d.Set("dns", FlattenStruct(node.DNS)); err != nil {
		return fmt.Errorf("error setting dns for node %s: %s", d.Id(), err)
	}

	if err := d.Set("datastore", FlattenStruct(node.Datastore)); err != nil {
		return fmt.Errorf("error setting datastore for node %s: %s", d.Id(), err)
	}

	addresses := map[string]interface{}{}
	readNodeAddresses(addresses, node.Addresses)
	if err := d.Set("private_addresses", addresses["private_addresses"]); err != nil {
		return fmt.Errorf("error setting private_addresses for node %s: %s", d.Id(), err)
	}
	if err := d.Set("public_addresses", addresses["public_addresses"]); err != nil {
		return fmt.Errorf("error setting public_addresses for node %s: %s", d.Id(), err)
	}
	if err := d.Set("port_access", addresses["port_access"]); err != nil {
		return fmt.Errorf("error setting port_access for node %s: %s", d.Id(), err)
	}

	return nil
}

func readNodeReplicas(replicas []gobizfly.CloudDatabaseNode) []interface{} {
	var results []interface{}
	for _, repl := range replicas {
		n := map[string]interface{}{
			"availability_zone": repl.AvailabilityZone,
			"created_at":        repl.CreatedAt,
			"id":                repl.ID,
			"instance_id":       repl.InstanceID,
		}
		results = append(results, n)
	}
	return results
}
func readNodeAddresses(node map[string]interface{}, addrs gobizfly.CloudDatabaseAddresses) {
	var _private []string
	var _public []string

	for _, addr := range addrs.Private {
		node["port_access"] = addr.Port
		_private = append(_private, addr.IPAddress)
	}

	for _, addr := range addrs.Public {
		_public = append(_public, addr.IPAddress)
	}

	node["private_addresses"] = _private
	node["public_addresses"] = _public
}

func readNode(node gobizfly.CloudDatabaseNode) map[string]interface{} {
	result := map[string]interface{}{
		"availability_zone": node.AvailabilityZone,
		"created_at":        node.CreatedAt,
		// "description":       node.Description,
		// "flavor":            node.Flavor,
		"id": node.ID,
		// "instance_id":       node.InstanceID,
		// "message":           node.Message,
		"name":             node.Name,
		"operating_status": node.OperatingStatus,
		"region_name":      node.RegionName,
		"role":             node.Role,
		// "role":              node.Role,
		"status": node.Status,
		// "task_id":           node.TaskID,
	}
	// result["datastore"] = FlattenStruct(node.Datastore)
	// result["dns"] = FlattenStruct(node.DNS)

	// readNodeAddresses(result, node.Addresses)

	// result["volume"] = map[string]interface{}{
	// 	"size": node.Volume.Size,
	// 	"used": node.Volume.Used,
	// }
	return result
}
