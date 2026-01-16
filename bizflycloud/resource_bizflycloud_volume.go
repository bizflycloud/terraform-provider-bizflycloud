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
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/bizflycloud/gobizfly"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBizflyCloudVolume() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Create: resourceBizflyCloudVolumeCreate,
		Read:   resourceBizflyCloudVolumeRead,
		Update: resourceBizflyCloudVolumeUpdate,
		Delete: resourceBizflyCloudVolumeDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
		},
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"category": {
				Type:     schema.TypeString,
				Required: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceBizflyCloudVolumeCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	vcr := &gobizfly.VolumeCreateRequest{
		Name:             d.Get("name").(string),
		Size:             d.Get("size").(int),
		VolumeType:       d.Get("type").(string),
		VolumeCategory:   d.Get("category").(string),
		AvailabilityZone: d.Get("availability_zone").(string),
	}
	log.Printf("[DEBUG] Create Volume configuration #{vcr}")
	volume, err := client.CloudServer.Volumes().Create(context.Background(), vcr)
	if err != nil {
		return fmt.Errorf("error creating volume: %v", err)
	}
	d.SetId(volume.ID)
	err = resourceBizflyCloudVolumeRead(d, meta)
	if err != nil {
		return fmt.Errorf("error retrieving volume: %v", err)
	}
	return nil
}

func resourceBizflyCloudVolumeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	volume, err := client.CloudServer.Volumes().Get(context.Background(), d.Id())
	if err != nil {
		if errors.Is(err, gobizfly.ErrNotFound) {
			log.Printf("[WARN] VOlume (%s) is not found", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error retrieving volume: %v", err)
	}
	_ = d.Set("name", volume.Name)
	_ = d.Set("size", volume.Size)
	_ = d.Set("status", volume.Status)
	_ = d.Set("type", volume.VolumeType)
	_ = d.Set("category", volume.Category)
	_ = d.Set("created_at", volume.CreatedAt)
	_ = d.Set("availability_zone", volume.AvailabilityZone)
	_ = d.Set("user_id", volume.UserID)
	_ = d.Set("project_id", volume.ProjectID)
	return nil
}

func resourceBizflyCloudVolumeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	if d.HasChange("size") {
		// resize volume
		// TODO check state of extending task
		_, err := client.CloudServer.Volumes().ExtendVolume(context.Background(), d.Id(), d.Get("size").(int))
		if err != nil {
			if errors.Is(err, gobizfly.ErrNotFound) {
				log.Printf("[WARN] Volume is not found %s", d.Id())
				d.SetId("")
				return nil
			}
		}
	}
	return nil
}

func resourceBizflyCloudVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	err := client.CloudServer.Volumes().Delete(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("error delete volume: %v", err)
	}
	return nil
}
