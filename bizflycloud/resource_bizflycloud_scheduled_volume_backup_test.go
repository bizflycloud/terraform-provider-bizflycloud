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
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

func init() {
	resource.AddTestSweepers("bizflycloud_scheduled_volume_backup", &resource.Sweeper{
		Name: "bizflycloud_scheduled_volume_backup",
	})
}

func TestAccBizflyCloudScheduledVolumeBackup(t *testing.T) {
	var backup gobizfly.ExtendedBackup
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBizflyCloudScheduledVolumeBackupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBizflyCloudScheduledVolumeBackupConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBizflyCloudScheduledVolumeBackupExists("bizflycloud_scheduled_volume_backup.test", &backup),
					testAccCheckBizflyCloudScheduledVolumeBackupAttributes(&backup),
					resource.TestCheckResourceAttr("bizflycloud_scheduled_volume_backup.test", "name", "test"),
				),
			},
		},
	})
}

func testAccCheckBizflyCloudScheduledVolumeBackupExists(n string, backup *gobizfly.ExtendedBackup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Backup ID is set")
		}
		client := testAccProvider.Meta().(*CombinedConfig).gobizflyClient()
		retrieveBackup, err := client.ScheduledVolumeBackup.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			return err
		}
		if retrieveBackup.ID != rs.Primary.ID {
			return fmt.Errorf("Backup not found")
		}
		*backup = *retrieveBackup
		return nil
	}
}

func testAccCheckBizflyCloudScheduledVolumeBackupAttributes(backup *gobizfly.ExtendedBackup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if backup.Options.Size != "10" {
			return fmt.Errorf("Bad size: %s", backup.Options.Size)
		}
		if backup.ResourceID != "vol-0a2b3c4d" {
			return fmt.Errorf("Bad volume id: %s", backup.ResourceID)
		}
		return nil
	}
}

func testAccCheckBizflyCloudScheduledVolumeBackupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*CombinedConfig).gobizflyClient()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bizflycloud_scheduled_volume_backup" {
			continue
		}
		_, err := client.ScheduledVolumeBackup.Get(context.Background(), rs.Primary.ID)
		if err == nil {
			if !errors.Is(err, gobizfly.ErrNotFound) {
				return fmt.Errorf("Error: %v", err)
			}
		}
	}
	return nil
}

func testAccBizflyCloudScheduledVolumeBackupConfig() string {
	return `resource "bizflycloud_scheduled_volume_backup" "backup_test" {
  volume_id = "11a2e71b-8701-47a0-b247-41843db17e54"
  frequency = "2880"
  size = "4"
  scheduled_hour = 4
}`
}
