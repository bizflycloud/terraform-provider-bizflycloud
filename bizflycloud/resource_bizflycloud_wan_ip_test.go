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
	resource.AddTestSweepers("bizflycloud_wan_ip", &resource.Sweeper{
		Name: "bizflycloud_wan_ip",
	})
}

func TestAccBizflyCloudWanIP(t *testing.T) {
	var wanIP gobizfly.WanIP
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBizflyCloudWanIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBizflyCloudWanIPConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBizflyCloudWanIPExists("bizflycloud_wan_ip.test_wan_1", &wanIP),
					testAccCheckBizflyCloudWanIPAttributes(&wanIP),
					resource.TestCheckResourceAttr(
						"bizflycloud_wan_ip.test_wan_1", "name", fmt.Sprintf("test-%d", rInt)),
				),
			},
		},
	})
}

func testAccCheckBizflyCloudWanIPExists(n string, wanIP *gobizfly.WanIP) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No WAN IP ID is set")
		}
		client := testAccProvider.Meta().(*CombinedConfig).gobizflyClient()

		retrieveWanIP, err := client.WanIP.Get(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}
		if retrieveWanIP.ID != rs.Primary.ID {
			return fmt.Errorf("WAN IP not found")
		}
		*wanIP = *retrieveWanIP
		return nil
	}
}

func testAccCheckBizflyCloudWanIPAttributes(wanIP *gobizfly.WanIP) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if wanIP.Name != "test_wan_ip" {
			return fmt.Errorf("Bad network interface name: %s", wanIP.Name)
		}
		if wanIP.AvailabilityZone != "HN1" {
			return fmt.Errorf("Bad availability zone: %s", wanIP.AvailabilityZone)
		}
		return nil
	}
}

func testAccCheckBizflyCloudWanIPDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*CombinedConfig).gobizflyClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bizflycloud_wan_ip" {
			continue
		}

		_, err := client.WanIP.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			if !errors.Is(err, gobizfly.ErrNotFound) {
				return fmt.Errorf("Error: %v", err)
			}
		}
	}
	return nil
}

func testAccBizflyCloudWanIPConfig(rInt int) string {
	return fmt.Sprintf(`
resource "bizflycloud_wan_ip" "test_wan_1" {
  name = "sapd-wan-ip-%d"
  availability_zone = "HN1"
}
`, rInt)
}
