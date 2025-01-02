package bizflycloud

import (
	"context"
	"fmt"
	"github.com/bizflycloud/gobizfly"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBizflyCloudSimpleStorageAccessKey() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Create: resourceBizflyCloudSimpleStorageKeyCreate,
		Read:   resourceBizflyCloudSimpleStorageKeyRead,
		Delete: resourceBizflyCloudSimpleStorageKeyDelete,
		Schema: map[string]*schema.Schema{
			"subuser_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"access_key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"secret_key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
		},
	}
}

func resourceBizflyCloudSimpleStorageKeyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	kcr := gobizfly.KeyCreateRequest{
		SubuserId: d.Get("subuser_id").(string),
		AccessKey: d.Get("access_key").(string),
		SecretKey: d.Get("secret_key").(string),
	}
	_, err := client.CloudSimpleStorage.SimpleStorageKey().CreateAccessKey(context.Background(), &kcr)
	if err != nil {
		return fmt.Errorf("Error when creating simple store key: %v", err)
	}
	d.SetId(kcr.AccessKey)
	return resourceBizflyCloudSimpleStorageKeyRead(d, meta)
}

func resourceBizflyCloudSimpleStorageKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	keys, err := client.CloudSimpleStorage.SimpleStorageKey().ListAccessKey(context.Background(), &gobizfly.ListOptions{})
	if err != nil {
		return fmt.Errorf("Error retrieving simple store key: %v", err)
	}
	for _, key := range keys {
		_ = d.Set("name", key.AccessKey)
		_ = d.Set("location", key.User)
	}
	return nil
}

func resourceBizflyCloudSimpleStorageKeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	err := client.CloudSimpleStorage.SimpleStorageKey().DeleteAccessKey(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error deleting simple store key: %v", err)
	}
	return nil
}
