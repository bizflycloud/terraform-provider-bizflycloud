package bizflycloud

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"log"
	"time"

	"github.com/bizflycloud/gobizfly"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBizflyCloudVPCNetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourceBizflyCloudVPCNetworkCreate,
		Read:   resourceBizflyCloudVPCNetworkRead,
		Update: resourceBizflyCloudVPCNetworkUpdate,
		Delete: resourceBizflyCloudVPCNetworkDelete,
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

func resourceBizflyCloudVPCNetworkCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	cvp := &gobizfly.CreateVPCPayload{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		CIDR:        d.Get("cidr").(string),
		IsDefault:   d.Get("is_default").(bool),
	}
	network, err := client.VPC.Create(context.Background(), cvp)
	if err != nil {
		return fmt.Errorf("Error when create vpc network: %v", err)
	}
	d.SetId(network.ID)
	return resourceBizflyCloudVPCNetworkRead(d, meta)
}

func resourceBizflyCloudVPCNetworkRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	var vpc *gobizfly.VPC
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		var err error

		log.Printf("[DEBUG] Reading vpc: %s", d.Id())
		vpc, err = client.VPC.Get(context.Background(), d.Id())

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
		log.Printf("[WARN] vpc network %s is not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("Error read vpc network VPC %s: %w", d.Id(), err)
	}

	// Prevent panics.
	if vpc == nil {
		return fmt.Errorf("Error read vpc network (%s): empty response", d.Id())
	}

	d.SetId(vpc.ID)
	_ = d.Set("name", vpc.Name)
	_ = d.Set("description", vpc.Description)
	_ = d.Set("is_default", vpc.IsDefault)
	_ = d.Set("mtu", vpc.MTU)
	_ = d.Set("status", vpc.Status)
	_ = d.Set("created_at", vpc.CreatedAt)
	_ = d.Set("updated_at", vpc.UpdatedAt)
	_ = d.Set("cidr", vpc.Subnets[0].CIDR)
	_ = d.Set("mtu", vpc.MTU)

	if err := d.Set("availability_zones", readAvailabilityZones(vpc.AvailabilityZones)); err != nil {
		return fmt.Errorf("error setting availability_zones: %w", err)
	}

	if err := d.Set("subnets", readSubnets(vpc.Subnets)); err != nil {
		return fmt.Errorf("error setting subnets: %w", err)
	}

	return nil
}

func resourceBizflyCloudVPCNetworkUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	vpcOpts := VPCRequestBuilder(d)
	network, err := client.VPC.Update(context.Background(), d.Id(), &vpcOpts)
	if err != nil {
		return fmt.Errorf("Error when update vpc network: %s, %v", d.Id(), err)
	}
	d.SetId(network.ID)
	return resourceBizflyCloudVPCNetworkRead(d, meta)
}

func resourceBizflyCloudVPCNetworkDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	err := client.VPC.Delete(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error when delete vpc network: %v", err)
	}
	return nil
}
