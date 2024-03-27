package bizflycloud

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bizflycloud/gobizfly"
	"github.com/bizflycloud/terraform-provider-bizflycloud/constants"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"k8s.io/apimachinery/pkg/util/wait"
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
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"action": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(constants.ValidL7PolicyActions, false),
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
				Default:  1,
			},
			"rules": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: getL7PolicyRuleSchema(),
				},
			},
		},
	}
}

// L7 policy rule schema
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
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(constants.ValidACLsTypes, false),
		},
		"compare_type": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(constants.ValidACLsCompareType, false),
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

// Create L7 policy
func resourceBizflycloudLoadbalancerL7PolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	listenerId := d.Get("listener_id").(string)
	_, err := waitListenerActiveProvisioningStatus(client, listenerId)
	if err != nil {
		log.Printf("[ERROR] wait listener active provisioning status failed: %v", err)
		return err
	}
	createReq, err := getCreateL7PolicyFromConfig(d)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] create l7 policy payload for listener %s: %#v", listenerId, *createReq)
	l7Policy, createErr := client.L7Policy.Create(context.Background(), listenerId, createReq)
	if createErr != nil {
		return fmt.Errorf("Create l7 policy for listener %s error: %v", listenerId, createErr)
	}
	d.SetId(l7Policy.Id)
	return resourceBizflycloudLoadbalancerL7PolicyRead(d, meta)
}

// Read L7 policy
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

// Update L7 policy
func resourceBizflycloudLoadbalancerL7PolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	listenerId := d.Get("listener_id").(string)
	_, err := waitListenerActiveProvisioningStatus(client, listenerId)
	if err != nil {
		log.Printf("[ERROR] wait listener active provisioning status failed: %v", err)
		return err
	}
	updateReq, err := getUpdateL7PolicyFromConfig(d)
	if err != nil {
		return err
	}
	rules, err := getUpdateL7PolicyRulesFromConfig(d, client)
	if err != nil {
		return err
	}
	updateReq.Rules = rules
	policyId := d.Id()
	_, err = waitL7PolicyActiveProvisioningStatus(client, policyId)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] update l7 policy %s payload for listener %s: %#v", policyId, listenerId, *updateReq)
	_, createErr := client.L7Policy.Update(context.Background(), policyId, updateReq)
	if createErr != nil {
		return fmt.Errorf("Update l7 policy %s for listener %s error: %v", policyId, listenerId, createErr)
	}
	return resourceBizflycloudLoadbalancerL7PolicyRead(d, meta)
}

// Destroy L7 policy
func resourceBizflycloudLoadbalancerL7PolicyDelete(d *schema.ResourceData, meta interface{}) error {
	policyId := d.Id()
	listenerId := d.Get("listener_id").(string)
	client := meta.(*CombinedConfig).gobizflyClient()
	_, err := waitListenerActiveProvisioningStatus(client, listenerId)
	if err != nil {
		log.Printf("[ERROR] wait listener active provisioning status failed: %v", err)
		return err
	}
	err = client.L7Policy.Delete(context.Background(), policyId)
	if err != nil {
		return fmt.Errorf("Error when delete l7 policy: %v", err)
	}
	return nil
}

// get create l7 policy payload
func getCreateL7PolicyFromConfig(d *schema.ResourceData) (*gobizfly.CreateL7PolicyRequest, error) {
	position := d.Get("position").(int)
	positionStr := fmt.Sprintf("%v", position)
	rules, err := getCreateL7PolicyRulesFromConfig(d)
	if err != nil {
		return nil, err
	}
	createReq := gobizfly.CreateL7PolicyRequest{
		Name:     d.Get("name").(string),
		Action:   d.Get("action").(string),
		Position: positionStr,
		Rules:    rules,
	}
	action := d.Get("action").(string)
	switch action {
	case constants.RedirectToUrlAction:
		redirectUrl := d.Get("redirect_url").(string)
		if redirectUrl == "" {
			return nil, fmt.Errorf("Action %s with 'redirect_url' must be not empty string", action)
		}
		createReq.RedirectUrl = redirectUrl
	case constants.RedirectToPoolAction:
		redirectPoolId := d.Get("redirect_pool_id").(string)
		if redirectPoolId == "" {
			return nil, fmt.Errorf("Action %s with 'redirect_pool_id' must be not empty string", action)
		}
		createReq.RedirectPoolId = redirectPoolId
	case constants.RedirectPrefixAction:
		redirectPrefix := d.Get("redirect_prefix").(string)
		if redirectPrefix == "" {
			return nil, fmt.Errorf("Action %s with 'redirect_prefix' must be not empty string", action)
		}
		createReq.RedirectPrefix = &redirectPrefix
	case constants.RejectAction:
	default:
		return nil, fmt.Errorf("Invalid l7 policy action %s.", action)
	}
	return &createReq, nil
}

