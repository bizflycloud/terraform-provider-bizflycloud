package bizflycloud

import (
	"context"
	"fmt"
	"time"

	"github.com/bizflycloud/gobizfly"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBizflyCloudCDN() *schema.Resource {
	return &schema.Resource{
		Create: resourceBizflyCloudCDNCreate,
		Read:   resourceBizflyCloudCDNRead,
		Update: resourceBizflyCloudCDNUpdate,
		Delete: resourceBizflyCloudCDNDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		SchemaVersion: 1,
		Schema:        resourceCDNSchema(),
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceBizflyCloudCDNCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	czp := &gobizfly.CreateZonePayload{
		Name:        d.Get("name").(string),
		Required:    d.Get("required").(bool),
		Description: d.Get("description").(string),
	}
	zone, err := client.DNS.CreateZone(context.Background(), czp)
	if err != nil {
		return fmt.Errorf("Error when create dns zone: %v", err)
	}
	d.SetId(zone.ID)
	return resourceBizflyCloudDNSRead(d, meta)
}

func resourceBizflyCloudCDNRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	zone, err := client.DNS.GetZone(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error when get dns zone: %v", err)
	}
	_ = d.Set("name", zone.Name)
	_ = d.Set("active", zone.Active)
	_ = d.Set("created_at", zone.CreatedAt)
	_ = d.Set("update_at", zone.UpdatedAt)
	_ = d.Set("deleted", zone.Deleted)
	_ = d.Set("ttl", zone.TTL)
	_ = d.Set("tenant_id", zone.TenantId)

	if err := d.Set("nameserver", readNameServer(zone.NameServer)); err != nil {
		return fmt.Errorf("error setting nameserver: %w", err)
	}

	if err := d.Set("record_set", readRecordsSet(zone.RecordsSet)); err != nil {
		return fmt.Errorf("error setting record_set: %w", err)
	}

	return nil
}

func resourceBizflyCloudCDNUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceBizflyCloudCDNDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	err := client.CDN.Delete(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error when delete cdn resource : %v", err)
	}
	return nil
}
