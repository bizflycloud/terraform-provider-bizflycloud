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

	origins := d.Get("origin").(*schema.Set).List()
	origin := origins[0].(map[string]interface{})

	cdp := &gobizfly.CreateDomainPayload{
		Domain: d.Get("domain").(string),
		Origin: &gobizfly.Origin{
			Name:          origin["name"].(string),
			UpstreamHost:  origin["upstream_host"].(string),
			UpstreamProto: origin["upstream_proto"].(string),
			UpstreamAddrs: origin["upstream_addrs"].(string),
		},
	}
	cdr, err := client.CDN.Create(context.Background(), cdp)
	if err != nil {
		return fmt.Errorf("error when create cdn resource: %v", err)
	}
	domain := cdr.Domain
	d.SetId(domain.DomainID)
	return resourceBizflyCloudCDNRead(d, meta)
}

func resourceBizflyCloudCDNRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	domain, err := client.CDN.Get(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("error when get cdn resource: %v", err)
	}
	_ = d.Set("domain_cdn", domain.DomainCDN)
	_ = d.Set("domain_id", domain.DomainID)

	return nil
}

func resourceBizflyCloudCDNUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	domain, err := client.CDN.Get(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("error when get cdn resource: %v", err)
	}

	if d.HasChange("origin") {
		origins := d.Get("origin").(*schema.Set).List()
		origin := origins[0].(map[string]interface{})
		udp := &gobizfly.UpdateDomainPayload{
			Origin: &gobizfly.Origin{
				Name:          origin["name"].(string),
				UpstreamHost:  origin["upstream_host"].(string),
				UpstreamProto: origin["upstream_proto"].(string),
				UpstreamAddrs: origin["upstream_addrs"].(string),
			},
		}
		_, err := client.CDN.Update(context.Background(), domain.DomainID, udp)
		if err != nil {
			return fmt.Errorf("error when update cdn resource: %v", err)
		}
	}
	return resourceBizflyCloudCDNRead(d, meta)
}

func resourceBizflyCloudCDNDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	err := client.CDN.Delete(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("error when delete cdn resource : %v", err)
	}
	return nil
}
