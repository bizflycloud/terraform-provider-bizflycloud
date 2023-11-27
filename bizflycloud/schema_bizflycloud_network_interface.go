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
			Computed: true,
		},
		"attached_server": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"ip_address": {
			Type:     schema.TypeString,
			Required: true,
		},
		"security_groups": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
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
		"status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"fixed_ips": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: dataFixedIpsInfoSchema(),
			},
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
			Computed: true,
		},
		"server_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"fixed_ip": {
			Type:     schema.TypeString,
			Computed: true,
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
		"security_groups": {
			Type:     schema.TypeList,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Computed: true,
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
