package bizflycloud

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

func dataNetworkInterfaceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"network_id": {
			Type:     schema.TypeString,
			Optional: true,
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
		"mac_address": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"admin_state_up": {
			Type:     schema.TypeBool,
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
		"port_security_enabled": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"action": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"security_groups": {
			Type:     schema.TypeList,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
		},
		"fixed_ips": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: dataFixedIpsInfoSchema(),
			},
		},
		"created_at": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"updated_at": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func dataFixedIpInfoSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"subnet_id": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"ip_address": {
			Computed: true,
			Type:     schema.TypeString,
		},
	}
}

func dataFixedIpsInfoSchema() map[string]*schema.Schema {
	commonSchema := dataFixedIpInfoSchema()

	return commonSchema
}