// get create l7 policy rules payload
func getCreateL7PolicyRulesFromConfig(d *schema.ResourceData) ([]gobizfly.L7PolicyRuleRequest, error) {
	rules := make([]gobizfly.L7PolicyRuleRequest, 0)
	rulesLen := len(d.Get("rules").([]interface{}))
	for idx := 0; idx < rulesLen; idx++ {
		pattern := fmt.Sprintf("rules.%d.", idx)
		ruleType := d.Get(pattern + "type").(string)
		rule := gobizfly.L7PolicyRuleRequest{
			Invert:      d.Get(pattern + "invert").(bool),
			Type:        ruleType,
			CompareType: d.Get(pattern + "compare_type").(string),
		}

		ruleKey := d.Get(pattern + "key").(string)
		ruleValue := d.Get(pattern + "value").(string)
		err := validateL7PolicyRuleType(ruleType, ruleKey, ruleValue)
		if err != nil {
			return nil, fmt.Errorf("rules.%d: %v", idx, err)
		}
		rule.Key = ruleKey
		rule.Value = ruleValue
		rules = append(rules, rule)
	}
	return rules, nil
}

// validate l7 policy rule
func validateL7PolicyRuleType(ruleType, key, value string) error {
	var err error
	switch ruleType {
	case constants.ACLsTypeHostName:
		if key != "" {
			err = fmt.Errorf("Rule type %s with 'key' must have not value", ruleType)
		}
		if value == "" {
			err = fmt.Errorf("Rule type %s with 'value' must be not empty string", ruleType)
		}
	case constants.ACLsTypePath:
		if key != "" {
			err = fmt.Errorf("Rule type %s with 'key' must have not value", ruleType)
		}
		if value == "" {
			err = fmt.Errorf("Rule type %s with 'value' must be not empty string", ruleType)
		}
	case constants.ACLsTypeFileType:
		if key != "" {
			err = fmt.Errorf("Rule type %s with 'key' must have not value", ruleType)
		}
		if value == "" {
			err = fmt.Errorf("Rule type %s with 'value' must be not empty string", ruleType)
		}
	case constants.ACLsTypeHeader:
		if key == "" {
			err = fmt.Errorf("Rule type %s with 'key' must be not empty string", ruleType)
		}
		if value == "" {
			err = fmt.Errorf("Rule type %s with 'value' must be not empty string", ruleType)
		}
	default:
		err = fmt.Errorf("Invalid rule type %s", ruleType)
	}
	return err
}

// Update L7 policy
func getUpdateL7PolicyFromConfig(d *schema.ResourceData) (*gobizfly.UpdateL7PolicyRequest, error) {
	var err error
	updateReq := gobizfly.UpdateL7PolicyRequest{
		Name:     d.Get("name").(string),
		Action:   d.Get("action").(string),
		Position: d.Get("position").(int),
	}
	action := d.Get("action").(string)
	switch action {
	case constants.RedirectToUrlAction:
		redirectUrl := d.Get("redirect_url").(string)
		if redirectUrl == "" {
			return nil, fmt.Errorf("Action %s with 'redirect_url' must be not empty string", action)
		}
		updateReq.RedirectUrl = &redirectUrl
	case constants.RedirectToPoolAction:
		redirectPoolId := d.Get("redirect_pool_id").(string)
		if redirectPoolId == "" {
			return nil, fmt.Errorf("Action %s with 'redirect_pool_id' must be not empty string", action)
		}
		updateReq.RedirectPoolId = &redirectPoolId
	case constants.RedirectPrefixAction:
		redirectPrefix := d.Get("redirect_prefix").(string)
		if redirectPrefix == "" {
			return nil, fmt.Errorf("Action %s with 'redirect_prefix' must be not empty string", action)
		}
		updateReq.RedirectPrefix = &redirectPrefix
	case constants.RejectAction:
	default:
		return nil, fmt.Errorf("Invalid l7 policy action %s.", action)
	}

	if err != nil {
		return nil, err
	}
	return &updateReq, nil
}

