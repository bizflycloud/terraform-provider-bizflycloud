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
	resource.AddTestSweepers("bizflycloud_cloud_database_instance", &resource.Sweeper{
		Name: "bizflycloud_cloud_database_instance",
	})
}

func TestAccBizFlyCloudCloudDatabaseinstance_Basic(t *testing.T) {
	var instance gobizfly.CloudDatabaseInstance
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBizFlyCloudCloudDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBizFlyCloudCloudDatabaseInstanceBasicConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBizFlyCloudCloudDatabaseInstanceExists("bizflycloud_cloud_database_instance.foobar", &instance),
					testAccCheckBizFlyCloudCloudDatabaseInstanceAttributes(&instance),
					resource.TestCheckResourceAttr(
						"bizflycloud_cloud_database_instance.foobar", "name", fmt.Sprintf("tf-testAccCloudDatabaseinstance-%d", rInt)),
					resource.TestCheckResourceAttr(
						"bizflycloud_cloud_database_instance.foobar", "type", "HDD"),
					resource.TestCheckResourceAttrSet("bizflycloud_cloud_database_instance.foobar", "status"),
					resource.TestCheckResourceAttrSet("bizflycloud_cloud_database_instance.foobar", "created_at"),
				),
			},
		},
	})
}

func testAccCheckBizFlyCloudCloudDatabaseInstanceExists(n string, instance *gobizfly.CloudDatabaseInstance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no cloud database instance ID is set")
		}
		client := testAccProvider.Meta().(*CombinedConfig).gobizflyClient()

		retrieveInstance, err := client.CloudDatabase.Instances().Get(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}
		if retrieveInstance.ID != rs.Primary.ID {
			return fmt.Errorf("cloud database instance not found")
		}
		*instance = *retrieveInstance
		return nil
	}
}

func testAccCheckBizFlyCloudCloudDatabaseInstanceAttributes(instance *gobizfly.CloudDatabaseInstance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if instance.Datastore.Type != "MongoDB" {
			return fmt.Errorf("bad cloud database instance datastore type: %s", instance.Datastore.Type)
		}

		if instance.Datastore.VersionID != "b48a59df-7a71-49a2-838c-50c9369976bc" {
			return fmt.Errorf("bad cloud database instance datastore version_id: %s", instance.Datastore.VersionID)
		}

		if instance.Volume.Size != 70 {
			return fmt.Errorf("bad cloud database instance volume size: %d", instance.Volume.Size)
		}

		return nil
	}
}

func testAccCheckBizFlyCloudCloudDatabaseInstanceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*CombinedConfig).gobizflyClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bizflycloud_cloud_database_instance" {
			continue
		}

		// Try to find the Instance
		_, err := client.CloudDatabase.Instances().Get(context.Background(), rs.Primary.ID)

		// Wait
		if err != nil {
			if !errors.Is(err, gobizfly.ErrNotFound) {
				return fmt.Errorf(
					"error waiting for cloud database instance (%s) to be destroyed: %s",
					rs.Primary.ID, err)
			}
		}
	}

	return nil
}

func testAccBizFlyCloudCloudDatabaseInstanceBasicConfig(rInt int) string {
	return fmt.Sprintf(`
		resource "bizflycloud_cloud_database_instance" "foobar" {
            name = "tf-testAccCloudDatabaseInstance-%d"
		    flavor_name                  = "1c_2g"
		    instance_type                = "enterprise"
		    volume_size                  = 70
		    datastore_type               = "MongoDB"
		    datastore_version_id         = "b48a59df-7a71-49a2-838c-50c9369976bc"
		    network_ids                  = ["d9000861-9958-461d-889e-04eab30f6dfc"]
		    public_access                = false
		    availability_zone            = "HN1"
		    autoscaling_enable           = false
		    autoscaling_volume_threshold = 90
		    autoscaling_volume_limited   = 90
		}
		`, rInt)
}
