package bizflycloud

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

func resourceBizFlyCloudVPCNetwork() *schema.Resource {
	return &schema.Resource{
		Create:        resourceBizFlyCloudVPCNetworkCreate,
		Update:        resourceBizFlyCloudVPCNetworkUpdate,
		Delete:        resourceBizFlyCloudVPCNetworkDelete,
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cidr": {
				Type:     schema.TypeString,
				Required: true,
			},
			"is-default": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceBizFlyCloudVPCNetworkCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	return nil
}

func resourceBizFlyCloudVPCNetworkUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceBizFlyCloudVPCNetworkDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
