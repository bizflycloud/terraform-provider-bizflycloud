package bizflycloud

import (
	"context"
	"errors"
	"fmt"
	"github.com/bizflycloud/gobizfly"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"strings"
	"time"
)

func resourceBizflyCloudWanIP() *schema.Resource {
	return &schema.Resource{
		Schema: resourceWanIPSchema(),
		Create: resourceBizflyCloudWanIPCreate,
		Read:   resourceBizflyCloudWanIPRead,
		Update: resourceBizflyCloudWanIPUpdate,
		Delete: resourceBizflyCloudWanIPDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		SchemaVersion: 1,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceBizflyCloudWanIPCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	createPayload := &gobizfly.CreatePublicNetworkInterfacePayload{
		Name:             d.Get("name").(string),
		AvailabilityZone: d.Get("availability_zone").(string),
	}
	wanIP, err := client.CloudServer.PublicNetworkInterfaces().Create(context.Background(), createPayload)
	if err != nil {
		return fmt.Errorf("error when creating wan ip: %s", err)
	}
	d.SetId(wanIP.ID)
	firewallIDs := readStringArray(d.Get("firewall_ids").(*schema.Set).List())
	if err := attachFirewallsForPort(client, wanIP.ID, firewallIDs); err != nil {
		return fmt.Errorf("error when attaching firewalls: %s", err)
	}
	return resourceBizflyCloudWanIPRead(d, meta)
}

func resourceBizflyCloudWanIPRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	var wanIP *gobizfly.CloudServerPublicNetworkInterface

	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		var err error
		log.Printf("[DEBUG] Reading WAN IP : %s", d.Id())
		wanIP, err = client.CloudServer.PublicNetworkInterfaces().Get(context.Background(), d.Id())

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
	_ = d.Set("firewall_ids", readSecurityGroups(wanIP.SecurityGroups))
	_ = d.Set("billing_type", wanIP.BillingType)
	_ = d.Set("bandwidth", wanIP.Bandwidth)
	_ = d.Set("availability_zone", wanIP.AvailabilityZone)
	_ = d.Set("server_id", wanIP.DeviceID)
	return nil
}

func resourceBizflyCloudWanIPDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	err := client.CloudServer.PublicNetworkInterfaces().Delete(context.Background(), d.Id())
	if err != nil && !strings.Contains(err.Error(), "Resource not found") {
		return fmt.Errorf("error when deleting WAN IP: %s", err)
	}
	return nil
}

func resourceBizflyCloudWanIPUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	if d.HasChange("attached_server") {
		serverId := d.Get("attached_server").(string)
		if serverId == "" {
			updatePayload := &gobizfly.ActionPublicNetworkInterfacePayload{
				Action: "detach_server",
			}
			err := client.CloudServer.PublicNetworkInterfaces().Action(context.Background(), d.Id(), updatePayload)
			if err != nil {
				return fmt.Errorf("error when detaching server: %s", err)
			}
		} else {
			updatePayload := &gobizfly.ActionPublicNetworkInterfacePayload{
				Action:   "attach_server",
				ServerId: serverId,
			}
			err := client.CloudServer.PublicNetworkInterfaces().Action(context.Background(), d.Id(), updatePayload)
			if err != nil {
				return fmt.Errorf("error when attaching server: %s", err)
			}
		}
	}
	if d.HasChange("billing_type") {
		billingType := d.Get("billing_type").(string)
		if billingType == "paid" {
			updatePayload := &gobizfly.ActionPublicNetworkInterfacePayload{
				Action: "convert_to_paid",
			}
			err := client.CloudServer.PublicNetworkInterfaces().Action(context.Background(), d.Id(), updatePayload)
			if err != nil {
				return fmt.Errorf("error when converting to paid: %s", err)
			}
		}
	}
	if d.HasChange("firewall_ids") {
		if err := updateFirewallForNetworkInterface(d, client, d.Id()); err != nil {
			return fmt.Errorf("error when update firewall for network interface: %s, %v", d.Id(), err)
		}
	}
	return resourceBizflyCloudWanIPRead(d, meta)
}
