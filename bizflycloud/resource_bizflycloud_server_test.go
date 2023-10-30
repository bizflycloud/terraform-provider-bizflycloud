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
	resource.AddTestSweepers("bizflycloud_server", &resource.Sweeper{
		Name: "bizflycloud_server",
	})
}

func TestAccBizflyCloudServer_Basic(t *testing.T) {
	var server gobizfly.Server
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBizflyCloudServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBizflyCloudServerBasic_config(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBizflyCloudServerExists("bizflycloud_server.foobar", &server),
					testAccCheckBizflyCloudServerAttributes(&server),
					resource.TestCheckResourceAttr(
						"bizflycloud_server.foobar", "name", fmt.Sprintf("foo-%d", rInt)),
					resource.TestCheckResourceAttr(
						"bizflycloud_server.foobar", "flavor_name", "4c_2g"),
					resource.TestCheckResourceAttrSet("bizflycloud_server.foobar", "status"),
					resource.TestCheckResourceAttrSet("bizflycloud_server.foobar", "created_at"),
				),
			},
			{
				ResourceName:            "bizflycloud_server.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{""},
			},
		},
	})
}

func testAccCheckBizflyCloudServerExists(n string, server *gobizfly.Server) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Server ID is set")
		}
		client := testAccProvider.Meta().(*CombinedConfig).gobizflyClient()

		retrieveServer, err := client.Server.Get(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}
		if retrieveServer.ID != rs.Primary.ID {
			return fmt.Errorf("Server not found")
		}
		*server = *retrieveServer
		return nil
	}
}

func testAccCheckBizflyCloudServerAttributes(server *gobizfly.Server) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if server.Flavor.Name != "4c_2g" {
			return fmt.Errorf("Bad flavor name: %s", server.Flavor.Name)
		}

		if server.AvailabilityZone != "HN1" {
			return fmt.Errorf("Bad Availability zone name: %s", server.AvailabilityZone)
		}

		if server.KeyName != "sapd1" {
			return fmt.Errorf("Bad ssh key name: %s", server.KeyName)
		}

		return nil
	}
}

func testAccCheckBizflyCloudServerDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*CombinedConfig).gobizflyClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bizflycloud_server" {
			continue
		}

		// Try to find the server
		_, err := client.Server.Get(context.Background(), rs.Primary.ID)

		// Wait

		if err != nil {
			if !errors.Is(err, gobizfly.ErrNotFound) {
				return fmt.Errorf(
					"Error waiting for server (%s) to be destroyed: %s",
					rs.Primary.ID, err)
			}
		}
	}

	return nil
}

func testAccBizflyCloudServerBasic_config(rInt int) string {
	return fmt.Sprintf(`
resource "bizflycloud_server" "foobar" {
	name = "foo-%d"
	flavor_name = "4c_2g"
    ssh_key = "sapd1"
    os_type = "image"
    os_id = "5f218529-ce32-4cb6-8557-920b16307d35"
    category = "premium"
    availability_zone = "HN1"
    root_disk_type = "HDD"
    root_disk_size = 20    
}
`, rInt)
}
