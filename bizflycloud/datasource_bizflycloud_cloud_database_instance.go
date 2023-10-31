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
	"encoding/json"
	"fmt"
	"log"

	"github.com/bizflycloud/gobizfly"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceBizflyCloudDatabaseInstance() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceBizflyCloudDatabaseInstanceRead,
		Schema: dataCloudDatabaseInstanceSchema(),
	}
}

func dataSourceBizflyCloudDatabaseInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()

	if v, ok := d.GetOk("id"); ok {
		d.SetId(v.(string))
	}

	instanceID := d.Id()

	log.Printf("[DEBUG] Reading database instance: %s", instanceID)
	instance, err := client.CloudDatabase.Instances().Get(context.Background(), instanceID)

	log.Printf("[DEBUG] Checking for error: %s", err)
	if err != nil {
		return fmt.Errorf("error describing database Instance: %w", err)
	}

	log.Printf("[DEBUG] Found database Instance: %s", instanceID)
	log.Printf("[DEBUG] bizflycloud_cloud_database_instance - Single database Instance found: %s", instance.Name)

	d.SetId(instance.ID)
	_ = d.Set("created_at", instance.CreatedAt)
	_ = d.Set("enable_failover", instance.EnableFailover)
	_ = d.Set("instance_type", instance.InstanceType)
	_ = d.Set("public_access", instance.PublicAccess)
	_ = d.Set("status", instance.Status)
	_ = d.Set("task_id", instance.TaskID)
	_ = d.Set("volume_size", instance.Volume.Size)

	if err := d.Set("autoscaling", readDataCloudDatabaseAutoScaling(instance.AutoScaling)); err != nil {
		return fmt.Errorf("error setting autoscaling: %w", err)
	}

	if err := d.Set("datastore", FlattenStruct(instance.Datastore)); err != nil {
		return fmt.Errorf("error setting datastore: %w", err)
	}

	if err := d.Set("dns", FlattenStruct(instance.DNS)); err != nil {
		return fmt.Errorf("error setting dns: %w", err)
	}

	if err := d.Set("nodes", readNodes(instance.Nodes)); err != nil {
		return fmt.Errorf("error setting nodes: %w", err)
	}
	return nil
}

func readDataStore(d gobizfly.CloudDatabaseDatastore) map[string]interface{} {
	result := map[string]interface{}{
		"type":       d.Type,
		"name":       d.VersionName,
		"version_id": d.VersionID,
	}
	if d.ID != "" {
		result["id"] = d.ID
	}

	return result
}

func readDataCloudDatabaseAutoScaling(as gobizfly.CloudDatabaseAutoScaling) map[string]interface{} {
	m := map[string]interface{}{
		"volume_limited":   as.Volume.Limited,
		"volume_threshold": as.Volume.Threshold,
	}
	m["enable"] = 0
	if as.Enable {
		m["enable"] = 1
	}

	return m
}

func readNodes(nodes []gobizfly.CloudDatabaseNode) []map[string]interface{} {
	var results []map[string]interface{}
	for _, node := range nodes {
		n := readNode(node)
		results = append(results, n)
	}
	return results
}

// FlattenStruct - export to json
func FlattenStruct(structData interface{}) map[string]interface{} {
	var mapData map[string]interface{}
	data, _ := json.Marshal(structData)
	_ = json.Unmarshal(data, &mapData)
	return mapData
}
