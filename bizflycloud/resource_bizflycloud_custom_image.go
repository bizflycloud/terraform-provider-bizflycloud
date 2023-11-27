package bizflycloud

import (
	"context"
	"github.com/bizflycloud/gobizfly"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
)

func resourceBizflyCloudCustomImage() *schema.Resource {
	return &schema.Resource{
		Create: resourceBizflyCloudCustomImageCreate,
		Read:   resourceBizflyCloudCustomImageRead,
		Delete: resourceBizflyCloudCustomImageDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"disk_format": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"qcow2",
					"raw",
					"vdi",
					"vmdk",
					"vhd",
				}, false),
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"image_url": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"container_format": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"billing_plan": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"visibility": {
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
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceBizflyCloudCustomImageCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	req := gobizfly.CreateCustomImagePayload{
		Name:        d.Get("name").(string),
		DiskFormat:  d.Get("disk_format").(string),
		ImageURL:    d.Get("image_url").(string),
		Description: d.Get("description").(string),
	}

	resp, err := client.Server.CreateCustomImage(context.Background(), &req)
	if err != nil {
		log.Printf("[ERROR] Error create custom image: %v", err)
		return err
	}
	d.SetId(resp.Image.ID)
	return resourceBizflyCloudCustomImageRead(d, meta)
}

func resourceBizflyCloudCustomImageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	resp, err := client.Server.GetCustomImage(context.Background(), d.Id())
	if err != nil {
		log.Printf("[ERROR] Error read custom image %s: %v", d.Id(), err)
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
	err = d.Set("visibility", customImage.Visibility)
	if err != nil {
		return err
	}
	err = d.Set("created_at", customImage.CreatedAt)
	if err != nil {
		return err
	}
	err = d.Set("updated_at", customImage.UpdatedAt)
	if err != nil {
		return err
	}
	err = d.Set("description", customImage.Description)
	if err != nil {
		return err
	}
	return nil
}

func resourceBizflyCloudCustomImageDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	err := client.Server.DeleteCustomImage(context.Background(), d.Id())
	if err != nil {
		return err
	}
	return nil
}
