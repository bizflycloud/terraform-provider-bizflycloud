package bizflycloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKafkaSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"version_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"nodes": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"flavor": {
			Type:     schema.TypeString,
			Required: true,
		},
		"volume_size": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"vpc_network_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"availability_zone": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"public_access": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"created_at": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"obs_dashboard_url": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}
