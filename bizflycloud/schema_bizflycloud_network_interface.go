package bizflycloud

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

func dataNetworkInterfaceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"network_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"attached_server": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"fixed_ip": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
}

func resourceNetworkInterfaceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"network_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"attached_server": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"fixed_ip": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
}
