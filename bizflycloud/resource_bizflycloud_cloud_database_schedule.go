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

func resourceBizFlyCloudDatabaseSchedule() *schema.Resource {
	return &schema.Resource{
		Create: resourceBizFlyCloudCloudDatabaseScheduleCreate,
		Read:   resourceBizFlyCloudCloudDatabaseScheduleRead,
		Delete: resourceBizFlyCloudCloudDatabaseScheduleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(40 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},
		Schema: resourceCloudDatabaseScheduleSchema(),
	}
}

func resourceBizFlyCloudCloudDatabaseScheduleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()

	scc := &gobizfly.CloudDatabaseScheduleCreate{
		DayOfMonth:   readIntArray(d.Get("day_of_month").(*schema.Set).List()),
		DayOfWeek:    readIntArray(d.Get("day_of_week").(*schema.Set).List()),
		Hour:         readIntArray(d.Get("hour").(*schema.Set).List()),
		LimitBackup:  d.Get("limit_backup").(int),
		Minute:       readIntArray(d.Get("minute").(*schema.Set).List()),
		ScheduleName: d.Get("name").(string),
		ScheduleType: d.Get("schedule_type").(string),
	}

	// retry
	retry := maxRetry
	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		schedule, err := client.CloudDatabase.Schedules().Create(context.Background(), d.Get("node_id").(string), scc)

		if err != nil {
			retry -= 1
			if retry > 0 {
				time.Sleep(timeSleep)
				return resource.RetryableError(fmt.Errorf("[ERROR] create cloud database schedule %s failed: %s. Retrying", d.Get("name"), err))
			}

			return resource.NonRetryableError(fmt.Errorf("[ERROR] create cloud database schedule %s failed: %s. Can't retry", d.Get("name"), err))
		}

		log.Printf("[DEBUG] creating cloud database schedule %s", schedule.Name)
		d.SetId(schedule.ID)

		err = resourceBizFlyCloudCloudDatabaseScheduleRead(d, meta)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("[ERROR] create cloud database schedule (%s) failed: %s. Can't retry", d.Get("name").(string), err))
		}

		return nil
	})
}

func resourceBizFlyCloudCloudDatabaseScheduleRead(d *schema.ResourceData, meta interface{}) error {
	if err := dataSourceBizFlyCloudDatabaseScheduleRead(d, meta); err != nil {
		return err
	}
	return nil
}

func resourceBizFlyCloudCloudDatabaseScheduleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	id := d.Id()

	// retry
	retry := maxRetry
	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		_, err := client.CloudDatabase.Schedules().Delete(context.Background(), id, &gobizfly.CloudDatabaseScheduleDelete{})

		if err != nil {
			retry -= 1
			if retry > 0 {
				time.Sleep(timeSleep)
				return resource.RetryableError(fmt.Errorf("[ERROR] delete cloud database schedule %s failed: %v. Retrying", id, err))
			}

			return resource.NonRetryableError(fmt.Errorf("[ERROR] delete cloud database schedule %s failed: %v. Can't retry", id, err))
		}
		return nil
	})
}

func readIntArray(items []interface{}) []int {
	intArray := make([]int, 0)
	for i := 0; i < len(items); i++ {
		item := items[i].(int)
		intArray = append(intArray, item)
	}
	return intArray
}
