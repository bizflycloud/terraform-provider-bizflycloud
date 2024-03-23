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

func resourceBizflyCloudLoadBalancerPool() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Create: resourceBizflyCloudLoadBalancerPoolCreate,
		Update: resourceBizflyCloudLoadBalancerPoolUpdate,
		Read:   resourceBizflyCloudLoadBalancerPoolRead,
		Delete: resourceBizflyCloudLoadBalancerPoolDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"algorithm": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(constants.ValidAlgorithms, false),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(constants.ValidPoolProtocols, false),
			},
			"load_balancer_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"members": {
				Type:       schema.TypeList,
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
					Schema: getHealthMonitorSchema(),
				},
			},
			"persistent": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: getPersistentSchema(),
				},
			},
		},
	}
}

func loadbalancerMemberSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
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
		"network_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"backup": {
			Type:     schema.TypeBool,
			Required: false,
			Optional: true,
		},
		"operating_status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"provisioning_status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"subnet_id": {
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

func getHealthMonitorSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(constants.ValidMemberProtocols, false),
		},
		"timeout": {
			Type:     schema.TypeInt,
			Default:  5,
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
			Default:  5,
			Optional: true,
		},
		"http_method": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(constants.ValidHealthMonitorMethods, false),
		},
		"url_path": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "/",
		},
		"expected_code": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "200",
		},
		"operating_status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"provisioning_status": {
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

func getPersistentSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(constants.ValidStickySessions, false),
		},
		"cookie_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
}

func resourceBizflyCloudLoadBalancerPoolCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	lbID := d.Get("load_balancer_id").(string)
	_, err := waitLoadbalancerActiveProvisioningStatus(client, lbID, loadbalancerResource)
	if err != nil {
		return err
	}
	createReq := getCreatePoolPayloadFromConfig(d)
	log.Printf("[DEBUG] Create pool payload: %+v", createReq)
	pool, err := client.Pool.Create(context.Background(), lbID, &createReq)
	if err != nil {
		return fmt.Errorf("Error when create pool for loadbalancer %s: %+v", lbID, err)
	}
	_, err = waitPoolActiveProvisioningStatus(client, pool.ID)
	if err != nil {
		log.Printf("[ERROR] wait pool %s active provisioning status failed: %+v", pool.ID, err)
		return err
	}
	d.SetId(pool.ID)
	return resourceBizflyCloudLoadBalancerPoolRead(d, meta)
}

func resourceBizflyCloudLoadBalancerPoolRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	pool, err := client.Pool.Get(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error when retrieving load balancer pool %s: %v", d.Id(), err)
	}
	networkNameMap, err := getNetworkNameMap(client)
	log.Printf("[DEBUG] Network name mapping: %+v", networkNameMap)
	if err != nil {
		return fmt.Errorf("Error when list vpc networks: %v", err)
	}
	healthMonitor := convertHealthMonitor(pool.HealthMonitor)
	persistent := convertSessionPersistent(pool.SessionPersistence)
	members := convertMember(pool.Members, networkNameMap)

	_ = d.Set("name", pool.Name)
	_ = d.Set("algorithm", pool.LBAlgorithm)
	_ = d.Set("description", pool.Description)
	_ = d.Set("protocol", pool.Protocol)
	_ = d.Set("load_balancer_id", pool.LoadBalancers[0].ID)
	_ = d.Set("health_monitor", healthMonitor)
	_ = d.Set("persistent", persistent)
	_ = d.Set("members", members)
	return nil
}

func resourceBizflyCloudLoadBalancerPoolUpdate(d *schema.ResourceData, meta interface{}) error {
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
	return resourceBizflyCloudLoadBalancerPoolRead(d, meta)
}

func resourceBizflyCloudLoadBalancerPoolDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	lbID := d.Get("load_balancer_id").(string)
	_, _ = waitLoadbalancerActiveProvisioningStatus(client, lbID, loadbalancerResource)
	err := client.Pool.Delete(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error when deleting load balancer pool: %v", err)
	}
	return nil
}

func getCreatePoolPayloadFromConfig(d *schema.ResourceData) gobizfly.PoolCreateRequest {
	poolName := d.Get("name").(string)
	healthMonitorReq := getHealthMonitorPayloadFromConfig(d)
	poolMembersReq := getMembersPaylaodFromConfig(d)
	stickySessionReq := getSessionPersistentPayloadFromConfig(d)
	poolReq := gobizfly.PoolCreateRequest{
		Name:               &poolName,
		LBAlgorithm:        d.Get("algorithm").(string),
		Protocol:           d.Get("protocol").(string),
		HealthMonitor:      healthMonitorReq,
		Members:            poolMembersReq,
		SessionPersistence: stickySessionReq,
	}
	return poolReq
}

func getSessionPersistentPayloadFromConfig(d *schema.ResourceData) *gobizfly.SessionPersistence {
	sessionPersistents := d.Get("persistent").([]interface{})
	if len(sessionPersistents) == 0 {
		return nil
	}
	sessionPersistent := sessionPersistents[0].(map[string]interface{})
	sesstionType := sessionPersistent["type"].(string)
	sessionPersistentReq := gobizfly.SessionPersistence{
		Type: sesstionType,
	}
	if sesstionType == constants.AppCookie {
		cookieName := sessionPersistent["cookie_name"].(string)
		sessionPersistentReq.CookieName = &cookieName
	}
	return &sessionPersistentReq
}

