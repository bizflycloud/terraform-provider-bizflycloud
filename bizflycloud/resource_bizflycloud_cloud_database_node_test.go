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
	resource.AddTestSweepers("bizflycloud_cloud_database_node", &resource.Sweeper{
		Name: "bizflycloud_cloud_database_node",
	})
}

func TestAccBizFlyCloudCloudDatabasenode_Basic(t *testing.T) {
	var node gobizfly.CloudDatabaseNode
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBizFlyCloudCloudDatabaseNodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBizFlyCloudCloudDatabaseNodeBasicConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBizFlyCloudCloudDatabaseNodeExists("bizflycloud_cloud_database_node.foobar", &node),
					testAccCheckBizFlyCloudCloudDatabaseNodeAttributes(&node),
					resource.TestCheckResourceAttr(
						"bizflycloud_cloud_database_node.foobar", "name", fmt.Sprintf("tf-testAccCloudDatabasenode-%d", rInt)),
					resource.TestCheckResourceAttr(
						"bizflycloud_cloud_database_node.foobar", "type", "HDD"),
					resource.TestCheckResourceAttrSet("bizflycloud_cloud_database_node.foobar", "status"),
					resource.TestCheckResourceAttrSet("bizflycloud_cloud_database_node.foobar", "created_at"),
				),
			},
		},
	})
}

func testAccCheckBizFlyCloudCloudDatabaseNodeExists(n string, node *gobizfly.CloudDatabaseNode) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no cloud database node ID is set")
		}
		client := testAccProvider.Meta().(*CombinedConfig).gobizflyClient()

		retrieveNode, err := client.CloudDatabase.Nodes().Get(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}
		if retrieveNode.ID != rs.Primary.ID {
			return fmt.Errorf("cloud database node not found")
		}
		*node = *retrieveNode
		return nil
	}
}

func testAccCheckBizFlyCloudCloudDatabaseNodeAttributes(node *gobizfly.CloudDatabaseNode) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if node.Role != "secondary" {
			return fmt.Errorf("bad cloud database node role: %s", node.Role)
		}

		return nil
	}
}

func testAccCheckBizFlyCloudCloudDatabaseNodeDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*CombinedConfig).gobizflyClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bizflycloud_cloud_database_node" {
			continue
		}

		// Try to find the Node
		_, err := client.CloudDatabase.Nodes().Get(context.Background(), rs.Primary.ID)

		// Wait
		if err != nil {
			if !errors.Is(err, gobizfly.ErrNotFound) {
				return fmt.Errorf(
					"error waiting for cloud database node (%s) to be destroyed: %s",
					rs.Primary.ID, err)
			}
		}
	}

	return nil
}

func testAccBizFlyCloudCloudDatabaseNodeBasicConfig(rInt int) string {
	return fmt.Sprintf(`
		resource "bizflycloud_cloud_database_node" "foobar" {
            name = "tf-testAccCloudDatabaseNode-%d"
		    replica_of = "4ac68fb6-c623-4c30-aca7-fdca40f8d6a7"
			role = "secondary"
		}
		`, rInt)
}
