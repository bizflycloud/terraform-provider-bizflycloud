package bizflycloud

import (
	"context"
	"fmt"
	"time"

	"github.com/bizflycloud/gobizfly"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBizFlyCloudNetworkInterface() *schema.Resource {
	return &schema.Resource{
		Create: resourceBizFlyCloudNetworkInterfaceCreate,
		Read:   resourceBizFlyCloudNetworkInterfaceRead,
		Update: resourceBizFlyCloudNetworkInterfaceUpdate,
		Delete: resourceBizFlyCloudNetworkInterfaceDelete,
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

func resourceBizFlyCloudNetworkInterfaceCreate(d *schema.ResourceData, meta interface{}) error {
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

	actionPayload := &gobizfly.ActionNetworkInterfacePayload{
		Action:   d.Get("action").(string),
		ServerID: d.Get("server_id").(string),
	}
	if v, ok := d.GetOk("security_groups"); ok {
		for _, id := range v.([]interface{}) {
			actionPayload.SecurityGroups = append(actionPayload.SecurityGroups, id.(string))
		}
	}
	_, err = client.NetworkInterface.Action(context.Background(), d.Id(), actionPayload)
	if err != nil {
		return fmt.Errorf("Error when add firewall network interface: %v", err)
	}

	return resourceBizFlyCloudNetworkInterfaceRead(d, meta)
}

func resourceBizFlyCloudNetworkInterfaceRead(d *schema.ResourceData, meta interface{}) error {
	if err := dataSourceBizFlyCloudNetworkInterfaceRead(d, meta); err != nil {
		return err
	}
	return nil
}

func resourceBizFlyCloudNetworkInterfaceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	updatePayload := &gobizfly.UpdateNetworkInterfacePayload{
		Name: d.Get("name").(string),
	}
	networkInterface, err := client.NetworkInterface.Update(context.Background(), d.Id(), updatePayload)
	if err != nil {
		return fmt.Errorf("Error when update network interface: %s, %v", d.Id(), err)
	}
	d.SetId(networkInterface.ID)

	actionPayload := &gobizfly.ActionNetworkInterfacePayload{
		Action:   d.Get("action").(string),
		ServerID: d.Get("server_id").(string),
	}
	if v, ok := d.GetOk("security_groups"); ok {
		for _, id := range v.([]interface{}) {
			actionPayload.SecurityGroups = append(actionPayload.SecurityGroups, id.(string))
		}
	}
	_, err = client.NetworkInterface.Action(context.Background(), d.Id(), actionPayload)
	if err != nil {
		return fmt.Errorf("Error when add firewall network interface: %v", err)
	}

	return resourceBizFlyCloudNetworkInterfaceRead(d, meta)
}

func resourceBizFlyCloudNetworkInterfaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	err := client.NetworkInterface.Delete(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error when delete network interface: %v", err)
	}
	return nil
}
