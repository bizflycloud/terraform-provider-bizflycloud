package bizflycloud

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bizflycloud/gobizfly"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBizflySimpleStoreVersioning() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Create: resourceBizflyCloudSimpleStoreVersioningUpdate,
		Read:   resourceBizflyCloudSimpleStoreVersioningRead,
		Delete: func(d *schema.ResourceData, meta interface{}) error {
			log.Printf("[INFO] Delete operation is not supported for versioning. Ignoring delete request.")
			return nil
		},
		Update: resourceBizflyCloudSimpleStoreVersioningUpdate,
		Schema: map[string]*schema.Schema{
			"bucket_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"versioning": {
				Type:     schema.TypeBool,
				Required: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Update: schema.DefaultTimeout(20 * time.Minute),
		},
	}
}

func resourceBizflyCloudSimpleStoreVersioningUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("versioning") {
		client := meta.(*CombinedConfig).gobizflyClient()
		bucketName := d.Get("bucket_name").(string)
		versioning := d.Get("versioning").(bool)

		_, err := client.CloudSimpleStoreBucket.UpdateVersioning(context.Background(), versioning, bucketName)
		if err != nil {
			return fmt.Errorf("error updating simple store versioning: %v", err)
		}
		d.SetId(bucketName)
	}
	return resourceBizflyCloudSimpleStoreVersioningRead(d, meta)
}

func resourceBizflyCloudSimpleStoreVersioningRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	paramGetPath := gobizfly.ParamListWithBucketNameInfo{
		Versioning: "versioning",
		BucketName: d.Get("bucket_name").(string),
	}
	dataBuckets, err := client.CloudSimpleStoreBucket.ListWithBucketNameInfo(context.Background(), paramGetPath)
	if err != nil {
		return fmt.Errorf("Error when reading simple store Verioning: %v", err)
	}
	versioningEnabled := false
	if dataBuckets.Versioning.Status == "Enabled" {
		versioningEnabled = true
	}

	if err = d.Set("versioning", versioningEnabled); err != nil {
		return fmt.Errorf("Error setting versioning state: %v", err)
	}
	return nil
}
