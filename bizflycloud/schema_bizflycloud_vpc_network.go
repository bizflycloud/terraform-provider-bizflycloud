package bizflycloud

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

func dataVPCNetworkSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Computed: true,
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

func resourceVPCNetworkSchema() map[string]*schema.Schema {
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
		"availability_zones": {
			Type:     schema.TypeList,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Computed: true,
		},
		"status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"tags": {
			Type:     schema.TypeList,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Computed: true,
		},
		"created_at": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"updated_at": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"mtu": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"subnets": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: dataSubnetsInfoSchema(),
			},
		},
	}
}

func dataSubnetInfoSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"project_id": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"ip_version": {
			Computed: true,
			Type:     schema.TypeInt,
		},
		"gateway_ip": {
			Computed: true,
			Type:     schema.TypeString,
		},
	}
}

func dataSubnetsInfoSchema() map[string]*schema.Schema {
	commonSchema := dataSubnetInfoSchema()

	return commonSchema
}
