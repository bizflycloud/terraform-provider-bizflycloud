package bizflycloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

type ServerNetworkInterface struct {
	ID          string   `json:"id"`
	FirewallIDs []string `json:"firewall_ids,omitempty"`
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
		"password": {
			Type:     schema.TypeBool,
			Optional: true,
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
		"user_id": {
			Type:     schema.TypeString,
			Computed: true,
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
		"volume_ids": {
			Type:     schema.TypeSet,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
			Computed: true,
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
		"zone_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"locked": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"network_interface": {
			ConflictsWith: []string{"vpc_network_ids"},
			Type:          schema.TypeSet,
			ConfigMode:    schema.SchemaConfigModeAuto,
			Elem:          &schema.Resource{Schema: resourceServerNetworkInterfaceSchema()},
			Optional:      true,
		},
		"free_wan": {
			Type: schema.TypeSet,
			Elem: &schema.Resource{
				Schema: resourceServerFreeWANNetworkInterfaceSchema(),
			},
			Optional: true,
			DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				return false
			},
		},
		"vpc_network_ids": {
			ConflictsWith: []string{"network_interface"},
			Type:          schema.TypeSet,
			ConfigMode:    schema.SchemaConfigModeAttr,
			Elem:          &schema.Schema{Type: schema.TypeString},
			Optional:      true,
		},
	}
}

func resourceServerNetworkInterfaceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
		"firewall_ids": {
			Type:       schema.TypeSet,
			ConfigMode: schema.SchemaConfigModeAttr,
			Elem:       &schema.Schema{Type: schema.TypeString},
			Optional:   true,
		},
	}
}

func resourceServerFreeWANNetworkInterfaceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"ip_version": {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntInSlice([]int{4, 6}),
		},
		"firewall_ids": {
			Type:       schema.TypeSet,
			ConfigMode: schema.SchemaConfigModeAttr,
			Elem:       &schema.Schema{Type: schema.TypeString},
			Optional:   true,
		},
	}
}
