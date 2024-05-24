package bizflycloud

import (
	"context"
	"fmt"
	"log"
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
		return fmt.Errorf("invalid network id specified")
	}
	createPayload := &gobizfly.CreateNetworkInterfacePayload{
		Name:    d.Get("name").(string),
		FixedIP: d.Get("fixed_ip").(string),
	}
	networkInterface, err := client.CloudServer.NetworkInterfaces().Create(context.Background(), networkID, createPayload)
	if err != nil {
		return fmt.Errorf("error when create network interface: %v", err)
	}
	d.SetId(networkInterface.ID)
	firewallIDs := readStringArray(d.Get("firewall_ids").(*schema.Set).List())
	log.Printf("[DEBUG] Firewall IDs: %v", firewallIDs)
	if attachFirewallsForPort(client, networkInterface.ID, firewallIDs) != nil {
		return fmt.Errorf("error when attach firewall for network interface: %v", err)
	}
	return resourceBizflyCloudNetworkInterfaceRead(d, meta)
}

func resourceBizflyCloudNetworkInterfaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	networkInterface, err := client.CloudServer.NetworkInterfaces().Get(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error read network interface network %s: %w", d.Id(), err)
	}
	d.SetId(networkInterface.ID)
	_ = d.Set("name", networkInterface.Name)
	_ = d.Set("network_id", networkInterface.NetworkID)
	_ = d.Set("fixed_ip", networkInterface.FixedIps[0].IPAddress)
	_ = d.Set("status", networkInterface.Status)
	_ = d.Set("created_at", networkInterface.CreatedAt)
	_ = d.Set("updated_at", networkInterface.UpdatedAt)
	_ = d.Set("server_id", networkInterface.DeviceID)
	if err := d.Set("fixed_ips", readFixedIps(networkInterface.FixedIps)); err != nil {
		return fmt.Errorf("error setting fixed_ips: %w", err)
	}
	if err := d.Set("firewall_ids", networkInterface.SecurityGroups); err != nil {
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
		_, err := client.CloudServer.NetworkInterfaces().Update(context.Background(), d.Id(), updatePayload)
		if err != nil {
			return fmt.Errorf("error when update network interface: %s, %v", d.Id(), err)
		}
	}
	if d.HasChange("firewall_ids") {
		if err := updateFirewallForNetworkInterface(d, client, d.Id()); err != nil {
			return fmt.Errorf("error when update firewall for network interface: %s, %v", d.Id(), err)
		}
	}
	return resourceBizflyCloudNetworkInterfaceRead(d, meta)
}

func resourceBizflyCloudNetworkInterfaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	err := client.CloudServer.NetworkInterfaces().Delete(context.Background(), d.Id())
	if err != nil && strings.Contains(err.Error(), "Resource not found") {
		return fmt.Errorf("error when delete network interface: %v", err)
	}
	return nil
}

func updateFirewallForNetworkInterface(d *schema.ResourceData, client *gobizfly.Client,
	networkInterfaceID string) error {
	firewallIDs := readStringArray(d.Get("firewall_ids").(*schema.Set).List())
	if err := attachFirewallsForPort(client, networkInterfaceID, firewallIDs); err != nil {
		log.Printf("[ERROR] Error attaching firewalls %s to network interface %s: %s", firewallIDs, networkInterfaceID, err)
		return err
	}
	oldFirewalls, newFirewalls := d.GetChange("firewall_ids")
	oldFirewallIDs := newSet(oldFirewalls.(*schema.Set).List())
	newFirewallIDs := newSet(newFirewalls.(*schema.Set).List())
	addFirewallIds := leftDiff(newFirewallIDs, oldFirewallIDs)
	removeFirewallIds := leftDiff(oldFirewallIDs, newFirewallIDs)
	addFirewallIdsArray := make([]string, 0, len(addFirewallIds))
	removeFirewallIdsArray := make([]string, 0, len(removeFirewallIds))
	for id := range addFirewallIds {
		addFirewallIdsArray = append(addFirewallIdsArray, id)
	}
	for id := range removeFirewallIds {
		removeFirewallIdsArray = append(removeFirewallIdsArray, id)
	}

	log.Printf("[DEBUG] Add firewalls %s to network interface %s", addFirewallIdsArray, networkInterfaceID)
	if err := attachFirewallsForPort(client, networkInterfaceID, addFirewallIdsArray); err != nil {
		log.Printf("[ERROR] Error attaching firewalls %s to network interface %s: %s", addFirewallIds, networkInterfaceID, err)
		return err
	}
	log.Printf("[DEBUG] Remove firewalls %s from network interface %s", removeFirewallIdsArray, networkInterfaceID)
	if err := detachFirewallsForPort(client, networkInterfaceID, removeFirewallIdsArray); err != nil {
		log.Printf("[ERROR] Error detaching firewalls %s from network interface %s: %s", removeFirewallIds, networkInterfaceID, err)
		return err
	}
	return nil
}
