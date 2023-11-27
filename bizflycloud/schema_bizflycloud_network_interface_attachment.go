package bizflycloud

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

func resourceNetworkInterfaceAttachmentSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"network_interface_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"server_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"firewall_ids": {
			Type:     schema.TypeSet,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
		},
	}
}
