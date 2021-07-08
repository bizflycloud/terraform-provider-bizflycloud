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
	resource.AddTestSweepers("bizflycloud_vpc_network", &resource.Sweeper{
		Name: "bizflycloud_vpc_network",
	})
}

func TestAccBizFlyCloudVPC(t *testing.T) {
	var vpc gobizfly.VPC
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBizFlyCloudVpcNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBizFlyCloudVpcNetworkConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBizFlyCloudVpcNetworkExists("bizflycloud_vpc_network.abc", &vpc),
					testAccCheckBizFlyCloudVpcNetworkAttributes(&vpc),
					resource.TestCheckResourceAttr(
						"bizflycloud_vpc_network.abc", "name", fmt.Sprintf("foo-%d", rInt)),
				),
			},
		},
	})
}

func testAccCheckBizFlyCloudVpcNetworkExists(n string, vpc *gobizfly.VPC) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No volume ID is set")
		}
		client := testAccProvider.Meta().(*CombinedConfig).gobizflyClient()

		retrieveVpcNetwork, err := client.VPC.Get(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}
		if retrieveVpcNetwork.ID != rs.Primary.ID {
			return fmt.Errorf("VPC Network not found")
		}
		*vpc = *retrieveVpcNetwork
		return nil
	}
}

func testAccCheckBizFlyCloudVpcNetworkAttributes(vpc *gobizfly.VPC) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if vpc.Description != "test vpc network" {
			return fmt.Errorf("Bad vpc network description: %s", vpc.Description)
		}

		if vpc.IsDefault != false {
			return fmt.Errorf("Bad vpc network is default: %t", vpc.IsDefault)
		}

		return nil
	}
}

func testAccCheckBizFlyCloudVpcNetworkDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*CombinedConfig).gobizflyClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bizflycloud_vpc_network" {
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

func testAccBizFlyCloudVpcNetworkConfig(rInt int) string {
	return fmt.Sprintf(`
resource "bizflycloud_vpc_network" "abc" {
    name = "foo-%d"
    description = "test vpc network"
    is_default = false
}
`, rInt)
}
