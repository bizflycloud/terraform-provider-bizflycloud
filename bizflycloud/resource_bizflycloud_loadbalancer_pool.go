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
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
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

// Get member schema
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
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      1,
			ValidateFunc: validation.IntBetween(0, 256),
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

// Get health monitor schema
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
			ValidateFunc: validation.StringInSlice(constants.ValidHealthMonitorProtocols, false),
		},
		"timeout": {
			Type:     schema.TypeInt,
			Default:  5,
			Optional: true,
		},
		"max_retries": {
			Type:         schema.TypeInt,
			Default:      3,
			Optional:     true,
			ValidateFunc: validation.IntBetween(1, 10),
		},
		"max_retries_down": {
			Type:         schema.TypeInt,
			Default:      3,
			Optional:     true,
			ValidateFunc: validation.IntBetween(1, 10),
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
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringMatch(constants.ValidURLPathRegex, "url_path must start with '/'"),
		},
		"expected_code": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(constants.ValidHealthMonitorExceptedCodes, false),
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

// Get persistent schema
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

// Create pool resource
func resourceBizflyCloudLoadBalancerPoolCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	lbID := d.Get("load_balancer_id").(string)

	// Lock to serialize operations on the same load balancer
	mutex := getLoadBalancerMutex(lbID)
	mutex.Lock()
	defer mutex.Unlock()

	_, err := waitLoadbalancerActiveProvisioningStatus(client, lbID, loadbalancerResource)
	if err != nil {
		return err
	}
	createReq := getCreatePoolPayloadFromConfig(d)
	log.Printf("[DEBUG] Create pool payload: %+v", createReq)
	pool, err := client.CloudLoadBalancer.Pools().Create(context.Background(), lbID, &createReq)
	if err != nil {
		return fmt.Errorf("Error when create pool for loadbalancer %s: %+v", lbID, err)
	}
	if pool == nil {
		/// need retry here to get the pool object
		return fmt.Errorf("Error when create pool for loadbalancer %s: pool object is nil", lbID)
	}
	poolID := pool.ID
	err = checkLoadbalancerPoolActiveStatus(client, poolID, lbID)
	if err != nil {
		return err
	}
	d.SetId(poolID)

	// Update description if could not create pool with description
	description := d.Get("description").(string)
	if description != pool.Description {
		updateReq := &gobizfly.CloudLoadBalancerPoolUpdateRequest{
			Description: &description,
		}
		persistent := getSessionPersistentPayloadFromConfig(d)
		updateReq.SessionPersistence = persistent
		_, err = client.CloudLoadBalancer.Pools().Update(context.Background(), poolID, updateReq)
		if err != nil {
			log.Printf("[ERROR] Update pool %s failed: %v", poolID, err)
			return err
		}

		err = checkLoadbalancerPoolActiveStatus(client, poolID, lbID)
		if err != nil {
			return err
		}
	}

	// Create health monitor
	healthMonitorPayload := getCreateHealthMonitorPayloadFromConfig(d)
	if healthMonitorPayload != nil {
		_, err = client.CloudLoadBalancer.HealthMonitors().Create(context.Background(), poolID, healthMonitorPayload)
		if err != nil {
			log.Printf("[ERROR] Create health monitor for pool %s failed: %+v", poolID, err)
			return err
		}
		err = checkLoadbalancerPoolActiveStatus(client, poolID, lbID)
		if err != nil {
			return err
		}
	}
	return resourceBizflyCloudLoadBalancerPoolRead(d, meta)
}

// Read pool resource
func resourceBizflyCloudLoadBalancerPoolRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	pool, err := client.CloudLoadBalancer.Pools().Get(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error when retrieving load balancer pool %s: %v", d.Id(), err)
	}
	if pool == nil {
		return fmt.Errorf("Error when retrieving load balancer pool %s: pool object is nil", d.Id())
	}
	if len(pool.LoadBalancers) == 0 {
		return fmt.Errorf("Error when retrieving load balancer pool %s: pool has no load balancers", d.Id())
	}
	healthMonitor := convertHealthMonitor(pool.HealthMonitor)
	persistent := convertSessionPersistent(pool.SessionPersistence)
	members := convertMember(pool.Members)

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

