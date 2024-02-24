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
		"domain_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"origin": {
			Type:     schema.TypeSet,
			Required: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"upstream_addrs": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "A domain name or IP of your origin source. Specify a port if custom.",
					},
					"upstream_host": {
						Type:     schema.TypeString,
						Required: true,
					},
					"upstream_proto": {
						Type:     schema.TypeString,
						Optional: true,
						Default:  "http",
					},
					"name": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
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
