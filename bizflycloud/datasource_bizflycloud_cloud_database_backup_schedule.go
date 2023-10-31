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

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceBizflyCloudDatabaseBackupSchedule() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceBizflyCloudDatabaseBackupScheduleRead,
		Schema: dataCloudDatabaseBackupScheduleSchema(),
	}
}

func dataSourceBizflyCloudDatabaseBackupScheduleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()

	if v, ok := d.GetOk("id"); ok {
		d.SetId(v.(string))
	}

	scheduleID := d.Id()

	log.Printf("[DEBUG] Reading database schedule: %s", scheduleID)
	schedule, err := client.CloudDatabase.BackupSchedules().Get(context.Background(), scheduleID)

	log.Printf("[DEBUG] Checking for error: %s", err)
	if err != nil {
		return fmt.Errorf("error describing database schedule: %w", err)
	}

	log.Printf("[DEBUG] Found database backup schedule: %s", scheduleID)
	log.Printf("[DEBUG] bizflycloud_cloud_database_backup_schedule - Single database backup schedule found: %s", schedule.Name)

	d.SetId(schedule.ID)
	_ = d.Set("cron_expression", schedule.CronExpression)
	_ = d.Set("first_execution_time", schedule.FirstExecutionTime)
	_ = d.Set("instance_id", schedule.InstanceID)
	_ = d.Set("instance_name", schedule.InstanceName)
	_ = d.Set("limit_backup", schedule.LimitBackup)
	_ = d.Set("name", schedule.Name)
	_ = d.Set("next_execution_time", schedule.NextExecutionTime)
	_ = d.Set("node_id", schedule.NodeID)
	_ = d.Set("node_name", schedule.NodeName)
	return nil
}
