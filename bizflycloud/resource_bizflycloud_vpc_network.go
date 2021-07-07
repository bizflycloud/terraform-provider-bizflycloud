package bizflycloud

import (
	"context"
	"fmt"
	"time"

	"github.com/bizflycloud/gobizfly"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBizFlyCloudVPCNetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourceBizFlyCloudVPCNetworkCreate,
		Read:   resourceBizFlyCloudVPCNetworkRead,
		Update: resourceBizFlyCloudVPCNetworkUpdate,
		Delete: resourceBizFlyCloudVPCNetworkDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		SchemaVersion: 1,
		Schema:        resourceVPCNetworkSchema(),
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func vpcRequestBuilder(d *schema.ResourceData) gobizfly.UpdateVPCPayload {
	vpcOpts := gobizfly.UpdateVPCPayload{}
	if v, ok := d.GetOk("name"); ok {
		vpcOpts.Name = v.(string)
	}
	if v, ok := d.GetOk("description"); ok {
		vpcOpts.Description = v.(string)
	}
	if v, ok := d.GetOk("cidr"); ok {
		vpcOpts.CIDR = v.(string)
	}
	if v, ok := d.GetOk("is_default"); ok {
		vpcOpts.IsDefault = v.(bool)
	}
	return vpcOpts
}

func resourceBizFlyCloudVPCNetworkCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	cvp := &gobizfly.CreateVPCPayload{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		CIDR:        d.Get("cidr").(string),
		IsDefault:   d.Get("is_default").(bool),
	}
	network, err := client.VPC.Create(context.Background(), cvp)
	if err != nil {
		return fmt.Errorf("Error creating VPC: %v", err)
	}
	d.SetId(network.ID)
	return resourceBizFlyCloudVPCNetworkRead(d, meta)
}

func resourceBizFlyCloudVPCNetworkRead(d *schema.ResourceData, meta interface{}) error {
	if err := dataSourceBizFlyCloudVPCNetworkRead(d, meta); err != nil {
		return err
	}

	return nil
}

func resourceBizFlyCloudVPCNetworkUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	vpcOpts := vpcRequestBuilder(d)
	network, err := client.VPC.Update(context.Background(), d.Id(), &vpcOpts)
	if err != nil {
		return fmt.Errorf("Error when update vpc network: %s, %v", d.Id(), err)
	}
	d.SetId(network.ID)
	return resourceBizFlyCloudVPCNetworkRead(d, meta)
}

func resourceBizFlyCloudVPCNetworkDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	err := client.VPC.Delete(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error delete vpc network: %v", err)
	}
	return nil
}
