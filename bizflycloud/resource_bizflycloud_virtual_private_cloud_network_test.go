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
	resource.AddTestSweepers("bizflycloud_virtual_private_cloud_network", &resource.Sweeper{
		Name: "bizflycloud_virtual_private_cloud_network",
	})
}

func TestAccBizFlyCloudVirtualPrivateCloud(t *testing.T) {
	var vpc gobizfly.VPC
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBizFlyCloudVirtualPrivateCloudNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBizFlyCloudVirtualPrivateCloudNetworkConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBizFlyCloudVirtualPrivateCloudNetworkExists("bizflycloud_virtual_private_cloud_network.abc", &vpc),
					testAccCheckBizFlyCloudVirtualPrivateCloudNetworkAttributes(&vpc),
					resource.TestCheckResourceAttr(
						"bizflycloud_virtual_private_cloud_network.abc", "name", fmt.Sprintf("foo-%d", rInt)),
				),
			},
		},
	})
}

func testAccCheckBizFlyCloudVirtualPrivateCloudNetworkExists(n string, vpc *gobizfly.VPC) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No volume ID is set")
		}
		client := testAccProvider.Meta().(*CombinedConfig).gobizflyClient()

		retrieveVirtualPrivateCloudNetwork, err := client.VPC.Get(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}
		if retrieveVirtualPrivateCloudNetwork.ID != rs.Primary.ID {
			return fmt.Errorf("Virtual private cloud network not found")
		}
		*vpc = *retrieveVirtualPrivateCloudNetwork
		return nil
	}
}

func testAccCheckBizFlyCloudVirtualPrivateCloudNetworkAttributes(vpc *gobizfly.VPC) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if vpc.Description != "test virtual private cloud network" {
			return fmt.Errorf("Bad virtual private cloud network description: %s", vpc.Description)
		}

		if vpc.IsDefault != false {
			return fmt.Errorf("Bad virtual private cloud network is default: %t", vpc.IsDefault)
		}

		return nil
	}
}

func testAccCheckBizFlyCloudVirtualPrivateCloudNetworkDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*CombinedConfig).gobizflyClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bizflycloud_virtual_private_cloud_network" {
			continue
		}

		_, err := client.VPC.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			if !errors.Is(err, gobizfly.ErrNotFound) {
				return fmt.Errorf("Error: %v", err)
			}
		}
	}
	return nil
}

func testAccBizFlyCloudVirtualPrivateCloudNetworkConfig(rInt int) string {
	return fmt.Sprintf(`
resource "bizflycloud_virtual_private_cloud_network" "abc" {
    name = "foo-%d"
    description = "test virtual private cloud network"
    is_default = false
}
`, rInt)
}
