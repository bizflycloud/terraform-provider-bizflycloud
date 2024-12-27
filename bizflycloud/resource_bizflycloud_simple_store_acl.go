package bizflycloud

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bizflycloud/gobizfly"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBizflyCloudSimpleStoreAcl() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Create: resourceBizflyCloudSimpleStoreAclUpdate,
		Read:   resourceBizflyCloudSimpleStoreAclRead,
		Delete: func(d *schema.ResourceData, meta interface{}) error {
			log.Printf("[INFO] Delete operation is not supported for ACL. Ignoring delete request.")
			return nil
		},
		Update: resourceBizflyCloudSimpleStoreAclUpdate,
		Schema: map[string]*schema.Schema{
			"bucket_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"acl": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Update: schema.DefaultTimeout(20 * time.Minute),
		},
	}
}

func resourceBizflyCloudSimpleStoreAclUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("acl") {
		client := meta.(*CombinedConfig).gobizflyClient()
		bucketName := d.Get("bucket_name").(string)
		acl := d.Get("acl").(string)

		_, err := client.CloudSimpleStorage.UpdateAcl(context.Background(), acl, bucketName)
		if err != nil {
			return fmt.Errorf("error updating simple store ACL: %v", err)
		}
		d.SetId(bucketName)
	}
	return resourceBizflyCloudSimpleStoreAclRead(d, meta)
}

func resourceBizflyCloudSimpleStoreAclRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	paramListBucketInfo := gobizfly.ParamListWithBucketNameInfo{
		Acl:        "acl",
		BucketName: d.Get("bucket_name").(string),
	}
	dataBuckets, err := client.CloudSimpleStorage.ListWithBucketNameInfo(context.Background(), paramListBucketInfo)
	if err != nil {
		return fmt.Errorf("Error when reading simple store Acl: %v", err)
	}

	_ = d.Set("owner", dataBuckets.Acl.Owner)
	_ = d.Set("message", dataBuckets.Acl.Message)
	_ = d.Set("grants", dataBuckets.Acl.Grants)

	return nil
}
