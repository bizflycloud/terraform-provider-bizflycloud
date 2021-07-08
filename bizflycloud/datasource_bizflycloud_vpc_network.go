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

func dataSourceBizFlyCloudVPCNetwork() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceBizFlyCloudVPCNetworkRead,
		Schema: dataVPCNetworkSchema(),
	}
}

func dataSourceBizFlyCloudVPCNetworkRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()

	if v, ok := d.GetOk("id"); ok {
		d.SetId(v.(string))
	}

	var network *gobizfly.VPC

	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		var err error

		log.Printf("[DEBUG] Reading vpc network: %s", d.Id())
		network, err = client.VPC.Get(context.Background(), d.Id())

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
		log.Printf("[WARN] vpc network %s is not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("Error read vpc network %s: %w", d.Id(), err)
	}

	// Prevent panics.
	if network == nil {
		return fmt.Errorf("Error read vpc network (%s): empty response", d.Id())
	}

	d.SetId(network.ID)
	_ = d.Set("name", network.Name)
	_ = d.Set("description", network.Description)
	_ = d.Set("is_default", network.IsDefault)

	return nil
}
