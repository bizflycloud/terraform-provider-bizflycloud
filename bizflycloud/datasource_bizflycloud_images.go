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
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceBizFlyCloudImages() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceBizFlyCloudImageRead,
		Schema: imageSchema(),
	}
}

func dataSourceBizFlyCloudImageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	osImages, err := client.Server.ListOSImages(context.Background())
	if err != nil {
		return err
	}
	distribution, okDist := d.GetOk("distribution")
	version, okVer := d.GetOk("version")
	if okDist && okVer {
		for _, image := range osImages {
			if !strings.EqualFold(strings.ToLower(image.OSDistribution), strings.ToLower(distribution.(string))) {
				continue
			}
			for _, v := range image.Version {
				if !strings.EqualFold(strings.ToLower(v.Name), strings.ToLower(version.(string))) {
					continue
				}
				d.SetId(v.ID)
				break
			}
		}
	} else {
		return fmt.Errorf("Distribution and Version must be set")
	}
	return nil
}