// Change pool resource
func resourceBizflyCloudLoadBalancerPoolUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	poolID := d.Id()
	lbID := d.Get("load_balancer_id").(string)

	// Lock to serialize operations on the same load balancer
	mutex := getLoadBalancerMutex(lbID)
	mutex.Lock()
	defer mutex.Unlock()

	err := checkLoadbalancerPoolActiveStatus(client, poolID, lbID)
	if err != nil {
		return err
	}
	// Update pool
	updateReq := getUpdatePoolPayloadFromConfig(d)
	_, err = client.CloudLoadBalancer.Pools().Update(context.Background(), poolID, &updateReq)
	if err != nil {
		log.Printf("[ERROR] Update pool %s failed: %+v", poolID, err)
		return err
	}
	err = checkLoadbalancerPoolActiveStatus(client, poolID, lbID)
	if err != nil {
		return err
	}

	// update health monitor of pool
	if d.HasChange("health_monitor") {
		healthMonitorID, healthMonitorReq := getUpdateHealthMonitorFromConfig(d)
		if healthMonitorReq != nil {
			if healthMonitorID == "" {
				// Create health monitor for pool
				createHealthMonitorReq := getCreateHealthMonitorPayloadFromConfig(d)
				_, err = client.CloudLoadBalancer.HealthMonitors().Create(context.Background(), poolID, createHealthMonitorReq)
				if err != nil {
					log.Printf("[ERROR] Create health monitor %s failed: %v", healthMonitorID, err)
					return err
				}
			} else {
				// Update health monitor for pool
				_, err = client.CloudLoadBalancer.HealthMonitors().Update(context.Background(), healthMonitorID, healthMonitorReq)
				if err != nil {
					log.Printf("[ERROR] Update health monitor %s failed: %v", healthMonitorID, err)
					return err
				}
			}

			err = checkLoadbalancerPoolActiveStatus(client, poolID, lbID)
			if err != nil {
				return err
			}
		}
	}

	// Update merbers of pool
	if d.HasChange("members") {
		oldMembers, newMembers := d.GetChange("members")
		// Delete all current members of pool
		for _, member := range oldMembers.([]interface{}) {
			castedMember := member.(map[string]interface{})
			memberID := castedMember["id"].(string)
			if err != nil {
				log.Printf("[ERROR] check pool active status before delete member")
				return err
			}
			err = client.CloudLoadBalancer.Members().Delete(context.Background(), poolID, memberID)
			if err != nil {
				return fmt.Errorf("Error when delete old member %s: %v", memberID, err)
			}

			err = checkLoadbalancerPoolActiveStatus(client, poolID, lbID)
			if err != nil {
				return err
			}
		}
		// Create new members for pool
		for _, member := range newMembers.([]interface{}) {
			castedMember := member.(map[string]interface{})
			updateMemberReq := &gobizfly.CloudLoadBalancerMemberCreateRequest{
				Name:         castedMember["name"].(string),
				Weight:       castedMember["weight"].(int),
				Address:      castedMember["address"].(string),
				ProtocolPort: castedMember["protocol_port"].(int),
				Backup:       castedMember["backup"].(bool),
			}
			_, err = client.CloudLoadBalancer.Members().Create(context.Background(), poolID, updateMemberReq)
			if err != nil {
				return fmt.Errorf("Error when create new member: %v", err)
			}

			err = checkLoadbalancerPoolActiveStatus(client, poolID, lbID)
			if err != nil {
				return err
			}
		}
	}
	return resourceBizflyCloudLoadBalancerPoolRead(d, meta)
}

// Delete pool resource
func resourceBizflyCloudLoadBalancerPoolDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	lbID := d.Get("load_balancer_id").(string)

	// Lock to serialize operations on the same load balancer
	mutex := getLoadBalancerMutex(lbID)
	mutex.Lock()
	defer mutex.Unlock()

	_, _ = waitLoadbalancerActiveProvisioningStatus(client, lbID, loadbalancerResource)
	err := client.CloudLoadBalancer.Pools().Delete(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error when deleting load balancer pool: %v", err)
	}
	return nil
}

