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
	"testing"

	"github.com/bizflycloud/gobizfly"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func init() {
	resource.AddTestSweepers("bizflycloud_cloud_database_schedule", &resource.Sweeper{
		Name: "bizflycloud_cloud_database_schedule",
	})
}

func TestAccBizFlyCloudCloudDatabaseschedule_Basic(t *testing.T) {
	var schedule gobizfly.CloudDatabaseSchedule
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBizFlyCloudCloudDatabaseScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBizFlyCloudCloudDatabaseScheduleBasicConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBizFlyCloudCloudDatabaseScheduleExists("bizflycloud_cloud_database_schedule.foobar", &schedule),
					testAccCheckBizFlyCloudCloudDatabaseScheduleAttributes(&schedule),
					resource.TestCheckResourceAttr(
						"bizflycloud_cloud_database_schedule.foobar", "name", fmt.Sprintf("tf-testAccCloudDatabaseschedule-%d", rInt)),
					resource.TestCheckResourceAttr(
						"bizflycloud_cloud_database_schedule.foobar", "type", "HDD"),
					resource.TestCheckResourceAttrSet("bizflycloud_cloud_database_schedule.foobar", "status"),
					resource.TestCheckResourceAttrSet("bizflycloud_cloud_database_schedule.foobar", "created_at"),
				),
			},
		},
	})
}

func testAccCheckBizFlyCloudCloudDatabaseScheduleExists(n string, schedule *gobizfly.CloudDatabaseSchedule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no cloud database schedule ID is set")
		}
		client := testAccProvider.Meta().(*CombinedConfig).gobizflyClient()

		retrieveSchedule, err := client.CloudDatabase.Schedules().Get(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}
		if retrieveSchedule.ID != rs.Primary.ID {
			return fmt.Errorf("cloud database schedule not found")
		}
		*schedule = *retrieveSchedule
		return nil
	}
}

func testAccCheckBizFlyCloudCloudDatabaseScheduleAttributes(schedule *gobizfly.CloudDatabaseSchedule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if schedule.NodeID != "9ed6c0fb-205a-45fb-9d95-80d101affbbb" {
			return fmt.Errorf("bad cloud database schedule node source: %s", schedule.NodeID)
		}

		if schedule.LimitBackup != 1 {
			return fmt.Errorf("bad cloud database schedule limit backup: %d", schedule.LimitBackup)
		}

		return nil
	}
}

func testAccCheckBizFlyCloudCloudDatabaseScheduleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*CombinedConfig).gobizflyClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bizflycloud_cloud_database_schedule" {
			continue
		}

		// Try to find the Schedule
		_, err := client.CloudDatabase.Schedules().Get(context.Background(), rs.Primary.ID)

		// Wait
		if err != nil {
			if !errors.Is(err, gobizfly.ErrNotFound) {
				return fmt.Errorf(
					"error waiting for cloud database schedule (%s) to be destroyed: %s",
					rs.Primary.ID, err)
			}
		}
	}

	return nil
}

func testAccBizFlyCloudCloudDatabaseScheduleBasicConfig(rInt int) string {
	return fmt.Sprintf(`
		resource "bizflycloud_cloud_database_schedule" "foobar" {
            name = "tf-testAccCloudDatabaseSchedule-%d"
			limit_backup = 1
			schedule_type = "monthly"
			minute = [20, 50]
			hour = [7]
			day_of_month = [10]
			node_id = "9ed6c0fb-205a-45fb-9d95-80d101affbbb"
		}
		`, rInt)
}
