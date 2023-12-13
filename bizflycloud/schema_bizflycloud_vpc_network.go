package bizflycloud

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

func dataVPCNetworkSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cidr": {
			Type:     schema.TypeString,
			Required: true,
		},
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"subnets": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: dataSubnetsInfoSchema(),
			},
		},
		"is_default": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"status": {
			Type:     schema.TypeString,
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
		"availability_zones": {
			Type:     schema.TypeList,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Computed: true,
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
			Computed: true,
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
		"allocation_pools": {
			Computed: true,
			Type:     schema.TypeList,
			Elem: &schema.Resource{
				Schema: dataAllocationPoolsInfoSchema(),
			},
		},
	}
}

func dataSubnetsInfoSchema() map[string]*schema.Schema {
	commonSchema := dataSubnetInfoSchema()

	return commonSchema
}

func dataAllocationPoolsInfoSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"end": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"start": {
			Computed: true,
			Type:     schema.TypeString,
		},
	}
}
