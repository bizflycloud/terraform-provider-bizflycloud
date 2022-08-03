// This file is part of terraform-provider-bizflycloud
//
// Copyright (C) 2022  Bizfly Cloud
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
	"errors"
	"fmt"
	"github.com/bizflycloud/gobizfly"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

func datasourceBizflyCloudVolumeTypes() *schema.Resource {
	return &schema.Resource{
		Read:   datasourceBizflyCloudVolumeTypesRead,
		Schema: dataSourceVolumeTypeSchema(),
	}
}

func datasourceBizflyCloudVolumeTypesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	volumeTypes, err := client.Volume.ListVolumeTypes(context.Background(), &gobizfly.ListVolumeTypesOptions{})
	if err != nil {
		return err
	}
	name := d.Get("name").(string)
	category := d.Get("category").(string)
	err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
		for _, volumeType := range volumeTypes {
			if volumeType.Name == name && volumeType.Category == category {
				d.SetId(volumeType.Type)
				d.Set("name", volumeType.Name)
				d.Set("category", volumeType.Category)
				d.Set("type", volumeType.Type)
				d.Set("availability_zones", volumeType.AvailabilityZones)
				return nil
			}
		}
		return resource.RetryableError(err)
	})
	if !d.IsNewResource() && errors.Is(err, gobizfly.ErrNotFound) {
		log.Printf("[WARN] Volume Type %s is not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("error read Volume Type %s: %w", d.Id(), err)
	}
	return nil
}
