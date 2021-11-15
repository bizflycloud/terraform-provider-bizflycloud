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

func datasourceBizFlyCloudKubernetesControllerVersions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBizFlyCloudKubernetesVersion,
		Schema: map[string]*schema.Schema{
			"version": {
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
func dataSourceBizFlyCloudKubernetesVersion(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	resp, err := client.KubernetesEngine.GetKubernetesVersion(context.Background())
	if err != nil {
		return err
	}
	versionName := d.Get("version").(string)
	for _, controllerVersion := range resp.ControllerVersions {
		if controllerVersion.K8SVersion == versionName {
			d.SetId(controllerVersion.ID)
			break
		}
	}

	versionID := d.Get("id")
	if versionID == "" {
		return fmt.Errorf("Version %s not found", versionName)
	}
	return nil
}
