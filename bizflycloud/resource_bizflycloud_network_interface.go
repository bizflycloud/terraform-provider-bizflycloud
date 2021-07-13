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

	cnp := &gobizfly.CreateNetworkInterfacePayload{
		AttachedServer: d.Get("attached_server").(string),
		FixedIP:        d.Get("fixed_ip").(string),
		Name:           d.Get("name").(string),
		NetworkID:      d.Get("network_id").(string),
	}

	networkInterface, err := client.NetworkInterface.CreateNetworkInterface(context.Background(), cnp.NetworkID, cnp)
	if err != nil {
		return fmt.Errorf("Error when create vpc network: %v", err)
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
	return nil
}

func resourceBizFlyCloudNetworkInterfaceDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
