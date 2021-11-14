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

func dataSourceBizflyClouldSSHKey() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Read: dataSourceBizflyClouldSSHKeyRead,
	}
}

func dataSourceBizflyClouldSSHKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	sshKey, err := client.SSHKey.Get(context.Background(), d.Get("name").(string))
	if err != nil {
		return err
	}
	d.SetId(sshKey.Name)
	err = d.Set("fingerprint", sshKey.FingerPrint)
	if err != nil {
		return err
	}
	err = d.Set("public_key", sshKey.PublicKey)
	if err != nil {
		return err
	}
	return nil

}
