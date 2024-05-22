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

	"github.com/bizflycloud/gobizfly"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBizflyCloudSSHKey() *schema.Resource {
	return &schema.Resource{
		Create:        resourceBizflyCloudSSHKeyCreate,
		Read:          resourceBizflyCloudSSHKeyRead,
		Delete:        resourceBizflyCloudSSHKeyDelete,
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
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceBizflyCloudSSHKeyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	resp, err := client.CloudServer.SSHKeys().Create(context.Background(), &gobizfly.SSHKeyCreateRequest{
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

func resourceBizflyCloudSSHKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	sshkeys, err := client.CloudServer.SSHKeys().List(context.Background(), &gobizfly.ListOptions{})
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

func resourceBizflyCloudSSHKeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	_, err := client.CloudServer.SSHKeys().Delete(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error deleting ssh key: %v", err)
	}
	return nil
}
