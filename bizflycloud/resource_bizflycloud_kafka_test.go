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
	resource.AddTestSweepers("bizflycloud_kafka", &resource.Sweeper{
		Name: "bizflycloud_kafka",
	})
}

func TestAccBizflyCloudKafka_Basic(t *testing.T) {
	var cluster gobizfly.ClusterResponse
	rInt := acctest.RandInt()
	resourceName := "bizflycloud_kafka.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBizflyCloudKafkaDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBizflyCloudKafkaBasicConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBizflyCloudKafkaExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("tf-kafka-test-%d", rInt)),
					resource.TestCheckResourceAttr(resourceName, "nodes", "1"),
					resource.TestCheckResourceAttr(resourceName, "flavor", "2c_4g"),
					resource.TestCheckResourceAttr(resourceName, "volume_size", "10"),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", "HN1"),
					resource.TestCheckResourceAttr(resourceName, "public_access", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "version_id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
				),
			},
		},
	})
}

func TestAccBizflyCloudKafka_Update(t *testing.T) {
	var cluster gobizfly.ClusterResponse
	rInt := acctest.RandInt()
	resourceName := "bizflycloud_kafka.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBizflyCloudKafkaDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBizflyCloudKafkaBasicConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBizflyCloudKafkaExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "nodes", "1"),
					resource.TestCheckResourceAttr(resourceName, "flavor", "2c_4g"),
					resource.TestCheckResourceAttr(resourceName, "volume_size", "10"),
				),
			},
			{
				Config: testAccBizflyCloudKafkaUpdateConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBizflyCloudKafkaExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "nodes", "2"),
					resource.TestCheckResourceAttr(resourceName, "flavor", "4c_8g"),
					resource.TestCheckResourceAttr(resourceName, "volume_size", "20"),
				),
			},
		},
	})
}

func TestAccBizflyCloudKafka_ResizeFlavor(t *testing.T) {
	var cluster gobizfly.ClusterResponse
	rInt := acctest.RandInt()
	resourceName := "bizflycloud_kafka.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBizflyCloudKafkaDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBizflyCloudKafkaBasicConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBizflyCloudKafkaExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "flavor", "2c_4g"),
				),
			},
			{
				Config: testAccBizflyCloudKafkaResizeFlavorConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBizflyCloudKafkaExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "flavor", "4c_8g"),
				),
			},
		},
	})
}

func TestAccBizflyCloudKafka_ResizeVolume(t *testing.T) {
	var cluster gobizfly.ClusterResponse
	rInt := acctest.RandInt()
	resourceName := "bizflycloud_kafka.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBizflyCloudKafkaDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBizflyCloudKafkaBasicConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBizflyCloudKafkaExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "volume_size", "10"),
				),
			},
			{
				Config: testAccBizflyCloudKafkaResizeVolumeConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBizflyCloudKafkaExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "volume_size", "20"),
				),
			},
		},
	})
}

func TestAccBizflyCloudKafka_AddNode(t *testing.T) {
	var cluster gobizfly.ClusterResponse
	rInt := acctest.RandInt()
	resourceName := "bizflycloud_kafka.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBizflyCloudKafkaDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBizflyCloudKafkaBasicConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBizflyCloudKafkaExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "nodes", "1"),
				),
			},
			{
				Config: testAccBizflyCloudKafkaAddNodeConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBizflyCloudKafkaExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "nodes", "3"),
				),
			},
		},
	})
}

func testAccCheckBizflyCloudKafkaExists(n string, cluster *gobizfly.ClusterResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Kafka cluster ID is set")
		}

		client := testAccProvider.Meta().(*CombinedConfig).gobizflyClient()

		retrievedCluster, err := client.Kafka.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			return err
		}

		if retrievedCluster.ID != rs.Primary.ID {
			return fmt.Errorf("Kafka cluster not found")
		}

		*cluster = *retrievedCluster
		return nil
	}
}

func testAccCheckBizflyCloudKafkaDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*CombinedConfig).gobizflyClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bizflycloud_kafka" {
			continue
		}

		// Try to find the cluster
		_, err := client.Kafka.Get(context.Background(), rs.Primary.ID)

		if err != nil {
			if !errors.Is(err, gobizfly.ErrNotFound) {
				return fmt.Errorf(
					"Error waiting for Kafka cluster (%s) to be destroyed: %s",
					rs.Primary.ID, err)
			}
		}
	}

	return nil
}

func testAccBizflyCloudKafkaBasicConfig(rInt int) string {
	return fmt.Sprintf(`
resource "bizflycloud_kafka" "foobar" {
    name              = "tf-kafka-test-%d"
    version_id        = "10d9c0ee-5b3b-4a9e-9f21-4416fbb9e8ef"
    nodes             = 1
    flavor            = "2c_4g"
    volume_size       = 10
    availability_zone = "HN1"
    vpc_network_id    = ""
    public_access     = false
}
`, rInt)
}

func testAccBizflyCloudKafkaUpdateConfig(rInt int) string {
	return fmt.Sprintf(`
resource "bizflycloud_kafka" "foobar" {
    name              = "tf-kafka-test-%d"
    version_id        = "10d9c0ee-5b3b-4a9e-9f21-4416fbb9e8ef"
    nodes             = 2
    flavor            = "4c_8g"
    volume_size       = 20
    availability_zone = "HN1"
    vpc_network_id    = ""
    public_access     = false
}
`, rInt)
}

func testAccBizflyCloudKafkaResizeFlavorConfig(rInt int) string {
	return fmt.Sprintf(`
resource "bizflycloud_kafka" "foobar" {
    name              = "tf-kafka-test-%d"
    version_id        = "10d9c0ee-5b3b-4a9e-9f21-4416fbb9e8ef"
    nodes             = 1
    flavor            = "4c_8g"
    volume_size       = 10
    availability_zone = "HN1"
    vpc_network_id    = ""
    public_access     = false
}
`, rInt)
}

func testAccBizflyCloudKafkaResizeVolumeConfig(rInt int) string {
	return fmt.Sprintf(`
resource "bizflycloud_kafka" "foobar" {
    name              = "tf-kafka-test-%d"
    version_id        = "10d9c0ee-5b3b-4a9e-9f21-4416fbb9e8ef"
    nodes             = 1
    flavor            = "2c_4g"
    volume_size       = 20
    availability_zone = "HN1"
    vpc_network_id    = ""
    public_access     = false
}
`, rInt)
}

func testAccBizflyCloudKafkaAddNodeConfig(rInt int) string {
	return fmt.Sprintf(`
resource "bizflycloud_kafka" "foobar" {
    name              = "tf-kafka-test-%d"
    version_id        = "10d9c0ee-5b3b-4a9e-9f21-4416fbb9e8ef"
    nodes             = 3
    flavor            = "2c_4g"
    volume_size       = 10
    availability_zone = "HN1"
    vpc_network_id    = ""
    public_access     = false
}
`, rInt)
}
