package bizflycloud

import (
	"context"
	"errors"
	"fmt"
	"github.com/bizflycloud/gobizfly"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"time"
)

func resourceBizFlyCloudWanIP() *schema.Resource {
	return &schema.Resource{
		Schema: resourceWanIPSchema(),
		Create: resourceBizFlyCloudWanIPCreate,
		Read:   resourceBizFlyCloudWanIPRead,
		Update: resourceBizFlyCloudWanIPUpdate,
		Delete: resourceBizFlyCloudWanIPDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		SchemaVersion: 1,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceBizFlyCloudWanIPCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	createPayload := &gobizfly.CreateWanIpPayload{
		Name:             d.Get("name").(string),
		AttachedServer:   d.Get("attached_server").(string),
		AvailabilityZone: d.Get("availability_zone").(string),
	}
	wanIP, err := client.WanIP.Create(context.Background(), createPayload)
	if err != nil {
		return fmt.Errorf("error when creating wan ip: %s", err)
	}
	d.SetId(wanIP.ID)
	return resourceBizFlyCloudWanIPRead(d, meta)
}

func resourceBizFlyCloudWanIPRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	var wanIP *gobizfly.WanIP

	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		var err error
		log.Printf("[DEBUG] Reading WAN IP : %s", d.Id())
		wanIP, err = client.WanIP.Get(context.Background(), d.Id())

		// Retry on any API "not found" errors, but only on new resources.
		if d.IsNewResource() && errors.Is(err, gobizfly.ErrNotFound) {
			return resource.RetryableError(err)
		}
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if !d.IsNewResource() && errors.Is(err, gobizfly.ErrNotFound) {
		log.Printf("[WARN] WAN IP %s is not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error read WAN IP network %s: %w", d.Id(), err)
	}

	// Prevent panics.
	if wanIP == nil {
		return fmt.Errorf("error read WAN IP (%s): empty response", d.Id())
	}
	d.SetId(wanIP.ID)
	_ = d.Set("name", wanIP.Name)
	_ = d.Set("network_id", wanIP.NetworkID)
	_ = d.Set("ip_address", wanIP.IpAddress)
	_ = d.Set("ip_version", wanIP.IpVersion)
	_ = d.Set("status", wanIP.Status)
	_ = d.Set("created_at", wanIP.CreatedAt)
	_ = d.Set("updated_at", wanIP.UpdatedAt)
	_ = d.Set("security_groups", readSecurityGroups(wanIP.SecurityGroups))
	_ = d.Set("billing_type", wanIP.BillingType)
	_ = d.Set("bandwidth", wanIP.Bandwidth)
	_ = d.Set("availability_zone", wanIP.AvailabilityZone)
	return nil
}

func resourceBizFlyCloudWanIPDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	err := client.WanIP.Delete(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("error when deleting WAN IP: %s", err)
	}
	return nil
}

func resourceBizFlyCloudWanIPUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	if d.HasChange("attached_server") {
		serverId := d.Get("attached_server").(string)
		if serverId == "" {
			updatePayload := &gobizfly.ActionWanIpPayload{
				Action: "detach_server",
			}
			err := client.WanIP.Action(context.Background(), d.Id(), updatePayload)
			if err != nil {
				return fmt.Errorf("error when detaching server: %s", err)
			}
		} else {
			updatePayload := &gobizfly.ActionWanIpPayload{
				Action:   "attach_server",
				ServerId: serverId,
			}
			err := client.WanIP.Action(context.Background(), d.Id(), updatePayload)
			if err != nil {
				return fmt.Errorf("error when attaching server: %s", err)
			}
		}
	}
	if d.HasChange("billing_type") {
		billingType := d.Get("billing_type").(string)
		if billingType == "paid" {
			updatePayload := &gobizfly.ActionWanIpPayload{
				Action: "convert_to_paid",
			}
			err := client.WanIP.Action(context.Background(), d.Id(), updatePayload)
			if err != nil {
				return fmt.Errorf("error when converting to paid: %s", err)
			}
		}
	}
	return resourceBizFlyCloudWanIPRead(d, meta)
}
