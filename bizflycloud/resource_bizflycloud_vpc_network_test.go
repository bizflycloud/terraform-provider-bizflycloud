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

func TestAccBizflyCloudVPC(t *testing.T) {
	var vpc gobizfly.VPCNetwork
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBizflyCloudVPCNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBizflyCloudVPCNetworkConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBizflyCloudVPCNetworkExists("bizflycloud_vpc_network.abc", &vpc),
					testAccCheckBizflyCloudVPCNetworkAttributes(&vpc),
					resource.TestCheckResourceAttr(
						"bizflycloud_vpc_network.abc", "name", fmt.Sprintf("foo-%d", rInt)),
				),
			},
		},
	})
}

func testAccCheckBizflyCloudVPCNetworkExists(n string, vpc *gobizfly.VPCNetwork) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No volume ID is set")
		}
		client := testAccProvider.Meta().(*CombinedConfig).gobizflyClient()

		retrieveVPCNetwork, err := client.CloudServer.VPCNetworks().Get(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}
		if retrieveVPCNetwork.ID != rs.Primary.ID {
			return fmt.Errorf("vpc network not found")
		}
		*vpc = *retrieveVPCNetwork
		return nil
	}
}

func testAccCheckBizflyCloudVPCNetworkAttributes(vpc *gobizfly.VPCNetwork) resource.TestCheckFunc {
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

func testAccCheckBizflyCloudVPCNetworkDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*CombinedConfig).gobizflyClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bizflycloud_vpc_network" {
			continue
		}

		_, err := client.CloudServer.VPCNetworks().Get(context.Background(), rs.Primary.ID)
		if err != nil {
			if !errors.Is(err, gobizfly.ErrNotFound) {
				return fmt.Errorf("Error: %v", err)
			}
		}
	}
	return nil
}

func testAccBizflyCloudVPCNetworkConfig(rInt int) string {
	return fmt.Sprintf(`
resource "bizflycloud_vpc_network" "abc" {
    name = "foo-%d"
    description = "test vpc network"
    is_default = false
}
`, rInt)
}
