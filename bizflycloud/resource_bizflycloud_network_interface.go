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
		Update: resourceBizFlyCloudNetworkInterfacekUpdate,
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

func NetworkInterfaceRequestBuilder(d *schema.ResourceData) gobizfly.UpdateNetworkInterfacePayload {
	networkInterfaceOpts := gobizfly.UpdateNetworkInterfacePayload{}
	if v, ok := d.GetOk("name"); ok {
		networkInterfaceOpts.Name = v.(string)
	}
	return networkInterfaceOpts
}

func resourceBizFlyCloudNetworkInterfaceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()

	networkID := d.Get("network_id").(string)
	if networkID == "" {
		return fmt.Errorf("Invalid network id specified")
	}

	cnp := &gobizfly.CreateNetworkInterfacePayload{
		AttachedServer: d.Get("attached_server").(string),
		FixedIP:        d.Get("fixed_ip").(string),
		Name:           d.Get("name").(string),
	}

	networkInterface, err := client.NetworkInterface.CreateNetworkInterface(context.Background(), networkID, cnp)
	if err != nil {
		return fmt.Errorf("Error when create network interface: %v", err)
	}
	d.SetId(networkInterface.ID)
	return resourceBizFlyCloudNetworkInterfaceRead(d, meta)
}

func resourceBizFlyCloudNetworkInterfaceRead(d *schema.ResourceData, meta interface{}) error {
	if err := dataSourceBizFlyCloudNetworkInterfaceRead(d, meta); err != nil {
		return err
	}
	return nil
}

func resourceBizFlyCloudNetworkInterfacekUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()

	networkID := d.Get("network_id").(string)
	if networkID == "" {
		return fmt.Errorf("Invalid network id specified")
	}

	unp := &gobizfly.UpdateNetworkInterfacePayload{
		Name: d.Get("name").(string),
	}

	networkInterface, err := client.NetworkInterface.UpdateNetworkInterface(context.Background(), networkID, d.Id(), unp)
	if err != nil {
		return fmt.Errorf("Error when update network interface: %s, %v", d.Id(), err)
	}
	d.SetId(networkInterface.ID)
	return resourceBizFlyCloudNetworkInterfaceRead(d, meta)
}

func resourceBizFlyCloudNetworkInterfaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()

	networkID := d.Get("network_id").(string)
	if networkID == "" {
		return fmt.Errorf("Invalid network id specified")
	}

	err := client.NetworkInterface.DeleteNetworkInterface(context.Background(), networkID, d.Id())
	if err != nil {
		return fmt.Errorf("Error when delete network interface: %v", err)
	}
	return nil
}
