package bizflycloud

import (
	"context"
	"fmt"
	"github.com/bizflycloud/gobizfly"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

func resourceBizflyCloudNetworkInterfaceAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceBizflyCloudNetworkInterfaceAttachmentCreate,
		Read:   resourceBizflyCloudNetworkInterfaceAttachmentRead,
		Update: resourceBizflyCloudNetworkInterfaceAttachmentUpdate,
		Delete: resourceBizflyCloudNetworkInterfaceAttachmentDelete,
		Schema: resourceNetworkInterfaceAttachmentSchema(),
	}
}

func resourceBizflyCloudNetworkInterfaceAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	serverID := d.Get("server_id").(string)
	networkInterfaceID := d.Get("network_interface_id").(string)
	firewallIDs := readStringArray(d.Get("firewall_ids").(*schema.Set).List())
	if err := attachServerForPort(client, serverID, networkInterfaceID); err != nil {
		log.Printf("[ERROR] Error attaching server %s to network interface %s: %s", serverID, networkInterfaceID, err)
		return err
	}
	if err := attachFirewallsForPort(client, networkInterfaceID, firewallIDs); err != nil {
		log.Printf("[ERROR] Error attaching firewalls %s to network interface %s: %s", firewallIDs, networkInterfaceID, err)
		return err
	}
	d.SetId(networkInterfaceID)
	return nil
}

func resourceBizflyCloudNetworkInterfaceAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	networkInterfaceID := d.Id()
	firewalls, err := client.Firewall.List(context.Background(), &gobizfly.ListOptions{})
	if err != nil {
		return fmt.Errorf("error listing firewalls: %v", err)
	}
	port, err := client.NetworkInterface.Get(context.Background(), networkInterfaceID)
	if err != nil {
		log.Printf("[ERROR] Error reading network interface %s: %s", networkInterfaceID, err)
		return err
	}
	userFirewallIDs := make([]string, len(firewalls))
	for i, firewall := range firewalls {
		userFirewallIDs[i] = firewall.ID
	}

	d.SetId(port.ID)
	_ = d.Set("server_id", port.DeviceID)
	_ = d.Set("network_interface_id", port.ID)
	_ = d.Set("firewall_ids", filterUserFirewalls(port.SecurityGroups, userFirewallIDs))
	return nil
}

func resourceBizflyCloudNetworkInterfaceAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	networkInterfaceID := d.Id()
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

func resourceBizflyCloudNetworkInterfaceAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	networkInterfaceID := d.Id()
	firewalls, err := client.Firewall.List(context.Background(), &gobizfly.ListOptions{})
	if err != nil {
		return fmt.Errorf("error listing firewalls: %v", err)
	}
	userFirewallIDs := make([]string, len(firewalls))
	for i, firewall := range firewalls {
		userFirewallIDs[i] = firewall.ID
	}
	if err := detachFirewallsForPort(client, networkInterfaceID, userFirewallIDs); err != nil {
		log.Printf("[ERROR] Error detaching firewalls %s from network interface %s: %s", userFirewallIDs, networkInterfaceID, err)
		return err
	}
	if err := detachServerForPort(client, networkInterfaceID); err != nil {
		log.Printf("[ERROR] Error detaching server from network interface %s: %s", networkInterfaceID, err)
		return err
	}
	return nil
}
