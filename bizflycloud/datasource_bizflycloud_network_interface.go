package bizflycloud

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/bizflycloud/gobizfly"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceBizFlyCloudNetworkInterface() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceBizFlyCloudNetworkInterfaceRead,
		Schema: dataNetworkInterfaceSchema(),
	}
}

func dataSourceBizFlyCloudNetworkInterfaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()

	gnp := &gobizfly.GetNetworkInterfacePayload{
		NetworkID: d.Get("network_id").(string),
	}

	if v, ok := d.GetOk("id"); ok {
		d.SetId(v.(string))
	}

	var networkInterface *gobizfly.NetworkInterface

	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		var err error

		log.Printf("[DEBUG] Reading network interface: %s", d.Id())
		networkInterface, err = client.NetworkInterface.GetNetworkInterface(context.Background(), gnp.NetworkID, d.Id())

		// Retry on any API "not found" errors, but only on new resources.
		if d.IsNewResource() && errors.Is(err, gobizfly.ErrNotFound) {
			return resource.RetryableError(err)
		}
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})

	// Prevent confusing Terraform error messaging to operators by
	// Only ignoring API "not found" errors if not a new resource
	if !d.IsNewResource() && errors.Is(err, gobizfly.ErrNotFound) {
		log.Printf("[WARN] network interface %s is not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("Error read network interface network %s: %w", d.Id(), err)
	}

	// Prevent panics.
	if networkInterface == nil {
		return fmt.Errorf("Error read network interface network (%s): empty response", d.Id())
	}

	d.SetId(networkInterface.ID)
	_ = d.Set("name", networkInterface.Name)
	_ = d.Set("attached_server", networkInterface.AttachedServer)
	_ = d.Set("fixed_ip", networkInterface.FixedIps)

	return nil
}
