package bizflycloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceWanIPSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"ip_address": {
			Type:     schema.TypeString,
			Required: true,
		},
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"ip_version": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"network_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"security_groups": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"billing_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"availability_zone": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"bandwidth": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	}
}

func resourceWanIPSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"network_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"tenant_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"device_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"security_groups": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"description": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"bandwidth": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"billing_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"availability_zone": {
			Type:     schema.TypeString,
			Required: true,
		},
		"ip_address": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"ip_version": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"attached_server": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
}
