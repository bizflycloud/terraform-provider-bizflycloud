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
	"github.com/bizflycloud/gobizfly"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBizFlyCloudDatabaseBackup() *schema.Resource {
	return &schema.Resource{
		Create: resourceBizFlyCloudCloudDatabaseBackupCreate,
		Read:   resourceBizFlyCloudCloudDatabaseBackupRead,
		Delete: resourceBizFlyCloudCloudDatabaseBackupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(40 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},
		Schema: resourceCloudDatabaseBackupSchema(),
	}
}

func resourceBizFlyCloudCloudDatabaseBackupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()

	bac := &gobizfly.CloudDatabaseBackupCreate{
		BackupName: d.Get("name").(string),
		NodeID:     d.Get("node_id").(string),
		ParentID:   d.Get("parent_id").(string),
		InstanceID: d.Get("instance_id").(string),
	}

	var resourceType string
	var resourceID string

	if d.Get("instance_id").(string) != "" {
		resourceType = "instances"
		resourceID = d.Get("instance_id").(string)
	} else if d.Get("node_id").(string) != "" {
		resourceType = "nodes"
		resourceID = d.Get("node_id").(string)
	} else {
		return fmt.Errorf("[ERROR] create cloud database backup %s failed: not found resource to backup", d.Get("name"))
	}

	// retry
	retry := maxRetry
	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		backup, err := client.CloudDatabase.Backups().Create(context.Background(), resourceType, resourceID, bac)

		if err != nil {
			retry -= 1
			if retry > 0 {
				time.Sleep(timeSleep)
				return resource.RetryableError(fmt.Errorf("[ERROR] create cloud database backup %s failed: %s. Retrying", d.Get("name"), err))
			}

			return resource.NonRetryableError(fmt.Errorf("[ERROR] create cloud database backup %s failed: %s. Can't retry", d.Get("name"), err))
		}

		log.Printf("[DEBUG] creating cloud database backup %s", backup.Name)
		d.SetId(backup.ID)

		// wait for cloud database backup to become active
		_, err = waitForCloudDatabaseBackupCreate(d, meta)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("[ERROR] create cloud database backup (%s) failed: %s. Can't retry", d.Get("name").(string), err))
		}

		err = resourceBizFlyCloudCloudDatabaseBackupRead(d, meta)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("[ERROR] read cloud database backup (%s) failed: %s. Can't retry", d.Get("name").(string), err))
		}

		return nil
	})

}

func resourceBizFlyCloudCloudDatabaseBackupRead(d *schema.ResourceData, meta interface{}) error {
	if err := dataSourceBizFlyCloudDatabaseBackupRead(d, meta); err != nil {
		return err
	}
	return nil
}

func resourceBizFlyCloudCloudDatabaseBackupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	id := d.Id()

	// retry
	retry := maxRetry
	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		_, err := client.CloudDatabase.Backups().Delete(context.Background(), id)

		if err != nil {
			retry -= 1
			if retry > 0 {
				time.Sleep(timeSleep)
				return resource.RetryableError(fmt.Errorf("[ERROR] delete cloud database backup %s failed: %v. Retrying", id, err))
			}

			return resource.NonRetryableError(fmt.Errorf("[ERROR] delete cloud database backup %s failed: %v. Can't retry", id, err))
		}

		log.Printf("[DEBUG] delete cloud database backup %s success", id)

		_, err = waitForCloudDatabaseBackupDelete(d, meta)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("[ERROR] delete cloud database backup %s failed: %v. Can't retry", id, err))
		}
		return nil
	})
}

func waitForCloudDatabaseBackupCreate(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	log.Printf("[INFO] Waiting for cloud database instance (%s) to be ready", d.Get("name").(string))
	stateConf := &resource.StateChangeConf{
		Pending:        []string{"false", "BUILDING"},
		Target:         []string{"true", "COMPLETED"},
		Refresh:        newCloudDatabaseBackupStateRefreshFunc(d, meta),
		Timeout:        d.Timeout(schema.TimeoutCreate),
		Delay:          30 * time.Second,
		MinTimeout:     20 * time.Second,
		NotFoundChecks: 50,
	}
	return stateConf.WaitForState()
}

func waitForCloudDatabaseBackupDelete(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	log.Printf("[INFO] Waiting for cloud database instance (%s) to be ready", d.Get("name").(string))
	stateConf := &resource.StateChangeConf{
		Pending:        []string{"false"},
		Target:         []string{"true"},
		Refresh:        deleteCloudDatabaseBackupStateRefreshFunc(d, meta),
		Timeout:        d.Timeout(schema.TimeoutCreate),
		Delay:          30 * time.Second,
		MinTimeout:     20 * time.Second,
		NotFoundChecks: 50,
	}
	return stateConf.WaitForState()
}

func newCloudDatabaseBackupStateRefreshFunc(d *schema.ResourceData, meta interface{}) resource.StateRefreshFunc {
	client := meta.(*CombinedConfig).gobizflyClient()

	return func() (interface{}, string, error) {
		err := resourceBizFlyCloudCloudDatabaseBackupRead(d, meta)
		if err != nil {
			return nil, "", err
		}

		if attr, ok := d.GetOk("status"); ok {
			bac, err := client.CloudDatabase.Backups().Get(context.Background(), d.Id())
			if err != nil {
				return nil, "", fmt.Errorf("[ERROR] Retrieving cloud database backup %s error: %v", d.Id(), err)
			}
			switch attr := attr.(type) {
			case bool:
				return &bac, strconv.FormatBool(attr), nil
			default:
				return &bac, attr.(string), nil
			}
		}
		log.Print("ok")
		return nil, "", nil
	}
}

func deleteCloudDatabaseBackupStateRefreshFunc(d *schema.ResourceData, meta interface{}) resource.StateRefreshFunc {
	client := meta.(*CombinedConfig).gobizflyClient()

	return func() (interface{}, string, error) {
		bac, err := client.CloudDatabase.Backups().Get(context.Background(), d.Id())

		if errors.Is(err, gobizfly.ErrNotFound) {
			return bac, "true", nil
		}

		return bac, "", nil

	}
}
