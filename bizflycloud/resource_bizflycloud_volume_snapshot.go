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

func resourceBizflyCloudVolumeSnapshot() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Create:        resourceBizflyCloudVolumeSnapshotCreate,
		Read:          resourceBizflyCloudVolumeSnapshotRead,
		Delete:        resourceBizflyCloudVolumeSnapshotDelete,
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"volume_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"size": {
				Type:     schema.TypeInt,
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
			"snapshot_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceBizflyCloudVolumeSnapshotCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	scr := gobizfly.SnapshotCreateRequest{
		Name:     d.Get("name").(string),
		VolumeId: d.Get("volume_id").(string),
		Force:    true,
	}
	snapshot, err := client.CloudServer.Snapshots().Create(context.Background(), &scr)
	if err != nil {
		return fmt.Errorf("Error creating snapshot: %v", err)
	}
	d.SetId(snapshot.Id)
	_ = d.Set("volume_id", snapshot.VolumeId)
	return resourceBizflyCloudVolumeSnapshotRead(d, meta)
}

func resourceBizflyCloudVolumeSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	snapshot, err := client.CloudServer.Snapshots().Get(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving snapshot %s: %v", d.Id(), err)
	}
	_ = d.Set("name", snapshot.Name)
	_ = d.Set("size", snapshot.Size)
	_ = d.Set("status", snapshot.Status)
	_ = d.Set("volume_id", snapshot.VolumeId)
	_ = d.Set("snapshot_type", snapshot.SnapshotType)
	_ = d.Set("type", snapshot.Type)
	_ = d.Set("availability_zone", snapshot.ZoneName)
	_ = d.Set("region_name", snapshot.RegionName)
	_ = d.Set("created_at", snapshot.CreateAt)
	_ = d.Set("updated_at", snapshot.UpdatedAt)
	return nil
}

func resourceBizflyCloudVolumeSnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	err := client.CloudServer.Snapshots().Delete(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error deleting volume snapshot %s: %v", d.Id(), err)
	}
	return nil
}
