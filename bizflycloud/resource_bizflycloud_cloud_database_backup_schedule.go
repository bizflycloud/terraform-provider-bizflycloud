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
	"log"
	"time"

	"github.com/bizflycloud/gobizfly"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBizflyCloudDatabaseBackupSchedule() *schema.Resource {
	return &schema.Resource{
		Create: resourceBizflyCloudCloudDatabaseBackupScheduleCreate,
		Delete: resourceBizflyCloudCloudDatabaseBackupScheduleDelete,
		Read:   resourceBizflyCloudCloudDatabaseBackupScheduleRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(40 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},
		Schema: resourceCloudDatabaseBackupScheduleSchema(),
	}
}

func resourceBizflyCloudCloudDatabaseBackupScheduleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()

	scc := &gobizfly.CloudDatabaseBackupScheduleCreate{
		CronExpression: d.Get("cron_expression").(string),
		LimitBackup:    d.Get("limit_backup").(int),
		Name:           d.Get("name").(string),
	}

	// retry
	retry := maxRetry
	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		schedule, err := client.CloudDatabase.BackupSchedules().Create(context.Background(), d.Get("node_id").(string), scc)

		if err != nil {
			retry--
			if retry > 0 {
				time.Sleep(timeSleep)
				return resource.RetryableError(fmt.Errorf("[ERROR] create cloud database schedule %s failed: %s. Retrying", d.Get("name"), err))
			}

			return resource.NonRetryableError(fmt.Errorf("[ERROR] create cloud database schedule %s failed: %s. Can't retry", d.Get("name"), err))
		}

		log.Printf("[DEBUG] creating cloud database schedule %s", schedule.Name)
		d.SetId(schedule.ID)

		err = resourceBizflyCloudCloudDatabaseBackupScheduleRead(d, meta)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("[ERROR] create cloud database schedule (%s) failed: %s. Can't retry", d.Get("name").(string), err))
		}

		return nil
	})
}

func resourceBizflyCloudCloudDatabaseBackupScheduleRead(d *schema.ResourceData, meta interface{}) error {
	if err := dataSourceBizflyCloudDatabaseBackupScheduleRead(d, meta); err != nil {
		return err
	}
	return nil
}

func resourceBizflyCloudCloudDatabaseBackupScheduleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	id := d.Id()

	// retry
	retry := maxRetry
	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		_, err := client.CloudDatabase.BackupSchedules().Delete(context.Background(), id, &gobizfly.CloudDatabaseBackupScheduleDelete{})

		if err != nil {
			retry--
			if retry > 0 {
				time.Sleep(timeSleep)
				return resource.RetryableError(fmt.Errorf("[ERROR] delete cloud database schedule %s failed: %v. Retrying", id, err))
			}

			return resource.NonRetryableError(fmt.Errorf("[ERROR] delete cloud database schedule %s failed: %v. Can't retry", id, err))
		}
		return nil
	})
}
