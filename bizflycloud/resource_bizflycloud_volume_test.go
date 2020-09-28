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

func TestAccBizFlyCloudVolume_Basic(t *testing.T) {
	var volume gobizfly.Volume
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBizFlyCloudVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBizFlyCloudVolumeBasic_config(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBizFlyCloudVolumeExists("bizflycloud_volume.foobar", &volume),
					testAccCheckBizFlyCloudVolumeAttributes(&volume),
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

func testAccCheckBizFlyCloudVolumeExists(n string, volume *gobizfly.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No volume ID is set")
		}
		client := testAccProvider.Meta().(*CombinedConfig).gobizflyClient()

		retrieveVolume, err := client.Volume.Get(context.Background(), rs.Primary.ID)

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

func testAccCheckBizFlyCloudVolumeAttributes(volume *gobizfly.Volume) resource.TestCheckFunc {
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

func testAccCheckBizFlyCloudVolumeDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*CombinedConfig).gobizflyClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bizflycloud_volume" {
			continue
		}

		// Try to find the volume
		_, err := client.Volume.Get(context.Background(), rs.Primary.ID)

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

func testAccBizFlyCloudVolumeBasic_config(rInt int) string {
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
