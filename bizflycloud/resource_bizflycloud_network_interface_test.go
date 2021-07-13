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
	resource.AddTestSweepers("bizflycloud_network_interface", &resource.Sweeper{
		Name: "bizflycloud_network_interface",
	})
}

func TestAccBizFlyCloudNetworkInterface(t *testing.T) {
	var networkInterface gobizfly.NetworkInterface
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBizFlyCloudNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBizFlyCloudNetworkInterfaceConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBizFlyCloudNetworkInterfaceExists("bizflycloud_network_interface.abc", &networkInterface),
					testAccCheckBizFlyCloudNetworkInterfaceAttributes(&networkInterface),
					resource.TestCheckResourceAttr(
						"bizflycloud_network_interface.abc", "name", fmt.Sprintf("test-%d", rInt)),
				),
			},
		},
	})
}

func testAccCheckBizFlyCloudNetworkInterfaceExists(n string, networkInterface *gobizfly.NetworkInterface) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No volume ID is set")
		}
		client := testAccProvider.Meta().(*CombinedConfig).gobizflyClient()

		retrieveNetworkInterface, err := client.NetworkInterface.GetNetworkInterface(context.Background(), "fcda80f4-88ee-4708-a55f-3c6bcdf0585e", rs.Primary.ID)

		if err != nil {
			return err
		}
		if retrieveNetworkInterface.ID != rs.Primary.ID {
			return fmt.Errorf("network interface not found")
		}
		*networkInterface = *retrieveNetworkInterface
		return nil
	}
}

func testAccCheckBizFlyCloudNetworkInterfaceAttributes(networkInterface *gobizfly.NetworkInterface) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if networkInterface.Description != "test network interface" {
			return fmt.Errorf("Bad network interface description: %s", networkInterface.Description)
		}

		if networkInterface.AdminStateUp != false {
			return fmt.Errorf("Bad network interface is admin_state_up: %t", networkInterface.AdminStateUp)
		}

		return nil
	}
}

func testAccCheckBizFlyCloudNetworkInterfaceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*CombinedConfig).gobizflyClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bizflycloud_network_interface" {
			continue
		}

		_, err := client.NetworkInterface.GetNetworkInterface(context.Background(), "fcda80f4-88ee-4708-a55f-3c6bcdf0585e", rs.Primary.ID)
		if err != nil {
			if !errors.Is(err, gobizfly.ErrNotFound) {
				return fmt.Errorf("Error: %v", err)
			}
		}
	}
	return nil
}

func testAccBizFlyCloudNetworkInterfaceConfig(rInt int) string {
	return fmt.Sprintf(`
resource "bizflycloud_network_interface" "abc" {
    name = "test-%d"
    network_id = "fcda80f4-88ee-4708-a55f-3c6bcdf0585e"
}
`, rInt)
}
