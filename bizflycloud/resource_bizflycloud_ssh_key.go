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

	"github.com/bizflycloud/gobizfly"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBizFlyCloudSSHKey() *schema.Resource {
	return &schema.Resource{
		Create:        resourceBizFlyCloudSSHKeyCreate,
		Read:          resourceBizFlyCloudSSHKeyRead,
		Delete:        resourceBizFlyCloudSSHKeyDelete,
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"public_key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceBizFlyCloudSSHKeyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	resp, err := client.SSHKey.Create(context.Background(), &gobizfly.SSHKeyCreateRequest{
		Name:      d.Get("name").(string),
		PublicKey: d.Get("public_key").(string),
	})
	if err != nil {
		return fmt.Errorf("Error creating SSH key: %v", err)
	}
	d.SetId(resp.Name)
	_ = d.Set("finger_print", resp.FingerPrint)
	return nil
}

func resourceBizFlyCloudSSHKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	sshkeys, err := client.SSHKey.List(context.Background(), &gobizfly.ListOptions{})
	if err != nil {
		return fmt.Errorf("Error retrieving ssh keys: %v", err)
	}
	for _, sshkey := range sshkeys {
		if sshkey.SSHKeyPair.Name == d.Id() {
			// found ssh key
			_ = d.Set("name", sshkey.SSHKeyPair.Name)
			_ = d.Set("public_key", sshkey.SSHKeyPair.PublicKey)
			_ = d.Set("fingerprint", sshkey.SSHKeyPair.FingerPrint)
		}
	}
	return nil
}

func resourceBizFlyCloudSSHKeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	_, err := client.SSHKey.Delete(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error deleting ssh key: %v", err)
	}
	return nil
}
