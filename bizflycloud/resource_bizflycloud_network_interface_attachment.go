package bizflycloud

import (
	"context"
	"fmt"
	"log"

	"github.com/bizflycloud/gobizfly"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBizflyCloudNetworkInterfaceAttachment() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Create: resourceBizflyCloudNetworkInterfaceAttachmentCreate,
		Read:   resourceBizflyCloudNetworkInterfaceAttachmentRead,
		Delete: resourceBizflyCloudNetworkInterfaceAttachmentDelete,
		Schema: resourceNetworkInterfaceAttachmentSchema(),
	}
}

func resourceBizflyCloudNetworkInterfaceAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	serverID := d.Get("server_id").(string)
	networkInterfaceID := d.Get("network_interface_id").(string)
	if err := attachServerForPort(client, serverID, networkInterfaceID); err != nil {
		log.Printf("[ERROR] Error attaching server %s to network interface %s: %s", serverID, networkInterfaceID, err)
		return err
	}
	d.SetId(networkInterfaceID)
	return nil
}

func resourceBizflyCloudNetworkInterfaceAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	networkInterfaceID := d.Id()
	port, err := client.CloudServer.NetworkInterfaces().Get(context.Background(), networkInterfaceID)
	if err != nil {
		log.Printf("[ERROR] Error reading network interface %s: %s", networkInterfaceID, err)
		return err
	}
	d.SetId(port.ID)
	_ = d.Set("server_id", port.DeviceID)
	_ = d.Set("network_interface_id", port.ID)
	_ = d.Set("firewall_ids", port.SecurityGroups)
	return nil
}

func resourceBizflyCloudNetworkInterfaceAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	networkInterfaceID := d.Id()
	firewalls, err := client.CloudServer.Firewalls().List(context.Background(), &gobizfly.ListOptions{})
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
