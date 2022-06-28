package bizflycloud

import (
	"context"
	"fmt"
	"time"

	"github.com/bizflycloud/gobizfly"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBizFlyCloudLoadBalancerPool() *schema.Resource {
	return &schema.Resource{
		Create: resourceBizFlyCloudLoadBalancerPoolCreate,
		Update: resourceBizFlyCloudLoadBalancerPoolUpdate,
		Read:   resourceBizFlyCloudLoadBalancerPoolRead,
		Delete: resourceBizFlyCloudLoadBalancerPoolDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"algorithm": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Required: true,
			},
			"load_balancer_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"members": {
				Type:       schema.TypeSet,
				Optional:   true,
				Computed:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: loadbalancerMemberSchema(),
				},
			},
			"health_monitor": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Default:  "pool-monitor",
							Optional: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							//	TODO Validate type
						},
						"timeout": {
							Type:     schema.TypeInt,
							Default:  3,
							Optional: true,
						},
						"max_retries": {
							Type:     schema.TypeInt,
							Default:  3,
							Optional: true,
						},
						"max_retries_down": {
							Type:     schema.TypeInt,
							Default:  3,
							Optional: true,
						},
						"delay": {
							Type:     schema.TypeInt,
							Default:  3,
							Optional: true,
						},
						"http_method": {
							Type:     schema.TypeString,
							Optional: true,
							//	TODO validate http method
						},
						"url_path": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"expected_code": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"persistent": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"cookie_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}
func resourceBizFlyCloudLoadBalancerPoolCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*CombinedConfig).gobizflyClient()
	lbID := d.Get("load_balancer_id").(string)
	_, _ = waitLoadbalancerActiveProvisioningStatus(client, lbID, loadbalancerResource)
	poolName := d.Get("name").(string)

	pcr := gobizfly.PoolCreateRequest{
		Name:        &poolName,
		LBAlgorithm: d.Get("algorithm").(string),
		Protocol:    d.Get("protocol").(string),
	}
	// Get session persistent
	sessionPersistent := d.Get("persistent").([]interface{})
	if len(sessionPersistent) == 1 {
		sp := sessionPersistent[0].(map[string]interface{})
		spType := sp["type"].(string)
		var sess gobizfly.SessionPersistence
		sess.Type = spType
		if spType == "APP_COOKIE" {
			cookieName := sp["cookie_name"].(string)
			sess.CookieName = &cookieName
		}
		pcr.SessionPersistence = &sess
	}
	pool, err := client.Pool.Create(context.Background(), lbID, &pcr)
	if err != nil {
		return fmt.Errorf("Error when creating load balancer pool: %v", err)
	}
	d.SetId(pool.ID)
	_, _ = waitLoadbalancerActiveProvisioningStatus(client, pool.ID, poolResource)
	// create member
	if v, ok := d.GetOk("members"); ok {
		mcr := flatternMembers(v.(*schema.Set))
		for _, m := range mcr {
			_, _ = waitLoadbalancerActiveProvisioningStatus(client, pool.ID, poolResource)
			_, err := client.Member.Create(context.Background(), pool.ID, &m)
			if err != nil {
				return fmt.Errorf("Error when creating member %s: %v", m.Address, err)
			}
		}
	}
	// create health monitor
	err = createHealthMonitor(client, d.Get("health_monitor").([]interface{}), pool.ID)
	if err != nil {
		return err
	}

	return resourceBizFlyCloudLoadBalancerPoolRead(d, meta)
}

