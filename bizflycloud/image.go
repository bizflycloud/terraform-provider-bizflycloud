package bizflycloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func imageSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Description: "Id of the image",
			Computed:    true,
		},
		"distribution": {
			Type:        schema.TypeString,
			Description: "OS Distribution",
			Required:    true,
		},
		"version": {
			Type:        schema.TypeString,
			Description: "OS Version",
			Required:    true,
		},
	}
}