// get update l7 policy rules payload
func getUpdateL7PolicyRulesFromConfig(d *schema.ResourceData, client *gobizfly.Client) ([]gobizfly.UpdateL7PolicyRuleRequest, error) {
	_, newRules := d.GetChange("rules")
	rulesReq, err := castUpdateL7PolicyRulesRequest(newRules)
	if err != nil {
		return nil, err
	}
	policyId := d.Id()
	updateRulesReq := make([]gobizfly.UpdateL7PolicyRuleRequest, 0)
	for _, ruleReq := range rulesReq {
		if ruleReq.ID != "" {
			updateRulesReq = append(updateRulesReq, ruleReq)
			continue
		}
		newRulePayload := gobizfly.L7PolicyRuleRequest{
			Invert:      ruleReq.Invert,
			Type:        ruleReq.Type,
			CompareType: ruleReq.CompareType,
			Key:         ruleReq.Key,
			Value:       ruleReq.Value,
		}
		log.Printf("[DEBUG] Create l7 policy payload: policyId=%s - payload=%#v", policyId, newRulePayload)
		newRule, err := client.L7Policy.CreateL7PolicyRule(context.Background(), policyId, newRulePayload)
		if err != nil {
			log.Printf("[Error] create l7 policy %s rule: %v", policyId, err)
			return nil, fmt.Errorf("Error create l7 policy %s rule: %v", policyId, err)
		}
		updateRuleReq := gobizfly.UpdateL7PolicyRuleRequest{
			ID:                  newRule.Id,
			L7PolicyRuleRequest: newRulePayload,
		}
		updateRulesReq = append(updateRulesReq, updateRuleReq)
	}
	return updateRulesReq, nil
}

func castUpdateL7PolicyRulesRequest(rules interface{}) ([]gobizfly.UpdateL7PolicyRuleRequest, error) {
	results := make([]gobizfly.UpdateL7PolicyRuleRequest, 0)
	castedRules := rules.([]interface{})
	for idx, rule := range castedRules {
		castedRule := rule.(map[string]interface{})
		ruleType := castedRule["type"].(string)
		result := gobizfly.UpdateL7PolicyRuleRequest{
			ID: castedRule["id"].(string),
			L7PolicyRuleRequest: gobizfly.L7PolicyRuleRequest{
				Invert:      castedRule["invert"].(bool),
				Type:        ruleType,
				CompareType: castedRule["compare_type"].(string),
			},
		}
		ruleKey := castedRule["key"].(string)
		ruleValue := castedRule["value"].(string)
		err := validateL7PolicyRuleType(ruleType, ruleKey, ruleValue)
		if err != nil {
			return nil, fmt.Errorf("rules.%d: %v", idx, err)
		}
		result.Key = ruleKey
		result.Value = ruleValue
		results = append(results, result)
	}
	return results, nil
}

// Parse l7 policy rules
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

// Check listener active status
func waitListenerActiveProvisioningStatus(client *gobizfly.Client, listenerId string) (*gobizfly.Listener, error) {
	var lb *gobizfly.Listener
	backoff := wait.Backoff{
		Duration: loadbalancerActiveInitDelay,
		Factor:   loadbalancerActiveFactor,
		Steps:    loadbalancerActiveSteps,
	}
	err := wait.ExponentialBackoff(backoff, func() (bool, error) {
		listener, err := client.Listener.Get(context.Background(), listenerId)
		if err != nil {
			return true, err
		}
		if listener.ProvisoningStatus == activeStatus {
			return true, nil
		} else if listener.ProvisoningStatus == errorStatus {
			return true, fmt.Errorf("Listener %s has gone into ERROR state", listenerId)
		} else {
			return false, nil
		}

	})
	if err != nil {
		if err == wait.ErrWaitTimeout {
			err = fmt.Errorf("Listener failed to go into ACTIVE provisioning status within allotted time")
		}
		return nil, err
	}

	return lb, err
}

// Check L7 policy active status
func waitL7PolicyActiveProvisioningStatus(client *gobizfly.Client, policyId string) (*gobizfly.DetailL7Policy, error) {
	var policy *gobizfly.DetailL7Policy
	backoff := wait.Backoff{
		Duration: loadbalancerActiveInitDelay,
		Factor:   loadbalancerActiveFactor,
		Steps:    loadbalancerActiveSteps,
	}
	err := wait.ExponentialBackoff(backoff, func() (bool, error) {
		policy, err := client.L7Policy.Get(context.Background(), policyId)
		if err != nil {
			return true, err
		}
		if policy.ProvisioningStatus == activeStatus {
			return true, nil
		} else if policy.ProvisioningStatus == errorStatus {
			return true, fmt.Errorf("L7 policy %s has gone into ERROR state", policyId)
		} else {
			return false, nil
		}

	})
	if err != nil {
		if err == wait.ErrWaitTimeout {
			err = fmt.Errorf("L7 policy failed to go into ACTIVE provisioning status within allotted time")
		}
		return nil, err
	}

	return policy, err
}
