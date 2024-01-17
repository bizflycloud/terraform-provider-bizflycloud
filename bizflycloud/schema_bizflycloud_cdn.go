package bizflycloud

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

func resourceCDNSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"domain": {
			Type:     schema.TypeString,
			Required: true,
		},
		"domain_cdn": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"tenant_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"origin": {
			Type:     schema.TypeMap,
			Required: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			// Elem: &schema.Resource{
			// 	Schema: dataOriginSchema(),
			// },
		},
	}
}

// func dataOriginSchema() map[string]*schema.Schema {
// 	return map[string]*schema.Schema{
// 		"name": {
// 			Optional: true,
// 			Type:     schema.TypeString,
// 		},
// 		"origin_type": {
// 			Required: false,
// 			Optional: true,
// 			Type:     schema.TypeString,
// 			Default:  "custom_origin",
// 		},
// 		"upstream_addrs": {
// 			Required: true,
// 			Type:     schema.TypeString,
// 		},
// 		"upstream_host": {
// 			Required: true,
// 			Type:     schema.TypeString,
// 		},
// 		"upstream_proto": {
// 			Required: true,
// 			Type:     schema.TypeString,
// 		},
// 	}
// }
