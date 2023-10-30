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

const nodeID = "9ed6c0fb-205a-45fb-9d95-80d101affbbb"

func init() {
	resource.AddTestSweepers("bizflycloud_cloud_database_backup", &resource.Sweeper{
		Name: "bizflycloud_cloud_database_backup",
	})
}

func TestAccBizflyCloudCloudDatabaseBackup_Basic(t *testing.T) {
	var backup gobizfly.CloudDatabaseBackup
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBizflyCloudCloudDatabaseBackupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBizflyCloudCloudDatabaseBackupBasicConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBizflyCloudCloudDatabaseBackupExists("bizflycloud_cloud_database_backup.foobar", &backup),
					testAccCheckBizflyCloudCloudDatabaseBackupAttributes(&backup),
					resource.TestCheckResourceAttr(
						"bizflycloud_cloud_database_backup.foobar", "name", fmt.Sprintf("tf-testAccCloudDatabaseBackup-%d", rInt)),
					resource.TestCheckResourceAttr(
						"bizflycloud_cloud_database_backup.foobar", "type", "HDD"),
					resource.TestCheckResourceAttrSet("bizflycloud_cloud_database_backup.foobar", "status"),
					resource.TestCheckResourceAttrSet("bizflycloud_cloud_database_backup.foobar", "created_at"),
				),
			},
		},
	})
}

func testAccCheckBizflyCloudCloudDatabaseBackupExists(n string, backup *gobizfly.CloudDatabaseBackup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no cloud database backup ID is set")
		}
		client := testAccProvider.Meta().(*CombinedConfig).gobizflyClient()

		retrieveBackup, err := client.CloudDatabase.Backups().Get(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}
		if retrieveBackup.ID != rs.Primary.ID {
			return fmt.Errorf("cloud database backup not found")
		}
		*backup = *retrieveBackup
		return nil
	}
}

func testAccCheckBizflyCloudCloudDatabaseBackupAttributes(backup *gobizfly.CloudDatabaseBackup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if backup.NodeID != nodeID {
			return fmt.Errorf("bad cloud database backup source: %s", backup.NodeID)
		}

		return nil
	}
}

func testAccCheckBizflyCloudCloudDatabaseBackupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*CombinedConfig).gobizflyClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bizflycloud_cloud_database_backup" {
			continue
		}

		// Try to find the backup
		_, err := client.CloudDatabase.Backups().Get(context.Background(), rs.Primary.ID)

		// Wait
		if err != nil {
			if !errors.Is(err, gobizfly.ErrNotFound) {
				return fmt.Errorf(
					"error waiting for cloud database backup (%s) to be destroyed: %s",
					rs.Primary.ID, err)
			}
		}
	}

	return nil
}

func testAccBizflyCloudCloudDatabaseBackupBasicConfig(rInt int) string {
	return fmt.Sprintf(`
		resource "bizflycloud_cloud_database_backup" "foobar" {
            name = "tf-testAccCloudDatabaseBackup-%d"
            node_id = %s
		}
		`, rInt, nodeID)
}
