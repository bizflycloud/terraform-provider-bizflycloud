package bizflycloud

import (
	"context"
	"fmt"
	"log"

	"github.com/bizflycloud/gobizfly"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBizFlyCloudVPCNetwork() *schema.Resource {
	return &schema.Resource{
		Create:        resourceBizFlyCloudVPCNetworkCreate,
		Update:        resourceBizFlyCloudVPCNetworkUpdate,
		Delete:        resourceBizFlyCloudVPCNetworkDelete,
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cidr": {
				Type:     schema.TypeString,
				Required: true,
			},
			"is-default": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceBizFlyCloudVPCNetworkCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()

	log.Println("[DEBUG] creating vpc network")
	cvp := &gobizfly.CreateVPCPayload{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		CIDR:        d.Get("cidr").(string),
		IsDefault:   d.Get("is-default").(bool),
	}
	log.Printf("[DEBUG] Create vpc network configuration: %#v\n", cvp)
	data, err := client.VPC.Create(context.Background(), cvp)
	if err != nil {
		return fmt.Errorf("Error creating vpc network: %v", err)
	}
	log.Println("[DEBUG] set id " + data.ID)
	d.SetId(data.ID)

	return nil
}

func resourceBizFlyCloudVPCNetworkUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceBizFlyCloudVPCNetworkDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	err := client.VPC.Delete(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error delete vpc network: %v", err)
	}
	return nil
}
