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

func dataSourceBizflyCloudKafkaFlavor() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBizflyCloudKafkaFlavorRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter flavors by name",
			},
			"flavors": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of available Kafka flavors",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vcpus": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"ram": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "RAM in MB",
						},
						"disk": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Disk size in GB",
						},
						"is_default": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"flavor_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"billing_plan": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"plan_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceBizflyCloudKafkaFlavorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()

	opts := &gobizfly.KafkaFlavorListOptions{}
	if name, ok := d.GetOk("name"); ok {
		opts.Name = name.(string)
	}

	flavors, err := client.Kafka.ListFlavor(context.Background(), opts)
	if err != nil {
		return fmt.Errorf("error retrieving Kafka flavors: %v", err)
	}

	var flavorList []map[string]interface{}
	for _, flavor := range flavors {
		f := map[string]interface{}{
			"id":           flavor.ID,
			"name":         flavor.Name,
			"code":         flavor.Code,
			"vcpus":        flavor.VCPUs,
			"ram":          flavor.RAM,
			"disk":         flavor.Disk,
			"is_default":   flavor.IsDefault,
			"flavor_type":  flavor.FlavorType,
			"billing_plan": flavor.BillingPlan,
			"plan_name":    flavor.PlanName,
			"description":  flavor.Description,
		}
		flavorList = append(flavorList, f)
	}

	if err := d.Set("flavors", flavorList); err != nil {
		return fmt.Errorf("error setting flavors: %v", err)
	}

	// Set a unique ID for the data source
	d.SetId(fmt.Sprintf("kafka-flavors-%s", d.Get("name").(string)))

	return nil
}
