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

func networkInterfaceRequestBuilder(d *schema.ResourceData) gobizfly.NetworkInterfaceRequestPayload {
	networkInterfaceOpts := gobizfly.NetworkInterfaceRequestPayload{}
	if v, ok := d.GetOk("name"); ok {
		networkInterfaceOpts.Name = v.(string)
	}
	if v, ok := d.GetOk("attached_server"); ok {
		networkInterfaceOpts.AttachedServer = v.(string)
	}
	if v, ok := d.GetOk("fixed_ip"); ok {
		networkInterfaceOpts.FixedIP = v.(string)
	}
	if v, ok := d.GetOk("action"); ok {
		networkInterfaceOpts.Action = v.(string)
	}
	return networkInterfaceOpts
}

func resourceBizFlyCloudNetworkInterfaceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()

	networkID := d.Get("network_id").(string)
	if networkID == "" {
		return fmt.Errorf("Invalid network id specified")
	}

	networkInterfaceOpts := networkInterfaceRequestBuilder(d)
	networkInterface, err := client.NetworkInterface.CreateNetworkInterface(context.Background(), networkID, &networkInterfaceOpts)
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

func resourceBizFlyCloudNetworkInterfaceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()

	networkID := d.Get("network_id").(string)
	if networkID == "" {
		return fmt.Errorf("Invalid network id specified")
	}

	unp := &gobizfly.NetworkInterfaceRequestPayload{
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
