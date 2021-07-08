package bizflycloud

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

func dataVirtualPrivateCloudNetworkSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"cidr": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"is_default": {
			Type:     schema.TypeBool,
			Optional: true,
		},
	}
}

func resourceVirtualPrivateCloudNetworkSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"cidr": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"is_default": {
			Type:     schema.TypeBool,
			Optional: true,
		},
	}
}
