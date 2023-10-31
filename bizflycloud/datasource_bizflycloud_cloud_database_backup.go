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
	"log"
)

func datasourceBizflyCloudDatabaseBackup() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceBizflyCloudDatabaseBackupRead,
		Schema: dataCloudDatabaseBackupSchema(),
	}
}

func dataSourceBizflyCloudDatabaseBackupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()

	if v, ok := d.GetOk("id"); ok {
		d.SetId(v.(string))
	}

	backupID := d.Id()

	log.Printf("[DEBUG] Reading database backup: %s", backupID)
	backup, err := client.CloudDatabase.Backups().Get(context.Background(), backupID)

	log.Printf("[DEBUG] Checking for error: %s", err)
	if err != nil {
		return fmt.Errorf("error describing database backup: %w", err)
	}

	log.Printf("[DEBUG] Found database backup: %s", backupID)
	log.Printf("[DEBUG] bizflycloud_cloud_database_backup - Single database backup found: %s", backup.Name)

	d.SetId(backup.ID)
	_ = d.Set("created", backup.Created)
	_ = d.Set("description", backup.Description)
	_ = d.Set("name", backup.Name)
	_ = d.Set("node_id", backup.NodeID)
	_ = d.Set("parent_id", backup.ParentID)
	_ = d.Set("size", backup.Size)
	_ = d.Set("status", backup.Status)
	_ = d.Set("type", backup.Type)
	_ = d.Set("updated", backup.Updated)

	if err := d.Set("datastore", FlattenStruct(backup.Datastore)); err != nil {
		return fmt.Errorf("error setting datastore for backup %s: %s", d.Id(), err)
	}
	return nil
}
