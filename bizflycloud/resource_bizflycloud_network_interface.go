package bizflycloud

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bizflycloud/gobizfly"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBizflyCloudNetworkInterface() *schema.Resource {
	return &schema.Resource{
		Create: resourceBizflyCloudNetworkInterfaceCreate,
		Read:   resourceBizflyCloudNetworkInterfaceRead,
		Update: resourceBizflyCloudNetworkInterfaceUpdate,
		Delete: resourceBizflyCloudNetworkInterfaceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		SchemaVersion: 1,
		Schema:        resourceNetworkInterfaceSchema(),
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceBizflyCloudNetworkInterfaceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	networkID := d.Get("network_id").(string)
	if networkID == "" {
		return fmt.Errorf("Invalid network id specified")
	}
	createPayload := &gobizfly.CreateNetworkInterfacePayload{
		Name:           d.Get("name").(string),
		AttachedServer: d.Get("attached_server").(string),
		FixedIP:        d.Get("fixed_ip").(string),
	}
	networkInterface, err := client.NetworkInterface.Create(context.Background(), networkID, createPayload)
	if err != nil {
		return fmt.Errorf("Error when create network interface: %v", err)
	}
	d.SetId(networkInterface.ID)

	return resourceBizflyCloudNetworkInterfaceRead(d, meta)
}

func resourceBizflyCloudNetworkInterfaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	networkInterface, err := client.NetworkInterface.Get(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error read network interface network %s: %w", d.Id(), err)
	}
	d.SetId(networkInterface.ID)
	_ = d.Set("name", networkInterface.Name)
	_ = d.Set("network_id", networkInterface.NetworkID)
	_ = d.Set("attached_server", networkInterface.AttachedServer)
	_ = d.Set("fixed_ip", networkInterface.FixedIps[0].IPAddress)
	_ = d.Set("status", networkInterface.Status)
	_ = d.Set("created_at", networkInterface.CreatedAt)
	_ = d.Set("updated_at", networkInterface.UpdatedAt)
	if err := d.Set("fixed_ips", readFixedIps(networkInterface.FixedIps)); err != nil {
		return fmt.Errorf("error setting fixed_ips: %w", err)
	}

	if err := d.Set("security_groups", readSecurityGroups(networkInterface.SecurityGroups)); err != nil {
		return fmt.Errorf("error setting security_groups: %w", err)
	}
	return nil
}

func resourceBizflyCloudNetworkInterfaceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	if d.HasChange("name") {
		updatePayload := &gobizfly.UpdateNetworkInterfacePayload{
			Name: d.Get("name").(string),
		}
		_, err := client.NetworkInterface.Update(context.Background(), d.Id(), updatePayload)
		if err != nil {
			return fmt.Errorf("Error when update network interface: %s, %v", d.Id(), err)
		}
	}
	return resourceBizflyCloudNetworkInterfaceRead(d, meta)
}

func resourceBizflyCloudNetworkInterfaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	err := client.NetworkInterface.Delete(context.Background(), d.Id())
	if err != nil && strings.Contains(err.Error(), "Resource not found") {
		return fmt.Errorf("Error when delete network interface: %v", err)
	}
	return nil
}
