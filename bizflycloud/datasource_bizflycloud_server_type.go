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

func dataSourceBizflyCloudServerTypes() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceBizflyCloudServerTypeRead,
		Schema: dataSourceServerTypeSchema(),
	}
}

func dataSourceBizflyCloudServerTypeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	var matchServerType *gobizfly.ServerType
	err := resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
		serverTypes, err := client.CloudServer.ListServerTypes(context.Background())
		if err != nil {
			return resource.RetryableError(err)
		}
		for _, serverType := range serverTypes {
			if serverType.Name == d.Get("name").(string) {
				matchServerType = serverType
				break
			}
		}
		if matchServerType == nil {
			return resource.RetryableError(fmt.Errorf("Server type %s not found", d.Get("name").(string)))
		}
		return nil
	})
	if !d.IsNewResource() && errors.Is(err, gobizfly.ErrNotFound) {
		log.Printf("[WARN] Server type %s not found, removing from state", d.Get("name").(string))
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("error reading server type %s: %s", d.Get("name").(string), err)
	}
	d.SetId(matchServerType.ID)
	_ = d.Set("name", matchServerType.Name)
	_ = d.Set("enabled", matchServerType.Enabled)
	_ = d.Set("compute_class", matchServerType.ComputeClass)
	_ = d.Set("priority", matchServerType.Priority)
	return nil
}
