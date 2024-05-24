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
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

func init() {
	resource.AddTestSweepers("bizflycloud_kubernetes", &resource.Sweeper{
		Name: "bizflycloud_kubernetes",
	})
}

func TestAccBizflyCloudCluster(t *testing.T) {
	var cluster gobizfly.FullCluster
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBizflyCloudClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBizflyCloudClusterConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBizflyCloudClusterExists("bizflycloud_kubernetes.xyz", &cluster),
					testAccCheckBizflyCloudClusterAttributes(&cluster),
					resource.TestCheckResourceAttr(
						"bizflycloud_kubernetes.xyz", "name", fmt.Sprintf("foo-%d", rInt)),
				),
			},
		},
	})
}
func testAccCheckBizflyCloudClusterExists(n string, cluster *gobizfly.FullCluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No cluster ID is set: %s", rs.Primary.ID)
		}
		client := testAccProvider.Meta().(*CombinedConfig).gobizflyClient()
		retrieveCluster, err := client.KubernetesEngine.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			return err
		}
		if cluster.UID != rs.Primary.ID {
			return fmt.Errorf("Cluster not fount")
		}
		*cluster = *retrieveCluster
		return nil
	}
}

func testAccCheckBizflyCloudClusterAttributes(cluster *gobizfly.FullCluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if !cluster.AutoUpgrade {
			return fmt.Errorf("Bad autoupgrade %v", cluster.AutoUpgrade)
		}
		if len(cluster.Tags) != 1 {
			return fmt.Errorf("Bad tags %v", cluster.Tags)
		}
		return nil
	}
}

func testAccCheckBizflyCloudClusterDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*CombinedConfig).gobizflyClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bizflycloud_kubernetes" {
			continue
		}

		_, err := client.CloudServer.Volumes().Get(context.Background(), rs.Primary.ID)
		if err != nil {
			if !errors.Is(err, gobizfly.ErrNotFound) {
				return fmt.Errorf("Error: %v", err)
			}
		}
	}
	return nil
}

func testAccBizflyCloudClusterConfig(rInt int) string {
	return fmt.Sprintf(`
resource "bizflycloud_kubernetes" "xyz" {
	name = "foo-%d"
	version = "5f6425f3d0d3befd40e7a31f"
	auto_upgrade = false
	enable_cloud = true
	tags = ["string"]
	worker_pools = [
	{
		name = "pool"
		version = "v1.18.0"
		flavor = "8c_8g"
		profile_type = "premium"
		volume_type = "SSD"
		volume_size = 40,
		availability_zone = "HN1"
		desired_size = 1
		enable_autoscaling = true
		min_size = 1
		max_size = 3
		tags = ["string"]
}
]
`, rInt)
}
