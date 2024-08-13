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

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceBizflyCloudKubernetesControllerPackage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBizflyCloudKubernetesPackage,
		Schema: map[string]*schema.Schema{
			"provision_type": {
				Type: schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func dataSourceBizflyCloudKubernetesPackage(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	provisionType := d.Get("provision_type").(string)
	resp, err := client.KubernetesEngine.GetPackages(context.Background(), provisionType)
	if err != nil {
		return err
	}
	packageName := d.Get("name").(string)

	for _, pkg := range resp.Packages {
		if pkg.Name == packageName {
			d.SetId(pkg.ID)
			break
		}
	}

	packageID := d.Get("id")
	if packageID == "" {
		return fmt.Errorf("Package %s not found", packageName)
	}
	return nil
}
