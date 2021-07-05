package bizflycloud

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

func resourceBizFlyCloudVPCNetwork() *schema.Resource {
	return &schema.Resource{
		Create:        resourceBizFlyCloudVPCNetworkCreate,
		Update:        resourceBizFlyCloudVPCNetworkUpdate,
		Delete:        resourceBizFlyCloudVPCNetworkDelete,
		SchemaVersion: 1,
		Schema:        map[string]*schema.Schema{},
	}
}

func resourceBizFlyCloudVPCNetworkCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceBizFlyCloudVPCNetworkUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceBizFlyCloudVPCNetworkDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