// Get create pool payload from config
func getCreatePoolPayloadFromConfig(d *schema.ResourceData) gobizfly.CloudLoadBalancerPoolCreateRequest {
	poolName := d.Get("name").(string)
	poolMembersReq := getMembersPayloadFromConfig(d)
	stickySessionReq := getSessionPersistentPayloadFromConfig(d)
	poolReq := gobizfly.CloudLoadBalancerPoolCreateRequest{
		Name:               &poolName,
		LBAlgorithm:        d.Get("algorithm").(string),
		Protocol:           d.Get("protocol").(string),
		Members:            poolMembersReq,
		SessionPersistence: stickySessionReq,
	}
	return poolReq
}

// Get update pool payload from config
func getUpdatePoolPayloadFromConfig(d *schema.ResourceData) gobizfly.CloudLoadBalancerPoolUpdateRequest {
	poolReq := gobizfly.CloudLoadBalancerPoolUpdateRequest{}
	if d.HasChange("name") {
		_, newName := d.GetChange("name")
		castedNewName := newName.(string)
		poolReq.Name = &castedNewName
	}
	if d.HasChange("description") {
		_, newDescription := d.GetChange("description")
		castedNewDescription := newDescription.(string)
		poolReq.Description = &castedNewDescription
	}
	if d.HasChange("algorithm") {
		_, newAlgorithm := d.GetChange("algorithm")
		castedNewAlgorithm := newAlgorithm.(string)
		poolReq.LBAlgorithm = &castedNewAlgorithm
	}
	persistent := getSessionPersistentPayloadFromConfig(d)
	poolReq.SessionPersistence = persistent
	return poolReq
}

// Get create health monitor payload from config
func getCreateHealthMonitorPayloadFromConfig(d *schema.ResourceData) *gobizfly.CloudLoadBalancerHealthMonitorCreateRequest {
	healthMonitors := d.Get("health_monitor").([]interface{})
	if len(healthMonitors) == 0 {
		return nil
	}
	healthMonitor := healthMonitors[0].(map[string]interface{})
	healthMonitorType := healthMonitor["type"].(string)
	healthMonitorReq := gobizfly.CloudLoadBalancerHealthMonitorCreateRequest{
		Name:           healthMonitor["name"].(string),
		Type:           healthMonitorType,
		TimeOut:        healthMonitor["timeout"].(int),
		PoolID:         d.Id(),
		Delay:          healthMonitor["delay"].(int),
		MaxRetries:     healthMonitor["max_retries"].(int),
		MaxRetriesDown: healthMonitor["max_retries_down"].(int),
	}
	if (healthMonitorType != constants.TcpProtocol) &&
		(healthMonitorType != constants.UdpConnectProtocol) {

		healthMonitorReq.HTTPMethod = healthMonitor["http_method"].(string)
		healthMonitorReq.URLPath = healthMonitor["url_path"].(string)
		healthMonitorReq.ExpectedCodes = healthMonitor["expected_code"].(string)
	}
	return &healthMonitorReq
}

// Get update health monitor payload from config
func getUpdateHealthMonitorFromConfig(d *schema.ResourceData) (string, *gobizfly.CloudLoadBalancerHealthMonitorUpdateRequest) {
	healthMonitor := d.Get("health_monitor").([]interface{})
	healthMonitorLen := len(healthMonitor)
	if healthMonitorLen == 0 {
		return "", nil
	}

	healthMonitorMap := healthMonitor[0].(map[string]interface{})
	updateReq := gobizfly.CloudLoadBalancerHealthMonitorUpdateRequest{
		Name: healthMonitorMap["name"].(string),
	}
	timeoutInt := healthMonitorMap["timeout"].(int)
	if timeoutInt != 0 {
		updateReq.TimeOut = &timeoutInt
	}
	delayInt := healthMonitorMap["delay"].(int)
	if delayInt != 0 {
		updateReq.Delay = &delayInt
	}
	maxRetriesInt := healthMonitorMap["max_retries"].(int)
	if maxRetriesInt != 0 {
		updateReq.MaxRetries = &maxRetriesInt
	}
	maxRetriesDownInt := healthMonitorMap["max_retries_down"].(int)
	if maxRetriesDownInt != 0 {
		updateReq.MaxRetriesDown = &maxRetriesDownInt
	}

	healthMonitorType := healthMonitorMap["type"].(string)
	if (healthMonitorType != constants.TcpProtocol) &&
		(healthMonitorType != constants.UdpConnectProtocol) {

		httpMethodStr := healthMonitorMap["http_method"].(string)
		if httpMethodStr != "" {
			updateReq.HTTPMethod = &httpMethodStr
		}
		urlPathStr := healthMonitorMap["url_path"].(string)
		if urlPathStr != "" {
			updateReq.URLPath = &urlPathStr
		}
		expectedCodeStr := healthMonitorMap["expected_code"].(string)
		if expectedCodeStr != "" {
			updateReq.ExpectedCodes = &expectedCodeStr
		}
	}

	poolHealthMonitorID := healthMonitorMap["id"].(string)
	log.Printf("[DEBUG] Update health monitor %s payload: %+v", poolHealthMonitorID, updateReq)
	return poolHealthMonitorID, &updateReq
}

