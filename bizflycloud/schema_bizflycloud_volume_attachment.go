package bizflycloud

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

func resourceBizflyCloudVolumeAttachmentSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"volume_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"server_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
	}
}
