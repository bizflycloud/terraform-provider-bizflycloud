// This file is part of terraform-provider-bizflycloud
//
// Copyright (C) 2021  Bizfly Cloud
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>

package bizflycloud

import (
	"context"
	"fmt"

	"github.com/bizflycloud/gobizfly"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBizflyCloudFirewall() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Create: resourceBizflyCloudFirewallCreate,
		Read:   resourceBizflyCloudFirewallRead,
		Update: resourceBizflyCloudFirewallUpdate,
		Delete: resourceBizflyCloudFirewallDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rules_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"network_interface_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ingress": {
				Type:       schema.TypeSet,
				Optional:   true,
				Computed:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: firewallRuleSchema(),
				},
			},
			"egress": {
				Type:       schema.TypeSet,
				Optional:   true,
				Computed:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: firewallRuleSchema(),
				},
			},
			"network_interfaces": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
			},
		},
	}
}

func firewallRuleSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cidr": {
			Type:     schema.TypeString,
			Required: true,
		},
		"protocol": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"port_range": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
}

func firewallRequestBuilder(d *schema.ResourceData) gobizfly.FirewallRequestPayload {
	firewallOpts := gobizfly.FirewallRequestPayload{}
	if v, ok := d.GetOk("name"); ok {
		firewallOpts.Name = v.(string)
	}
	if v, ok := d.GetOk("network_interfaces"); ok {
		for _, id := range v.(*schema.Set).List() {
			firewallOpts.NetworkInterfaces = append(firewallOpts.NetworkInterfaces, id.(string))
		}
	}
	if v, ok := d.GetOk("ingress"); ok {
		firewallOpts.InBound = flatternFirewallRules(v.(*schema.Set))
	}
	if v, ok := d.GetOk("egress"); ok {
		firewallOpts.OutBound = flatternFirewallRules(v.(*schema.Set))
	}
	return firewallOpts
}
func resourceBizflyCloudFirewallCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()

	firewallOpts := firewallRequestBuilder(d)
	firewall, err := client.CloudServer.Firewalls().Create(context.Background(), &firewallOpts)
	if err != nil {
		return fmt.Errorf("Error when creating firewall: %v", err)
	}
	d.SetId(firewall.BaseFirewall.ID)
	return resourceBizflyCloudFirewallRead(d, meta)
}

func resourceBizflyCloudFirewallRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	firewall, err := client.CloudServer.Firewalls().Get(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error when retrieving fireall: %v", err)
	}
	_ = d.Set("name", firewall.BaseFirewall.Name)
	_ = d.Set("rules_count", firewall.BaseFirewall.RulesCount)
	_ = d.Set("network_interface_count", firewall.BaseFirewall.NetworkInterfaceCount)

	_ = d.Set("network_interfaces", flatternBizflyCloudNetworkInterfaces(firewall.NetworkInterface))
	if len(firewall.InBound) > 0 {
		_ = d.Set("ingress", convertFWRule(firewall.InBound))
	}
	if len(firewall.OutBound) > 0 {
		_ = d.Set("egress", convertFWRule(firewall.OutBound))
	}
	return nil
}

func resourceBizflyCloudFirewallDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	_, err := client.CloudServer.Firewalls().Delete(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error when deleting firewall: %v", err)
	}
	return nil
}

func resourceBizflyCloudFirewallUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()

	firewallOpts := firewallRequestBuilder(d)
	firewall, err := client.CloudServer.Firewalls().Update(context.Background(), d.Id(), &firewallOpts)
	if err != nil {
		return fmt.Errorf("Error when creating firewall: %v", err)
	}
	d.SetId(firewall.BaseFirewall.ID)
	return resourceBizflyCloudFirewallRead(d, meta)
}

func flatternFirewallRules(rules *schema.Set) []gobizfly.FirewallRuleCreateRequest {
	fwrules := []gobizfly.FirewallRuleCreateRequest{}
	for _, rawRule := range rules.List() {
		r := rawRule.(map[string]interface{})
		rule := gobizfly.FirewallRuleCreateRequest{
			Type:      "CUSTOM",
			CIDR:      r["cidr"].(string),
			PortRange: r["port_range"].(string),
			Protocol:  r["protocol"].(string),
		}
		fwrules = append(fwrules, rule)
	}
	return fwrules
}

func convertFWRule(rules []gobizfly.FirewallRule) []map[string]interface{} {
	result := make([]map[string]interface{}, len(rules))
	for i, v := range rules {
		result[i] = map[string]interface{}{
			"cidr":       v.CIDR,
			"port_range": v.PortRange,
			"protocol":   v.Protocol,
		}
	}
	return result
}

func flatternBizflyCloudNetworkInterfaces(networkInterfaces []*gobizfly.NetworkInterface) *schema.Set {
	flattenedNetworkInterfaces := schema.NewSet(schema.HashString, []interface{}{})
	for _, networkInterface := range networkInterfaces {
		flattenedNetworkInterfaces.Add(networkInterface.ID)
	}
	return flattenedNetworkInterfaces
}
