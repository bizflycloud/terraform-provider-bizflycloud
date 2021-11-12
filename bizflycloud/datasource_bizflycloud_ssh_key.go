// This file is part of terraform-provider-bizflycloud
//
// Copyright (C) 2021  BizFly Cloud
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
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceBizFlyCloudSSHKey() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBizflyCloudSSHKey,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"public_key": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceBizflyCloudSSHKey(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	name := d.Get("name").(string)
	resp, err := client.SSHKey.Get(context.Background(), name)
	if err != nil {

		return err
	}
	err = d.Set("public_key", resp.PublicKey)
	if err != nil {
		return err
	}
	err = d.Set("fingerprint", resp.FingerPrint)
	if err != nil {
		return err
	}
	err = d.Set("name", resp.Name)
	if err != nil {
		return err
	}
	return nil
}
