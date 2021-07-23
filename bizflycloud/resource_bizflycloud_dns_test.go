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
	resource.AddTestSweepers("bizflycloud_dns", &resource.Sweeper{
		Name: "bizflycloud_dns",
	})
}

func TestAccBizFlyCloudDNS(t *testing.T) {
	var extendedZone gobizfly.ExtendedZone
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBizFlyCloudDNSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBizFlyCloudDNSConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBizFlyCloudDNSExists("bizflycloud_dns.zone", &extendedZone),
					testAccCheckBizFlyCloudDNSAttributes(&extendedZone),
					resource.TestCheckResourceAttr(
						"bizflycloud_dns.zone", "name", fmt.Sprintf("test.%d", rInt)),
				),
			},
		},
	})
}

func testAccCheckBizFlyCloudDNSExists(n string, extendedZone *gobizfly.ExtendedZone) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No volume ID is set")
		}
		client := testAccProvider.Meta().(*CombinedConfig).gobizflyClient()

		retrieveDNS, err := client.DNS.GetZone(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}
		if retrieveDNS.ID != rs.Primary.ID {
			return fmt.Errorf("DNS zone not found")
		}
		*extendedZone = *retrieveDNS
		return nil
	}
}

func testAccCheckBizFlyCloudDNSAttributes(extendedZone *gobizfly.ExtendedZone) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if extendedZone.Name != "test dns zone" {
			return fmt.Errorf("Bad dns zone description: %s", extendedZone.Name)
		}

		if extendedZone.Active != false {
			return fmt.Errorf("Bad dns zone is active: %t", extendedZone.Active)
		}

		return nil
	}
}

func testAccCheckBizFlyCloudDNSDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*CombinedConfig).gobizflyClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bizflycloud_dns" {
			continue
		}

		_, err := client.DNS.GetZone(context.Background(), rs.Primary.ID)
		if err != nil {
			if !errors.Is(err, gobizfly.ErrNotFound) {
				return fmt.Errorf("Error: %v", err)
			}
		}
	}
	return nil
}

func testAccBizFlyCloudDNSConfig(rInt int) string {
	return fmt.Sprintf(`
resource "bizflycloud_dns" "zone" {
    name = "test.%d"
}
`, rInt)
}
