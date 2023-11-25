package bizflycloud

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

func resourceCDNSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"required": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"created_at": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"updated_at": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"deleted": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"tenant_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"ttl": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"active": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"nameserver": {
			Type:     schema.TypeList,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
		},
		"record_set": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: dataRecordSetInfoSchema(),
			},
		},
	}
}
