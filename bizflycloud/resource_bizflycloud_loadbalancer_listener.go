package bizflycloud

import (
	"context"
	"fmt"
	"time"

	"github.com/bizflycloud/gobizfly"
	"github.com/bizflycloud/terraform-provider-bizflycloud/constants"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceBizflyCloudLoadBalancerListener() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Create: resourceBizflyCloudLoadBalancerListenerCreate,
		Read:   resourceBizflyCloudLoadBalancerListenerRead,
		Update: resourceBizflyCloudLoadBalancerListenerUpdate,
		Delete: resourceBizflyCloudLoadBalancerListenerDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(constants.ValidListenerProtocols, false),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"default_pool_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"default_tls_ref": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"load_balancer_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"listener_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  5000,
			},
			"server_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  5000,
			},
			"server_connect_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  5000,
			},
			"operating_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"provisioning_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"l7policy_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
		},
	}
}

func resourceBizflyCloudLoadBalancerListenerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	lbID := d.Get("load_balancer_id").(string)

	// Lock to serialize operations on the same load balancer
	mutex := getLoadBalancerMutex(lbID)
	mutex.Lock()
	defer mutex.Unlock()

	_, _ = waitLoadbalancerActiveProvisioningStatus(client, lbID, loadbalancerResource)
	lName := d.Get("name").(string)
	lPoolDefaultID := d.Get("default_pool_id").(string)
	lPoolTLSRef := d.Get("default_tls_ref").(string)
	listenerTimeout := d.Get("listener_timeout").(int)
	serverTimeout := d.Get("server_timeout").(int)
	serverConnectTimeout := d.Get("server_connect_timeout").(int)
	lcr := gobizfly.CloudLoadBalancerListenerCreateRequest{
		Name:                   &lName,
		Protocol:               d.Get("protocol").(string),
		ProtocolPort:           d.Get("port").(int),
		DefaultPoolID:          &lPoolDefaultID,
		DefaultTLSContainerRef: &lPoolTLSRef,
		TimeoutClientData:      &listenerTimeout,
		TimeoutMemberData:      &serverTimeout,
		TimeoutMemberConnect:   &serverConnectTimeout,
	}
	listener, err := client.CloudLoadBalancer.Listeners().Create(context.Background(), lbID, &lcr)
	if err != nil {
		return fmt.Errorf("error creating listener: %v", err)
	}
	if listener == nil {
		return fmt.Errorf("error creating listener for loadbalancer %s: listener object is nil", lbID)
	}
	d.SetId(listener.ID)
	return resourceBizflyCloudLoadBalancerListenerRead(d, meta)
}

func resourceBizflyCloudLoadBalancerListenerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	listener, err := client.CloudLoadBalancer.Listeners().Get(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("error retrieving listener: %v", err)
	}
	if listener == nil {
		return fmt.Errorf("error retrieving listener %s: listener object is nil", d.Id())
	}
	if len(listener.LoadBalancers) == 0 {
		return fmt.Errorf("error retrieving listener %s: listener has no load balancers", d.Id())
	}
	l7policyIDs := make([]string, 0)
	for _, policy := range listener.L7Policies {
		l7policyIDs = append(l7policyIDs, policy.ID)
	}
	_ = d.Set("name", listener.Name)
	_ = d.Set("protocol", listener.Protocol)
	_ = d.Set("port", listener.ProtocolPort)
	_ = d.Set("description", listener.Description)
	_ = d.Set("default_pool_id", listener.DefaultPoolID)
	_ = d.Set("default_tls_ref", listener.DefaultTLSContainerRef)
	_ = d.Set("load_balancer_id", listener.LoadBalancers[0].ID)
	_ = d.Set("listener_timeout", listener.TimeoutClientData)
	_ = d.Set("server_timeout", listener.TimeoutMemberData)
	_ = d.Set("server_connect_timeout", listener.TimeoutMemberConnect)
	_ = d.Set("operating_status", listener.OperatingStatus)
	_ = d.Set("provisioning_status", listener.ProvisoningStatus)
	_ = d.Set("l7policy_ids", l7policyIDs)
	_ = d.Set("created_at", listener.CreatedAt)
	_ = d.Set("updated_at", listener.UpdatedAt)

	return nil
}

func resourceBizflyCloudLoadBalancerListenerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	lbID := d.Get("load_balancer_id").(string)

	// Lock to serialize operations on the same load balancer
	mutex := getLoadBalancerMutex(lbID)
	mutex.Lock()
	defer mutex.Unlock()

	_, _ = waitLoadbalancerActiveProvisioningStatus(client, lbID, loadbalancerResource)
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	tlsRef := d.Get("default_tls_ref").(string)
	poolID := d.Get("default_pool_id").(string)
	listenerTimeout := d.Get("listener_timeout").(int)
	serverTimeout := d.Get("server_timeout").(int)
	serverConnectTimeout := d.Get("server_connect_timeout").(int)
	lur := gobizfly.CloudLoadBalancerListenerUpdateRequest{
		Name:                   &name,
		Description:            &description,
		DefaultTLSContainerRef: &tlsRef,
		DefaultPoolID:          &poolID,
		TimeoutClientData:      &listenerTimeout,
		TimeoutMemberData:      &serverTimeout,
		TimeoutMemberConnect:   &serverConnectTimeout,
	}
	_, err := client.CloudLoadBalancer.Listeners().Update(context.Background(), d.Id(), &lur)
	if err != nil {
		return fmt.Errorf("error updating listener: %v", err)
	}
	return resourceBizflyCloudLoadBalancerListenerRead(d, meta)
}

func resourceBizflyCloudLoadBalancerListenerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	lbID := d.Get("load_balancer_id").(string)

	// Lock to serialize operations on the same load balancer
	mutex := getLoadBalancerMutex(lbID)
	mutex.Lock()
	defer mutex.Unlock()

	_, _ = waitLoadbalancerActiveProvisioningStatus(client, lbID, loadbalancerResource)
	err := client.CloudLoadBalancer.Listeners().Delete(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("error deleting listener: %v", err)
	}
	return nil
}
