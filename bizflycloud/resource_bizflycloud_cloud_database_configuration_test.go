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
	resource.AddTestSweepers("bizflycloud_cloud_database_configuration", &resource.Sweeper{
		Name: "bizflycloud_cloud_database_configuration",
	})
}

func TestAccBizFlyCloudCloudDatabaseConfiguration_Basic(t *testing.T) {
	var configuration gobizfly.CloudDatabaseConfiguration
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBizFlyCloudCloudDatabaseConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBizFlyCloudCloudDatabaseConfigurationBasicConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBizFlyCloudCloudDatabaseConfigurationExists("bizflycloud_cloud_database_configuration.foobar", &configuration),
					testAccCheckBizFlyCloudCloudDatabaseConfigurationAttributes(&configuration),
					resource.TestCheckResourceAttr(
						"bizflycloud_cloud_database_configuration.foobar", "name", fmt.Sprintf("tf-testAccCloudDatabaseconfiguration-%d", rInt)),
					resource.TestCheckResourceAttr(
						"bizflycloud_cloud_database_configuration.foobar", "type", "HDD"),
					resource.TestCheckResourceAttrSet("bizflycloud_cloud_database_configuration.foobar", "status"),
					resource.TestCheckResourceAttrSet("bizflycloud_cloud_database_configuration.foobar", "created_at"),
				),
			},
		},
	})
}

func testAccCheckBizFlyCloudCloudDatabaseConfigurationExists(n string, configuration *gobizfly.CloudDatabaseConfiguration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no cloud database configuration ID is set")
		}
		client := testAccProvider.Meta().(*CombinedConfig).gobizflyClient()

		retrieveConfiguration, err := client.CloudDatabase.Configurations().Get(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}
		if retrieveConfiguration.ID != rs.Primary.ID {
			return fmt.Errorf("cloud database configuration not found")
		}
		*configuration = *retrieveConfiguration
		return nil
	}
}

func testAccCheckBizFlyCloudCloudDatabaseConfigurationAttributes(configuration *gobizfly.CloudDatabaseConfiguration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if configuration.Datastore.Type != "MongoDB" {
			return fmt.Errorf("bad cloud database configuration datastore type: %s", configuration.Datastore.Type)
		}

		if configuration.Datastore.Name != "4.4.7" {
			return fmt.Errorf("bad cloud database configuration datastore name: %s", configuration.Datastore.Name)
		}

		return nil
	}
}

func testAccCheckBizFlyCloudCloudDatabaseConfigurationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*CombinedConfig).gobizflyClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bizflycloud_cloud_database_configuration" {
			continue
		}

		// Try to find the Configuration
		_, err := client.CloudDatabase.Configurations().Get(context.Background(), rs.Primary.ID)

		// Wait
		if err != nil {
			if !errors.Is(err, gobizfly.ErrNotFound) {
				return fmt.Errorf(
					"error waiting for cloud database configuration (%s) to be destroyed: %s",
					rs.Primary.ID, err)
			}
		}
	}

	return nil
}

func testAccBizFlyCloudCloudDatabaseConfigurationBasicConfig(rInt int) string {
	return fmt.Sprintf(`
		resource "bizflycloud_cloud_database_configuration" "foobar" {
            name = "tf-testAccCloudDatabaseConfiguration-%d"
            datastore_type = "MongoDB",
			datastore_version_name = "4.4.7",
			parameters = {
			    "auditLog.format" = "test"
			    "net.ipv6" = false,
			    "net.maxIncomingConnections" = 123
			}
		}
		`, rInt)
}
