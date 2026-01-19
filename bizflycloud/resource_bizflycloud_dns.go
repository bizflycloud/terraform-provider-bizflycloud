package bizflycloud

import (
	"context"
	"fmt"
	"time"

	"github.com/bizflycloud/gobizfly"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBizflyCloudDNS() *schema.Resource {
	return &schema.Resource{
		Create: resourceBizflyCloudDNSCreate,
		Read:   resourceBizflyCloudDNSRead,
		Update: resourceBizflyCloudDNSUpdate,
		Delete: resourceBizflyCloudDNSDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		SchemaVersion: 1,
		Schema:        resourceDNSSchema(),
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceBizflyCloudDNSCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	czp := &gobizfly.CreateZonePayload{
		Name:        d.Get("name").(string),
		Required:    d.Get("required").(bool),
		Description: d.Get("description").(string),
	}
	zone, err := client.DNS.CreateZone(context.Background(), czp)
	if err != nil {
		return fmt.Errorf("error creating dns zone: %v", err)
	}
	d.SetId(zone.ID)
	return resourceBizflyCloudDNSRead(d, meta)
}

func resourceBizflyCloudDNSRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	zone, err := client.DNS.GetZone(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("error getting dns zone: %v", err)
	}
	_ = d.Set("name", zone.Name)
	_ = d.Set("active", zone.Active)
	_ = d.Set("created_at", zone.CreatedAt)
	_ = d.Set("update_at", zone.UpdatedAt)
	_ = d.Set("deleted", zone.Deleted)
	_ = d.Set("ttl", zone.TTL)
	_ = d.Set("tenant_id", zone.TenantID)

	if err := d.Set("nameserver", readNameServer(zone.NameServer)); err != nil {
		return fmt.Errorf("error setting nameserver: %w", err)
	}

	if err := d.Set("record_set", readRecordsSet(zone.RecordsSet)); err != nil {
		return fmt.Errorf("error setting record_set: %w", err)
	}

	return nil
}

func resourceBizflyCloudDNSUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceBizflyCloudDNSDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	err := client.DNS.DeleteZone(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("error deleting dns zone: %v", err)
	}
	return nil
}

func readNameServer(nameServer []string) []string {
	var results []string
	results = append(results, nameServer...)

	return results
}

func readRecordsSet(recordsSet []gobizfly.Record) []map[string]interface{} {
	var results []map[string]interface{}
	for _, v := range recordsSet {
		results = append(results, map[string]interface{}{
			"id":   v.ID,
			"name": v.Name,
			"type": v.Type,
			"ttl":  v.TTL,
		})
	}
	return results
}
