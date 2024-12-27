package bizflycloud

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bizflycloud/gobizfly"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBizflySimpleStoreCors() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Create: resourceBizflyCloudSimpleStoreCorsUpdate,
		Read:   resourceBizflyCloudSimpleStoreCorsRead,
		Delete: func(d *schema.ResourceData, meta interface{}) error {
			log.Printf("[INFO] Delete operation is not supported for CORS. Ignoring delete request.")
			return nil
		},
		Update: resourceBizflyCloudSimpleStoreCorsUpdate,
		Schema: map[string]*schema.Schema{
			"bucket_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"rules": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allowed_origin": {
							Type:     schema.TypeString,
							Required: true,
						},
						"allowed_methods": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Required: true,
						},
						"allowed_headers": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Required: true,
						},
						"max_age_seconds": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Update: schema.DefaultTimeout(20 * time.Minute),
		},
	}
}

func resourceBizflyCloudSimpleStoreCorsUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()

	rulesData := d.Get("rules").(*schema.Set).List()
	rules := make([]gobizfly.Rule, len(rulesData))

	for i, rule := range rulesData {
		ruleMap := rule.(map[string]interface{})
		rules[i] = gobizfly.Rule{
			AllowedOrigin:  ruleMap["allowed_origin"].(string),
			AllowedMethods: convertInterfaceSliceToStringSlice(ruleMap["allowed_methods"].([]interface{})),
			AllowedHeaders: convertInterfaceSliceToStringSlice(ruleMap["allowed_headers"].([]interface{})),
			MaxAgeSeconds:  ruleMap["max_age_seconds"].(int),
		}
	}

	paramCors := gobizfly.ParamUpdateCors{
		Rules:      rules,
		BucketName: d.Get("bucket_name").(string),
	}

	_, err := client.CloudSimpleStorage.UpdateCors(context.Background(), &paramCors)
	if err != nil {
		return fmt.Errorf("error updating simple store CORS: %v", err)
	}

	d.SetId(d.Get("bucket_name").(string))
	return resourceBizflyCloudSimpleStoreCorsRead(d, meta)
}

func convertInterfaceSliceToStringSlice(input []interface{}) []string {
	output := make([]string, len(input))
	for i, v := range input {
		output[i] = v.(string)
	}
	return output
}

func resourceBizflyCloudSimpleStoreCorsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	paramGetPath := gobizfly.ParamListWithBucketNameInfo{
		Cors:       "cors",
		BucketName: d.Get("bucket_name").(string),
	}
	dataBuckets, err := client.CloudSimpleStorage.ListWithBucketNameInfo(context.Background(), paramGetPath)
	if err != nil {
		return fmt.Errorf("Error when reading simple store Cors: %v", err)
	}
	d.SetId(dataBuckets.Bucket.Name)

	return nil
}
