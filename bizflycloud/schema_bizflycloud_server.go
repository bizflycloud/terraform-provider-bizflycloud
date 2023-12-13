package bizflycloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ServerNetworkInterface struct {
	ID string `json:"id"`
}

type ServerWANNetworkInterface struct {
	Version     string   `json:"version,omitempty"`
	FirewallIDs []string `json:"firewall_ids,omitempty"`
}

func resourceServerSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"flavor_name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"ssh_key": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"category": {
			Type:     schema.TypeString,
			Required: true,
		},
		"os_type": {
			Type:     schema.TypeString,
			Required: true,
		},
		"os_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"root_disk_volume_type": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"root_disk_size": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"availability_zone": {
			Type:     schema.TypeString,
			Required: true,
		},
		"project_id": {
			Type:     schema.TypeString,
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
		"user_data": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"network_plan": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "free_datatransfer",
		},
		"billing_plan": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "saving_plan",
		},
		"is_available": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"locked": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"default_public_ipv4": {
			Type: schema.TypeList,
			Elem: &schema.Resource{
				Schema: resourceServerFreeWANNetworkInterfaceSchema(),
			},
			Optional: true,
		},
		"default_public_ipv6": {
			Type: schema.TypeList,
			Elem: &schema.Resource{
				Schema: resourceServerFreeWANNetworkInterfaceSchema(),
			},
			Optional: true,
		},
		"network_interface_ids": {
			Type:     schema.TypeSet,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Computed: true,
		},
		"vpc_network_ids": {
			Type:     schema.TypeSet,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Computed: true,
		},
		"volume_ids": {
			Type:     schema.TypeSet,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Computed: true,
		},
	}
}

func resourceServerFreeWANNetworkInterfaceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"firewall_ids": {
			Type:     schema.TypeSet,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
		},
		"enabled": {
			Type:     schema.TypeBool,
			Default:  true,
			Optional: true,
		},
		"ip_address": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}
