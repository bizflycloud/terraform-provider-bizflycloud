package bizflycloud

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/bizflycloud/gobizfly"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
				Optional: true,
			},
			"is_default": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(600 * time.Second),
		},
	}
}

func VPCRequestBuilder(d *schema.ResourceData) gobizfly.UpdateVPCPayload {
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
	client := meta.(*CombinedConfig).gobizflyClient()
	var network *gobizfly.VPC

	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		var err error
		network, err = client.VPC.Get(context.Background(), d.Id())

		// Retry on any API "not found" errors, but only on new resources.
		if d.IsNewResource() && errors.Is(err, gobizfly.ErrNotFound) {
			return resource.RetryableError(err)
		}
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})

	// Prevent confusing Terraform error messaging to operators by
	// Only ignoring API "not found" errors if not a new resource
	if !d.IsNewResource() && errors.Is(err, gobizfly.ErrNotFound) {
		log.Printf("[WARN] VPC network %s is not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("Error create vpc network %s: %w", d.Id(), err)
	}

	// Prevent panics.
	if network == nil {
		return fmt.Errorf("error create vpc network (%s): empty response", d.Id())
	}
	_ = d.Set("name", network.Name)
	_ = d.Set("description", network.Description)
	_ = d.Set("is_default", network.IsDefault)
	_ = d.Set("id", network.ID)
	return nil
}

func resourceBizFlyCloudVPCNetworkUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	vpcOpts := VPCRequestBuilder(d)
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
