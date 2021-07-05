package bizflycloud

import "github.com/hashicorp/terraform-plugin-sdk/helper/resource"

func init() {
	resource.AddTestSweepers("bizflycloud_vpc_network", &resource.Sweeper{
		Name: "bizflycloud_vpc_network",
	})
}
