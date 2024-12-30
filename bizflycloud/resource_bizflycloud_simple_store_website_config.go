package bizflycloud

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bizflycloud/gobizfly"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBizflyCloudSimpleStorageBucketWebsiteConfig() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Create: resourceBizflyCloudSimpleStorageBucketWebsiteConfigUpdate,
		Read:   resourceBizflyCloudSimpleStorageBucketWebsiteConfigRead,
		Delete: func(d *schema.ResourceData, meta interface{}) error {
			log.Printf("[INFO] Delete operation is not supported for web config. Ignoring delete request.")
			return nil
		},
		Update: resourceBizflyCloudSimpleStorageBucketWebsiteConfigUpdate,
		Schema: map[string]*schema.Schema{
			"bucket_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"index": {
				Type:     schema.TypeString,
				Required: true,
			},
			"error": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Update: schema.DefaultTimeout(20 * time.Minute),
		},
	}
}

func resourceBizflyCloudSimpleStorageBucketWebsiteConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("error") {
		client := meta.(*CombinedConfig).gobizflyClient()

		ParamUpdateWebConfig := gobizfly.ParamUpdateWebsiteConfig{
			Index:      d.Get("index").(string),
			Error:      d.Get("error").(string),
			BucketName: d.Get("bucket_name").(string),
		}
		_, err := client.CloudSimpleStorage.UpdateWebsiteConfig(context.Background(), &ParamUpdateWebConfig)
		if err != nil {
			return fmt.Errorf("error updating simple store website config: %v", err)
		}
		d.SetId(d.Get("bucket_name").(string))
	}
	return resourceBizflyCloudSimpleStorageBucketWebsiteConfigRead(d, meta)
}

func resourceBizflyCloudSimpleStorageBucketWebsiteConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	paramListBucketInfo := gobizfly.ParamListWithBucketNameInfo{
		WebsiteConfig: "website_config",
		BucketName:    d.Get("bucket_name").(string),
	}
	dataBucket, err := client.CloudSimpleStorage.ListWithBucketNameInfo(context.Background(), paramListBucketInfo)
	if err != nil {
		return fmt.Errorf("Error when reading simple store website config: %v", err)
	}

	if err = d.Set("index", dataBucket.WebsiteConfig.Index); err != nil {
		return fmt.Errorf("Error setting website config state: %v", err)
	}
	return nil
}
