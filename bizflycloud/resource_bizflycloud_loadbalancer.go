package bizflycloud

import (
	"context"
	"fmt"
	"time"

	"github.com/bizflycloud/gobizfly"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"k8s.io/apimachinery/pkg/util/wait"
)

const (
	loadbalancerActiveInitDelay = 1 * time.Second
	loadbalancerActiveFactor    = 1.2
	loadbalancerActiveSteps     = 19

	activeStatus = "ACTIVE"
	errorStatus  = "ERROR"

	loadbalancerResource = "loadbalancer"
	listenerResource     = "listener"
	poolResource         = "pool"
)

func resourceBizFlyCloudLoadBalancer() *schema.Resource {
	return &schema.Resource{
		Create: resourceBizFlyCloudLoadBalancerCreate,
		Read:   resourceBizFlyCloudLoadBalancerRead,
		Update: resourceBizFlyCloudLoadBalancerUpdate,
		Delete: resourceBizFlyCloudLoadBalancerDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"network_type": {
				Type:     schema.TypeString,
				Default:  "external",
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Default:  "medium",
				Optional: true,
			},
			"vip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"provisioning_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"operating_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pools": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
			},
			"listeners": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceBizFlyCloudLoadBalancerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	lbcr := gobizfly.LoadBalancerCreateRequest{
		Name:        d.Get("name").(string),
		Type:        d.Get("type").(string),
		NetworkType: d.Get("network_type").(string),
		Description: d.Get("description").(string),
	}
	lb, err := client.LoadBalancer.Create(context.Background(), &lbcr)
	if err != nil {
		return fmt.Errorf("Error when creating load balancer")
	}
	d.SetId(lb.ID)

	// TODO add pool and listener

	return resourceBizFlyCloudLoadBalancerRead(d, meta)
}

func resourceBizFlyCloudLoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	lb, err := waitLoadbalancerActiveProvisioningStatus(client, d.Id(), loadbalancerResource)
	if err != nil {
		return fmt.Errorf("Error when retrieving load balancer: %v", err)
	}
	d.Set("name", lb.Name)
	d.Set("description", lb.Description)
	d.Set("network_type", lb.NetworkType)
	d.Set("type", lb.Type)
	d.Set("vip_address", lb.VipAddress)
	d.Set("provisioning_status", lb.ProvisioningStatus)
	d.Set("operating_status", lb.OperatingStatus)
	pools := schema.NewSet(schema.HashString, []interface{}{})
	for _, v := range lb.Pools {
		pools.Add(v.ID)
	}

	if err := d.Set("pools", pools); err != nil {
		return fmt.Errorf("Error setting pools: %v", err)
	}
	listeners := schema.NewSet(schema.HashString, []interface{}{})
	for _, v := range lb.Pools {
		listeners.Add(v.ID)
	}
	if err := d.Set("listeners", listeners); err != nil {
		return fmt.Errorf("Error setting listeners: %v", err)
	}
	return nil
}

func resourceBizFlyCloudLoadBalancerUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("type") {
		// TODO update new type of load balancer
	}
	if d.HasChange("listeners") {
		//	TODO update listeners in a load balancer
	}
	if d.HasChange("pools") {
		//	TODO update pool in a load balancer
	}
	return nil
}

func resourceBizFlyCloudLoadBalancerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	lb, err := waitLoadbalancerActiveProvisioningStatus(client, d.Id(), loadbalancerResource)
	if err != nil {
		return fmt.Errorf("Error when retrieving load balancer: %v", err)
	}
	ldr := gobizfly.LoadBalancerDeleteRequest{
		ID:      lb.ID,
		Cascade: true,
	}
	_ = client.LoadBalancer.Delete(context.Background(), &ldr)
	return nil
}

func waitLoadbalancerActiveProvisioningStatus(client *gobizfly.Client, ID string, resourceType string) (*gobizfly.LoadBalancer, error) {
	backoff := wait.Backoff{
		Duration: loadbalancerActiveInitDelay,
		Factor:   loadbalancerActiveFactor,
		Steps:    loadbalancerActiveSteps,
	}
	err := wait.ExponentialBackoff(backoff, func() (bool, error) {
		switch resourceType {
		case loadbalancerResource:
			lb, err := client.LoadBalancer.Get(context.Background(), ID)
			if err != nil {
				return false, err
			}
			if lb.ProvisioningStatus == activeStatus {
				return true, nil
			} else if lb.ProvisioningStatus == errorStatus {
				return true, fmt.Errorf("loadbalancer %w has gone into ERROR state", ID)
			} else {
				return false, nil
			}
		case poolResource:
			pool, err := client.Pool.Get(context.Background(), ID)
			if err != nil {
				return false, err
			}
			if pool.ProvisoningStatus == activeStatus {
				return true, nil
			} else if pool.ProvisoningStatus == errorStatus {
				return true, fmt.Errorf("Pool %w has gone into ERROR state", ID)
			} else {
				return false, nil
			}
		case listenerResource:
			listener, err := client.Listener.Get(context.Background(), ID)
			if err != nil {
				return false, err
			}
			if listener.ProvisoningStatus == activeStatus {
				return true, nil
			} else if listener.ProvisoningStatus == errorStatus {
				return true, fmt.Errorf("Listener %w has gone into ERROR state", ID)
			} else {
				return false, nil
			}
		default:
			return false, nil
		}

	})

	if err == wait.ErrWaitTimeout {
		err = fmt.Errorf("loadbalancer failed to go into ACTIVE provisioning status within allotted time")
		return nil, err
	}
	lb, err := client.LoadBalancer.Get(context.Background(), ID)
	return lb, err
}
