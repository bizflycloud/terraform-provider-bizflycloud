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
	var matchNetworkInterface *gobizfly.NetworkInterface

	err := resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
		ipAddress := d.Get("ip_address").(string)
		networkInterfaces, err := client.NetworkInterface.List(context.Background(), &gobizfly.ListNetworkInterfaceOptions{})
		if d.IsNewResource() && errors.Is(err, gobizfly.ErrNotFound) {
			return resource.RetryableError(err)
		}
		if err != nil {
			return resource.NonRetryableError(err)
		}
		for _, networkInterface := range networkInterfaces {
			if networkInterface.FixedIps[0].IPAddress == ipAddress {
				matchNetworkInterface = networkInterface
			}
		}
		if matchNetworkInterface == nil {
			return resource.NonRetryableError(errors.New("no match network interface found"))
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
	if matchNetworkInterface == nil {
		return fmt.Errorf("Error read network interface network (%s): empty response", d.Id())
	}

	d.SetId(matchNetworkInterface.ID)
	_ = d.Set("name", matchNetworkInterface.Name)
	_ = d.Set("network_id", matchNetworkInterface.NetworkID)
	_ = d.Set("attached_server", matchNetworkInterface.AttachedServer)
	_ = d.Set("fixed_ip", matchNetworkInterface.FixedIps[0].IPAddress)
	_ = d.Set("status", matchNetworkInterface.Status)
	_ = d.Set("created_at", matchNetworkInterface.CreatedAt)
	_ = d.Set("updated_at", matchNetworkInterface.UpdatedAt)

	if err := d.Set("fixed_ips", readFixedIps(matchNetworkInterface.FixedIps)); err != nil {
		return fmt.Errorf("error setting fixed_ips: %w", err)
	}

	if err := d.Set("security_groups", readSecurityGroups(matchNetworkInterface.SecurityGroups)); err != nil {
		return fmt.Errorf("error setting security_groups: %w", err)
	}

	return nil
}

func readFixedIps(fixedIps []gobizfly.FixedIp) []map[string]interface{} {
	var results []map[string]interface{}
	for _, v := range fixedIps {
		results = append(results, map[string]interface{}{
			"subnet_id":  v.SubnetID,
			"ip_address": v.IPAddress,
		})
	}
	return results
}

func readSecurityGroups(securityGroups []string) []string {
	var results []string
	results = append(results, securityGroups...)

	return results
}
