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
	resource.AddTestSweepers("bizflycloud_volume", &resource.Sweeper{
		Name: "bizflycloud_volume",
	})
}

func TestAccBizflyCloudVolume_Basic(t *testing.T) {
	var volume gobizfly.Volume
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBizflyCloudVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBizflyCloudVolumeBasic_config(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBizflyCloudVolumeExists("bizflycloud_volume.foobar", &volume),
					testAccCheckBizflyCloudVolumeAttributes(&volume),
					resource.TestCheckResourceAttr(
						"bizflycloud_volume.foobar", "name", fmt.Sprintf("foo-%d", rInt)),
					resource.TestCheckResourceAttr(
						"bizflycloud_volume.foobar", "type", "HDD"),
					resource.TestCheckResourceAttrSet("bizflycloud_volume.foobar", "status"),
					resource.TestCheckResourceAttrSet("bizflycloud_volume.foobar", "created_at"),
				),
			},
		},
	})
}

func testAccCheckBizflyCloudVolumeExists(n string, volume *gobizfly.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No volume ID is set")
		}
		client := testAccProvider.Meta().(*CombinedConfig).gobizflyClient()

		retrieveVolume, err := client.CloudServer.Volumes().Get(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}
		if retrieveVolume.ID != rs.Primary.ID {
			return fmt.Errorf("Volume not found")
		}
		*volume = *retrieveVolume
		return nil
	}
}

func testAccCheckBizflyCloudVolumeAttributes(volume *gobizfly.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if volume.Size != 20 {
			return fmt.Errorf("Bad volume size: %d", volume.Size)
		}

		if volume.AvailabilityZone != "HN1" {
			return fmt.Errorf("Bad Availability zone name: %s", volume.AvailabilityZone)
		}

		if volume.VolumeType != "HDD" {
			return fmt.Errorf("Bad volume type: %s", volume.VolumeType)
		}

		return nil
	}
}

func testAccCheckBizflyCloudVolumeDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*CombinedConfig).gobizflyClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bizflycloud_volume" {
			continue
		}

		// Try to find the volume
		_, err := client.CloudServer.Volumes().Get(context.Background(), rs.Primary.ID)

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

func testAccBizflyCloudVolumeBasic_config(rInt int) string {
	return fmt.Sprintf(`
resource "bizflycloud_volume" "foobar" {
    name = "foo-%d"
    size = 20
    type = "HDD"
    category = "premium"
    availability_zone = "HN2"
}
`, rInt)
}
