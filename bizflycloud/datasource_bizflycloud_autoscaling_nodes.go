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

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceBizFlyCloudAutoscalingNodes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBizFlyCloudNodesRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"nodes": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
					Elem: schema.TypeString,
				},
				Computed: true,
			},
		},
	}
}

func dataSourceBizFlyCloudNodesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	clusterID, okID := d.GetOk("cluster_id")
	// Always list all nodes
	osNodes, err := client.AutoScaling.Nodes().List(context.Background(), clusterID.(string), true)
	if err != nil {
		return err
	}

	if okID {
		nodesResult := make([]map[string]interface{}, len(osNodes))
		for i, node := range osNodes {
			nodesResult[i] = map[string]interface{}{
				"name":         node.Name,
				"id":           node.ID,
				"profile_name": node.ProfileName,
				//TODO: Seting this directly won't output the final output, need someone else to address this later
				//"addresses":     node.Addresses,
				"profile_id":    node.ProfileID,
				"physical_id":   node.PhysicalID,
				"status":        node.Status,
				"status_reason": node.StatusReason,
			}
		}
		d.SetId(clusterID.(string))
		_ = d.Set("nodes", nodesResult)
	} else {
		return fmt.Errorf("Nodes ID must be set")
	}

	return nil
}
