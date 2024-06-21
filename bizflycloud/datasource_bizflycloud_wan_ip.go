package bizflycloud

import (
	"context"
	"errors"
	"fmt"
	"github.com/bizflycloud/gobizfly"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

func dataSourceBizflyCloudWanIP() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceBizflyCloudWanIPRead,
		Schema: dataSourceWanIPSchema(),
	}
}
func dataSourceBizflyCloudWanIPRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	var matchWanIP *gobizfly.CloudServerPublicNetworkInterface

	err := resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
		ipAddress := d.Get("ip_address").(string)
		wanIPs, err := client.CloudServer.PublicNetworkInterfaces().List(context.Background())
		if d.IsNewResource() && errors.Is(err, gobizfly.ErrNotFound) {
			return resource.RetryableError(err)
		}
		if err != nil {
			return resource.NonRetryableError(err)
		}
		for _, wanIP := range wanIPs {
			if wanIP.IPAddress == ipAddress {
				matchWanIP = wanIP
			}
		}
		if matchWanIP == nil {
			return resource.NonRetryableError(errors.New("no wan ip found"))
		}
		return nil
	})
	if !d.IsNewResource() && errors.Is(err, gobizfly.ErrNotFound) {
		log.Printf("[WARN] WAN IP %s is not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error read WAN IP network %s: %w", d.Id(), err)
	}

	d.SetId(matchWanIP.ID)
	_ = d.Set("name", matchWanIP.Name)
	_ = d.Set("network_id", matchWanIP.NetworkID)
	_ = d.Set("ip_address", matchWanIP.IPAddress)
	_ = d.Set("ip_version", matchWanIP.IpVersion)
	_ = d.Set("status", matchWanIP.Status)
	_ = d.Set("created_at", matchWanIP.CreatedAt)
	_ = d.Set("updated_at", matchWanIP.UpdatedAt)
	_ = d.Set("security_groups", readSecurityGroups(matchWanIP.SecurityGroups))
	_ = d.Set("billing_type", matchWanIP.BillingType)
	_ = d.Set("bandwidth", matchWanIP.Bandwidth)
	_ = d.Set("availability_zone", matchWanIP.AvailabilityZone)
	return nil
}
