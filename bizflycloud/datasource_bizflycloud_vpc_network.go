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

	var matchVPC *gobizfly.VPC
	cidr := d.Get("cidr")

	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		var err error

		log.Printf("[DEBUG] Reading vpc network: %s", d.Id())
		vpcs, err := client.VPC.List(context.Background())

		// Retry on any API "not found" errors, but only on new resources.
		if d.IsNewResource() && errors.Is(err, gobizfly.ErrNotFound) {
			return resource.RetryableError(err)
		}
		if err != nil {
			return resource.NonRetryableError(err)
		}
		for _, vpc := range vpcs {
			if vpc.Subnets[0].CIDR == cidr {
				matchVPC = vpc
			}
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
	if matchVPC == nil {
		return fmt.Errorf("Error read vpc network (%s): empty response", d.Id())
	}

	d.SetId(matchVPC.ID)
	_ = d.Set("name", matchVPC.Name)
	_ = d.Set("description", matchVPC.Description)
	_ = d.Set("is_default", matchVPC.IsDefault)
	_ = d.Set("mtu", matchVPC.MTU)
	_ = d.Set("status", matchVPC.Status)
	_ = d.Set("created_at", matchVPC.CreatedAt)
	_ = d.Set("updated_at", matchVPC.UpdatedAt)
	_ = d.Set("cidr", matchVPC.Subnets[0].CIDR)
	_ = d.Set("mtu", matchVPC.MTU)

	if err := d.Set("availability_zones", readAvailabilityZones(matchVPC.AvailabilityZones)); err != nil {
		return fmt.Errorf("error setting availability_zones: %w", err)
	}

	if err := d.Set("subnets", readSubnets(matchVPC.Subnets)); err != nil {
		return fmt.Errorf("error setting subnets: %w", err)
	}

	return nil
}

func readAvailabilityZones(availabilityZone []string) []interface{} {
	var results []interface{}
	for _, az := range availabilityZone {
		results = append(results, az)
	}
	return results
}

func readSubnets(subnets []gobizfly.Subnet) []map[string]interface{} {
	var results []map[string]interface{}
	for _, s := range subnets {
		results = append(results, map[string]interface{}{
			"project_id":       s.ProjectID,
			"ip_version":       s.IPVersion,
			"gateway_ip":       s.GatewayIP,
			"allocation_pools": flattenAllocationPools(s.AllocationPools),
		})
	}
	return results
}

func flattenAllocationPools(allocationPools []gobizfly.AllocationPool) []map[string]interface{} {
	var flatAllocationPools []map[string]interface{}
	for _, p := range allocationPools {
		flatAllocationPools = append(flatAllocationPools, map[string]interface{}{
			"start": p.Start,
			"end":   p.End,
		})
	}
	return flatAllocationPools
}
