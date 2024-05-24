package bizflycloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceBizflyCloudCustomImage() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"disk_format": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"container_format": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"billing_plan": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Read: dataSourceBizflyCloudCustomImageRead,
	}
}

func dataSourceBizflyCloudCustomImageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	resp, err := client.CloudServer.CustomImages().Get(context.Background(), d.Get("id").(string))
	if err != nil {
		return err
	}
	customImage := resp.Image
	d.SetId(customImage.ID)
	err = d.Set("name", customImage.Name)
	if err != nil {
		return err
	}
	err = d.Set("size", customImage.Size)
	if err != nil {
		return err
	}
	err = d.Set("disk_format", customImage.DiskFormat)
	if err != nil {
		return err
	}
	err = d.Set("container_format", customImage.ContainerFormat)
	if err != nil {
		return err
	}
	err = d.Set("billing_plan", customImage.BillingPlan)
	if err != nil {
		return err
	}
	return nil
}
