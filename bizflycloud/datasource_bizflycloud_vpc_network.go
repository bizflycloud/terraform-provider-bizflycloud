package bizflycloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceBizFlyCloudVPCNetwork() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceBizFlyCloudVPCNetworkRead,
		Schema: dataVPCNetworkSchema(),
	}
}

func dataSourceBizFlyCloudVPCNetworkRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()

	if v, ok := d.GetOk("id"); ok {
		d.SetId(v.(string))
	}

	log.Printf("[DEBUG] Reading VPC Network: %s", d.Id())
	network, err := client.VPC.Get(context.Background(), d.Id())

	log.Printf("[DEBUG] Checking for error: %s", err)
	if err != nil {
		return fmt.Errorf("error describing VPC Network: %w", err)
	}

	log.Printf("[DEBUG] Found VPC Network: %s", d.Id())
	log.Printf("[DEBUG] bizflycloud_vpc_network found: %s", network.Name)

	d.SetId(network.ID)
	_ = d.Set("name", network.Name)
	_ = d.Set("description", network.Description)
	_ = d.Set("is_default", network.IsDefault)
	return nil
}
