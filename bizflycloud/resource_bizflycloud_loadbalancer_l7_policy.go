package bizflycloud

import (
	"context"
	"fmt"
	"log"

	"github.com/bizflycloud/gobizfly"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBizflyCloudLoadBalancerL7Policy() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Create: resourceBizflycloudLoadbalancerL7PolicyCreate,
		Read:   resourceBizflycloudLoadbalancerL7PolicyRead,
		Update: resourceBizflycloudLoadbalancerL7PolicyUpdate,
		Delete: resourceBizflycloudLoadbalancerL7PolicyDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"action": {
				Type:     schema.TypeString,
				Required: true,
			},
			"redirect_pool_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"redirect_prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"redirect_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"listener_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"position": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"rules": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: getL7PolicyRuleSchema(),
				},
			},
		},
	}
}

func getL7PolicyRuleSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"invert": {
			Type:     schema.TypeBool,
			Required: true,
		},
		"type": {
			Type:     schema.TypeString,
			Required: true,
		},
		"compare_type": {
			Type:     schema.TypeString,
			Required: true,
		},
		"key": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"value": {
			Type:     schema.TypeString,
			Required: true,
		},
		"operating_status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"provisioning_status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"project_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"created_at": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"updated_at": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func resourceBizflycloudLoadbalancerL7PolicyCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceBizflycloudLoadbalancerL7PolicyRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*CombinedConfig).gobizflyClient()
	policyId := d.Id()
	log.Printf("[DEBUG] test read l7 policy %s", policyId)
	l7Policy, err := client.L7Policy.Get(context.Background(), policyId)
	if err != nil {
		return fmt.Errorf("Error when retrieving l7 policy %s: %v", policyId, err)
	}
	l7PolicyRules, err := client.L7Policy.ListL7PolicyRules(context.Background(), policyId)
	if err != nil {
		return fmt.Errorf("Error when listing l7 policy %s rules: %v", policyId, err)
	}
	rules := parseL7PolicyRules(l7PolicyRules)
	_ = d.Set("name", l7Policy.Name)
	_ = d.Set("action", l7Policy.Action)
	_ = d.Set("redirect_pool_id", l7Policy.RedirectPoolId)
	_ = d.Set("redirect_prefix", l7Policy.RedirectPrefix)
	_ = d.Set("redirect_url", l7Policy.RedirectUrl)
	_ = d.Set("listener_id", l7Policy.ListenerId)
	_ = d.Set("position", l7Policy.Position)
	_ = d.Set("rules", rules)
	return nil
}

func resourceBizflycloudLoadbalancerL7PolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceBizflycloudLoadbalancerL7PolicyDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func parseL7PolicyRules(rules []gobizfly.DetailL7PolicyRule) []map[string]interface{} {
	results := make([]map[string]interface{}, 0)
	for _, rule := range rules {
		result := map[string]interface{}{
			"id":                  rule.Id,
			"invert":              rule.Invert,
			"type":                rule.Type,
			"compare_type":        rule.CompareType,
			"key":                 rule.Key,
			"value":               rule.Value,
			"operating_status":    rule.OperatingStatus,
			"provisioning_status": rule.ProvisioningStatus,
			"project_id":          rule.ProjectId,
			"created_at":          rule.CreatedAt,
			"updated_at":          rule.UpdatedAt,
		}
		results = append(results, result)
	}
	return results
}
