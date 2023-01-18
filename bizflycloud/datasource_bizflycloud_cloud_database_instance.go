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

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceBizFlyCloudDatabaseInstance() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceBizFlyCloudDatabaseInstanceRead,
		Schema: resourceCloudDatabaseInstanceSchema(),
	}
}

func dataSourceBizFlyCloudDatabaseInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()

	if v, ok := d.GetOk("id"); ok {
		d.SetId(v.(string))
	}

	instanceID := d.Id()

	log.Printf("[DEBUG] Reading Database Instance: %s", instanceID)
	instance, err := client.CloudDatabase.Instances().Get(context.Background(), instanceID)

	log.Printf("[DEBUG] Checking for error: %s", err)
	if err != nil {
		return fmt.Errorf("error describing Database Instance: %w", err)
	}

	log.Printf("[DEBUG] Found Database Instance: %s", instanceID)
	log.Printf("[DEBUG] bizflycloud_cloud_database_instance - Single Database Instance found: %s", instance.Name)

	d.SetId(instance.ID)
	_ = d.Set("status", instance.Status)
	_ = d.Set("nodes", []map[string]interface{}{
		{
			"id": instance.Nodes[0].ID,
		},
	})

	if err := d.Set("dns", ConvertStruct(instance.DNS)); err != nil {
		return fmt.Errorf("error setting dns for instance %s: %s", d.Id(), err)
	}

	return nil
}

// ConvertStruct - export to json
func ConvertStruct(structData interface{}) map[string]interface{} {
	var mapData map[string]interface{}
	data, _ := json.Marshal(structData)
	_ = json.Unmarshal(data, &mapData)
	return mapData
}
