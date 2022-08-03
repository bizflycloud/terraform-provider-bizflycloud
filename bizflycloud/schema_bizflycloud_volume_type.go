package bizflycloud

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

func dataSourceVolumeTypeSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"category": {
			Type:     schema.TypeString,
			Required: true,
		},
		"type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"availability_zones": {
			Type:     schema.TypeSet,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Computed: true,
		},
	}
}
