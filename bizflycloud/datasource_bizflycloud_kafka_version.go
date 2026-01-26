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
	"fmt"

	"github.com/bizflycloud/gobizfly"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceBizflyCloudKafkaVersion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBizflyCloudKafkaVersionRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter versions by name",
			},
			"versions": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of available Kafka versions",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_default": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceBizflyCloudKafkaVersionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()

	opts := &gobizfly.KafkaVersionListOptions{}
	if name, ok := d.GetOk("name"); ok {
		opts.Name = name.(string)
	}

	versions, err := client.Kafka.ListVersion(context.Background(), opts)
	if err != nil {
		return fmt.Errorf("error retrieving Kafka versions: %v", err)
	}

	var versionList []map[string]interface{}
	for _, version := range versions {
		v := map[string]interface{}{
			"id":         version.ID,
			"code":       version.Code,
			"name":       version.Name,
			"is_default": version.IsDefault,
		}
		versionList = append(versionList, v)
	}

	if err := d.Set("versions", versionList); err != nil {
		return fmt.Errorf("error setting versions: %v", err)
	}

	// Set a unique ID for the data source
	d.SetId(fmt.Sprintf("kafka-versions-%s", d.Get("name").(string)))

	return nil
}