// Get session persistent payload from config
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

// Get members payload from config
func getMembersPayloadFromConfig(d *schema.ResourceData) []gobizfly.CloudLoadBalancerPoolMemberRequest {
	poolMembers := d.Get("members").([]interface{})
	if len(poolMembers) == 0 {
		return nil
	}

	poolMembersReq := make([]gobizfly.CloudLoadBalancerPoolMemberRequest, 0)
	for idx, member := range poolMembers {
		castedMember := member.(map[string]interface{})
		memberReq := gobizfly.CloudLoadBalancerPoolMemberRequest{
			ID:      idx,
			Name:    castedMember["name"].(string),
			Address: castedMember["address"].(string),
			Weight:  castedMember["weight"].(int),
			Port:    castedMember["protocol_port"].(int),
		}
		poolMembersReq = append(poolMembersReq, memberReq)
	}
	return poolMembersReq
}

// Convert members to update state
func convertMember(members []gobizfly.CloudLoadBalancerMember) []map[string]interface{} {
	result := make([]map[string]interface{}, len(members))
	for i, v := range members {
		result[i] = map[string]interface{}{
			"id":                  v.ID,
			"name":                v.Name,
			"weight":              v.Weight,
			"address":             v.Address,
			"protocol_port":       v.ProtocolPort,
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

// Convert health monitor to update state
func convertHealthMonitor(healthMonitor *gobizfly.CloudLoadBalancerHealthMonitor) []map[string]interface{} {
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
		"url_path":            healthMonitor.URLPath,
		"expected_code":       healthMonitor.ExpectedCodes,
		"operating_status":    healthMonitor.OperatingStatus,
		"provisioning_status": healthMonitor.ProvisioningStatus,
		"created_at":          healthMonitor.CreatedAt,
		"updated_at":          healthMonitor.UpdatedAt,
	}
	results = append(results, result)
	return results
}

// Convert session persistent to update state
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

// Wait loadbalancer pool active status
func waitPoolActiveProvisioningStatus(client *gobizfly.Client, poolID string) (*gobizfly.CloudLoadBalancerPool, error) {
	var pool *gobizfly.CloudLoadBalancerPool
	backoff := wait.Backoff{
		Duration: loadbalancerActiveInitDelay,
		Factor:   loadbalancerActiveFactor,
		Steps:    loadbalancerActiveSteps,
	}
	err := wait.ExponentialBackoff(backoff, func() (bool, error) {
		pool, err := client.CloudLoadBalancer.Pools().Get(context.Background(), poolID)
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

// Check loadbalancer and pool resource active status
func checkLoadbalancerPoolActiveStatus(client *gobizfly.Client, poolID, loadbalancerID string) error {
	_, err := waitPoolActiveProvisioningStatus(client, poolID)
	if err != nil {
		log.Printf("[ERROR] wait pool %s active status: %v", poolID, err)
		return err
	}
	_, err = waitLoadbalancerActiveProvisioningStatus(client, loadbalancerID, loadbalancerResource)
	if err != nil {
		log.Printf("[ERROR] wait loadbalancer %s active status: %v", loadbalancerID, err)
		return err
	}
	return nil
}