func resourceBizFlyCloudLoadBalancerPoolUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	var pur gobizfly.PoolUpdateRequest
	poolChanged := false
	if d.HasChange("persistent") {
		poolChanged = true
		// Get session persistent
		sessionPersistent := d.Get("persistent").([]interface{})
		if len(sessionPersistent) == 1 {
			sp := sessionPersistent[0].(map[string]interface{})
			spType := sp["type"].(string)
			var sess gobizfly.SessionPersistence
			sess.Type = spType
			if spType == "APP_COOKIE" {
				cookieName := sp["cookie_name"].(string)
				sess.CookieName = &cookieName
			}
			pur.SessionPersistence = &sess
		}
	}
	if d.HasChange("algorithm") {
		poolChanged = true
		algo := d.Get("algorithm").(string)
		pur.LBAlgorithm = &algo
	}
	if d.HasChange("description") {
		poolChanged = true
		desc := d.Get("description").(string)
		pur.Description = &desc
	}
	if poolChanged {
		_, _ = waitLoadbalancerActiveProvisioningStatus(client, d.Id(), poolResource)
		_, err := client.Pool.Update(context.Background(), d.Id(), &pur)
		if err != nil {
			return fmt.Errorf("Error when update pool %s, %v", d.Id(), err)
		}
	}
	if d.HasChange("health_monitor") {
		healthMonitors := d.Get("health_monitor").([]interface{})
		if len(healthMonitors) == 1 {
			_, _ = waitLoadbalancerActiveProvisioningStatus(client, d.Id(), poolResource)
			healthMonitor := healthMonitors[0].(map[string]interface{})
			// if health monitor is not created, create a new one
			if healthMonitor["id"].(string) == "" {
				err := createHealthMonitor(client, d.Get("health_monitor").([]interface{}), d.Id())
				if err != nil {
					return err
				}
			} else {
				hm, err := client.HealthMonitor.Get(context.Background(), healthMonitor["id"].(string))
				if err != nil {
					return fmt.Errorf("Error when get current health monitor: %s, %v", healthMonitor["id"].(string), err)
				}
				hmur := gobizfly.HealthMonitorUpdateRequest{
					Name:           healthMonitor["name"].(string),
					TimeOut:        healthMonitor["timeout"].(*int),
					MaxRetries:     healthMonitor["max_retries"].(*int),
					MaxRetriesDown: healthMonitor["max_retries_down"].(*int),
					Delay:          healthMonitor["delay"].(*int),
				}
				if hm.Type == "HTTP" {
					hmur.HTTPMethod = healthMonitor["http_method"].(*string)
					hmur.URLPath = healthMonitor["url_path"].(*string)
					hmur.ExpectedCodes = healthMonitor["expected_code"].(*string)
				}

				_, err = client.HealthMonitor.Update(context.Background(), healthMonitor["id"].(string), &hmur)
				if err != nil {
					return fmt.Errorf("Error when updating health monitor: %s, %v", healthMonitor["name"].(string), err)
				}
			}
		}
	}
	if d.HasChange("members") {
		// update member
		if v, ok := d.GetOk("members"); ok {
			mcr := flatternMembers(v.(*schema.Set))
			// get current members in load balancer
			currentMembers, err := client.Member.List(context.Background(), d.Id(), &gobizfly.ListOptions{})
			// workaround remove all member then re-add
			// TODO use batch update member when the api is available
			if err != nil {
				return fmt.Errorf("Error when get current member: %v", err)
			}
			for _, member := range currentMembers {
				_, _ = waitLoadbalancerActiveProvisioningStatus(client, d.Id(), poolResource)
				err := client.Member.Delete(context.Background(), d.Id(), member.ID)
				if err != nil {
					return fmt.Errorf("Error when delete old member: %v", err)
				}
			}
			for _, m := range mcr {
				_, _ = waitLoadbalancerActiveProvisioningStatus(client, d.Id(), poolResource)
				_, err := client.Member.Create(context.Background(), d.Id(), &m)
				if err != nil {
					return fmt.Errorf("Error when creating member %s: %v", m.Address, err)
				}
			}
		}
	}
	return resourceBizFlyCloudLoadBalancerPoolRead(d, meta)
}
func resourceBizFlyCloudLoadBalancerPoolRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	pool, err := client.Pool.Get(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error when retrieving load balancer pool: %v", err)
	}
	_ = d.Set("name", pool.Name)
	_ = d.Set("algorithm", pool.LBAlgorithm)
	_ = d.Set("description", pool.Description)
	_ = d.Set("protocol", pool.Protocol)
	_ = d.Set("load_balancer_id", pool.LoadBalancers[0].ID)
	members, err := client.Member.List(context.Background(), pool.ID, &gobizfly.ListOptions{})
	if err != nil {
		return fmt.Errorf("Error when getting pool member of pool: %s, %v", pool.Name, err)
	}
	if len(members) > 0 {
		_ = d.Set("members", convertMember(members))
	}
	_ = d.Set("health_monitor", convertHealthMonitor(pool.HealthMonitor))
	_ = d.Set("persistent", convertSessionPersistent(pool.SessionPersistence))
	return nil
}