func getMembersPaylaodFromConfig(d *schema.ResourceData) []gobizfly.PoolMemberRequest {
	poolMembers := d.Get("members").([]interface{})
	if len(poolMembers) == 0 {
		return nil
	}

	poolMembersReq := make([]gobizfly.PoolMemberRequest, 0)
	for idx, member := range poolMembers {
		castedMember := member.(map[string]interface{})
		memberReq := gobizfly.PoolMemberRequest{
			ID:          idx,
			Name:        castedMember["name"].(string),
			Address:     castedMember["address"].(string),
			Weight:      castedMember["weight"].(int),
			Port:        castedMember["protocol_port"].(int),
			NetworkName: castedMember["network_name"].(string),
		}
		poolMembersReq = append(poolMembersReq, memberReq)
	}
	return poolMembersReq
}

func getHealthMonitorPayloadFromConfig(d *schema.ResourceData) *gobizfly.PoolHealthMonitorRequest {
	healthMonitors := d.Get("health_monitor").([]interface{})
	if len(healthMonitors) == 0 {
		return nil
	}
	healthMonitor := healthMonitors[0].(map[string]interface{})
	healthMonitorReq := gobizfly.PoolHealthMonitorRequest{
		Name:           healthMonitor["name"].(string),
		Type:           healthMonitor["type"].(string),
		Timeout:        healthMonitor["timeout"].(int),
		Delay:          healthMonitor["delay"].(int),
		MaxRetries:     healthMonitor["max_retries"].(int),
		MaxRetriesDown: healthMonitor["max_retries_down"].(int),
		HttpMethod:     healthMonitor["http_method"].(string),
		UrlPath:        healthMonitor["url_path"].(string),
		ExpectedCodes:  healthMonitor["expected_code"].(string),
	}
	return &healthMonitorReq
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

func convertMember(members []gobizfly.Member, networkNameMap map[string]string) []map[string]interface{} {
	result := make([]map[string]interface{}, len(members))
	for i, v := range members {
		networkName := networkNameMap[v.SubnetID]
		result[i] = map[string]interface{}{
			"id":                  v.ID,
			"name":                v.Name,
			"weight":              v.Weight,
			"address":             v.Address,
			"protocol_port":       v.ProtocolPort,
			"network_name":        networkName,
			"backup":              v.Backup,
			"operating_status":    v.OperatingStatus,
			"provisioning_status": v.ProvisoningStatus,
			"subnet_id":           v.SubnetID,
			"project_id":          v.ProjectID,
			"created_at":          v.CreatedAt,
			"updated_at":          v.UpdatedAt,
		}
	}
	return result
}

func convertHealthMonitor(healthMonitor *gobizfly.HealthMonitor) []map[string]interface{} {
	results := make([]map[string]interface{}, 0, 1)
	if healthMonitor == nil {
		return nil
	}
	result := map[string]interface{}{
		"id":                  healthMonitor.ID,
		"name":                healthMonitor.Name,
		"type":                healthMonitor.Type,
		"timeout":             healthMonitor.TimeOut,
		"max_retries":         healthMonitor.MaxRetries,
		"max_retries_down":    healthMonitor.MaxRetriesDown,
		"delay":               healthMonitor.Delay,
		"http_method":         healthMonitor.HTTPMethod,
		"url_path":            healthMonitor.UrlPath,
		"expected_code":       healthMonitor.ExpectedCodes,
		"operating_status":    healthMonitor.OperatingStatus,
		"provisioning_status": healthMonitor.ProvisioningStatus,
		"created_at":          healthMonitor.CreatedAt,
		"updated_at":          healthMonitor.UpdatedAt,
	}
	results = append(results, result)
	return results
}

func convertSessionPersistent(sessionPersistent *gobizfly.SessionPersistence) []map[string]interface{} {
	results := make([]map[string]interface{}, 0, 1)
	if sessionPersistent == nil {
		return nil
	}
	result := map[string]interface{}{
		"type":        sessionPersistent.Type,
		"cookie_name": sessionPersistent.CookieName,
	}
	results = append(results, result)
	return results
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

// Check loadbalancer pool active status
func waitPoolActiveProvisioningStatus(client *gobizfly.Client, poolID string) (*gobizfly.Pool, error) {
	var pool *gobizfly.Pool
	backoff := wait.Backoff{
		Duration: loadbalancerActiveInitDelay,
		Factor:   loadbalancerActiveFactor,
		Steps:    loadbalancerActiveSteps,
	}
	err := wait.ExponentialBackoff(backoff, func() (bool, error) {
		pool, err := client.Pool.Get(context.Background(), poolID)
		if err != nil {
			return true, err
		}
		if pool.ProvisoningStatus == activeStatus {
			return true, nil
		} else if pool.ProvisoningStatus == errorStatus {
			return true, fmt.Errorf("Loadbalancer Pool %s has gone into ERROR state", poolID)
		} else {
			return false, nil
		}

	})
	if err != nil {
		if err == wait.ErrWaitTimeout {
			err = fmt.Errorf("Loadbalancer Pool %s failed to go into ACTIVE provisioning status within allotted time", poolID)
		}
		return nil, err
	}

	return pool, err
}

func getNetworkNameMap(client *gobizfly.Client) (map[string]string, error) {
	networkNameMap := make(map[string]string, 0)
	vpcNetworks, err := client.VPC.List(context.Background())
	if err != nil {
		return nil, err
	}
	for _, vpc := range vpcNetworks {
		vpcNetworkName := vpc.Name
		for _, subnet := range vpc.Subnets {
			subnetID := subnet.ID
			networkNameMap[subnetID] = vpcNetworkName
		}
	}

	return networkNameMap, nil
}
