package bizflycloud

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bizflycloud/gobizfly"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBizflyCloudSimpleStorageBucketVersioning() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Create: resourceBizflyCloudSimpleStorageBucketVersioningUpdate,
		Read:   resourceBizflyCloudSimpleStorageBucketVersioningRead,
		Delete: func(d *schema.ResourceData, meta interface{}) error {
			log.Printf("[INFO] Delete operation is not supported for versioning. Ignoring delete request.")
			return nil
		},
		Update: resourceBizflyCloudSimpleStorageBucketVersioningUpdate,
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

func resourceBizflyCloudSimpleStorageBucketVersioningUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("versioning") {
		client := meta.(*CombinedConfig).gobizflyClient()
		bucketName := d.Get("bucket_name").(string)
		versioning := d.Get("versioning").(bool)

		_, err := client.CloudSimpleStorage.UpdateVersioning(context.Background(), versioning, bucketName)
		if err != nil {
			return fmt.Errorf("error updating simple store versioning: %v", err)
		}
		d.SetId(bucketName)
	}
	return resourceBizflyCloudSimpleStorageBucketVersioningRead(d, meta)
}

func resourceBizflyCloudSimpleStorageBucketVersioningRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	paramListBucketInfo := gobizfly.ParamListWithBucketNameInfo{
		Versioning: "versioning",
		BucketName: d.Get("bucket_name").(string),
	}
	dataBucket, err := client.CloudSimpleStorage.ListWithBucketNameInfo(context.Background(), paramListBucketInfo)
	if err != nil {
		return fmt.Errorf("error when reading simple store Verioning: %v", err)
	}
	versioningEnabled := dataBucket.Versioning.Status == "Enabled"

	if err = d.Set("versioning", versioningEnabled); err != nil {
		return fmt.Errorf("error setting versioning state: %v", err)
	}
	return nil
}