func resourceBizFlyCloudLoadBalancerPoolDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	lbID := d.Get("load_balancer_id").(string)
	_, _ = waitLoadbalancerActiveProvisioningStatus(client, lbID, loadbalancerResource)
	err := client.Pool.Delete(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error when deleting load balancer pool: %v", err)
	}
	return nil
}

func loadbalancerMemberSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"weight": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  1,
		},
		"address": {
			Type:     schema.TypeString,
			Required: true,
		},
		"protocol_port": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"backup": {
			Type:     schema.TypeBool,
			Required: false,
			Optional: true,
		},
	}
}

func flatternMembers(rules *schema.Set) []gobizfly.MemberCreateRequest {
	members := []gobizfly.MemberCreateRequest{}
	for _, rawMember := range rules.List() {
		m := rawMember.(map[string]interface{})
		member := gobizfly.MemberCreateRequest{
			Name:         m["name"].(string),
			Address:      m["address"].(string),
			ProtocolPort: m["protocol_port"].(int),
			Weight:       m["weight"].(int),
			Backup:       m["backup"].(bool),
		}
		members = append(members, member)
	}
	return members
}

func convertMember(members []*gobizfly.Member) []map[string]interface{} {
	result := make([]map[string]interface{}, len(members))
	for i, v := range members {
		result[i] = map[string]interface{}{
			"name":          v.Name,
			"address":       v.Address,
			"protocol_port": v.ProtocolPort,
			"weight":        v.Weight,
			"backup":        v.Backup,
		}
	}
	return result
}

func convertHealthMonitor(healthMonitor *gobizfly.HealthMonitor) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, 1)
	if healthMonitor == nil {
		return nil
	}
	r := make(map[string]interface{})
	r["id"] = healthMonitor.ID
	r["name"] = healthMonitor.Name
	r["type"] = healthMonitor.Type
	r["timeout"] = healthMonitor.TimeOut
	r["max_retries"] = healthMonitor.MaxRetries
	r["max_retries_down"] = healthMonitor.MaxRetriesDown
	r["delay"] = healthMonitor.Delay
	r["http_method"] = healthMonitor.HTTPMethod
	r["expected_code"] = healthMonitor.ExpectedCodes
	result = append(result, r)
	return result
}

func convertSessionPersistent(sessionPersistent *gobizfly.SessionPersistence) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, 1)
	if sessionPersistent == nil {
		return nil
	}
	r := make(map[string]interface{})
	r["type"] = sessionPersistent.Type
	r["cookie_name"] = sessionPersistent.CookieName
	result = append(result, r)
	return result
}

func createHealthMonitor(client *gobizfly.Client, healthMonitor []interface{}, poolID string) error {
	if len(healthMonitor) == 1 {
		hm := healthMonitor[0].(map[string]interface{})
		hmType := hm["type"].(string)
		var hmcr gobizfly.HealthMonitorCreateRequest
		if hmType == "TCP" {
			hmcr.Name = hm["name"].(string)
			hmcr.Type = hmType
			hmcr.TimeOut = hm["timeout"].(int)
			hmcr.MaxRetriesDown = hm["max_retries_down"].(int)
			hmcr.MaxRetries = hm["max_retries"].(int)
			hmcr.Delay = hm["delay"].(int)
		} else {
			hmcr.HTTPMethod = hm["http_method"].(string)
			hmcr.URLPath = hm["url_path"].(string)
			hmcr.ExpectedCodes = hm["expected_code"].(string)
		}
		_, _ = waitLoadbalancerActiveProvisioningStatus(client, poolID, poolResource)
		_, err := client.HealthMonitor.Create(context.Background(), poolID, &hmcr)
		if err != nil {
			return fmt.Errorf("Error when creating health monitor for pool: %s, %v", poolID, err)
		}
	}
	return nil
}
