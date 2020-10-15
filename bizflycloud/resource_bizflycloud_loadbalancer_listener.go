package bizflycloud

import (
	"context"
	"fmt"

	"github.com/bizflycloud/gobizfly"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBizFlyCloudLoadBalancerListener() *schema.Resource {
	return &schema.Resource{
		Create: resourceBizFlyCloudLoadBalancerListenerCreate,
		Read:   resourceBizFlyCloudLoadBalancerListenerRead,
		Update: resourceBizFlyCloudLoadBalancerListenerUpdate,
		Delete: resourceBizFlyCloudLoadBalancerListenerDelete,
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
				Type:     schema.TypeString,
				Required: true,
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
		},
	}
}

func resourceBizFlyCloudLoadBalancerListenerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	lcr := gobizfly.ListenerCreateRequest{
		Name:                   d.Get("name").(*string),
		Protocol:               d.Get("protocol").(string),
		ProtocolPort:           d.Get("port").(int),
		Description:            d.Get("description").(*string),
		DefaultPoolID:          d.Get("default_pool_id").(*string),
		DefaultTLSContainerRef: d.Get("default_tls_ref").(*string),
	}
	listener, err := client.Listener.Create(context.Background(), d.Get("load_balancer_id").(string), &lcr)
	if err != nil {
		return fmt.Errorf("Error when creating listener: %v", err)
	}
	d.SetId(listener.ID)
	return resourceBizFlyCloudLoadBalancerListenerRead(d, meta)
}

func resourceBizFlyCloudLoadBalancerListenerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	listener, err := client.Listener.Get(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error when retrieving listener: %v", err)
	}
	d.Set("name", listener.Name)
	d.Set("protocol", listener.Protocol)
	d.Set("port", listener.ProtocolPort)
	d.Set("description", listener.Description)
	d.Set("default_pool_id", listener.DefaultPoolID)
	d.Set("default_tls_ref", listener.DefaultTLSContainerRef)
	d.Set("load_balancer_id", listener.LoadBalancers[0].ID)
	return nil
}

func resourceBizFlyCloudLoadBalancerListenerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	lur := gobizfly.ListenerUpdateRequest{
		Name:                   d.Get("name").(*string),
		Description:            d.Get("description").(*string),
		DefaultTLSContainerRef: d.Get("default_tls_ref").(*string),
		DefaultPoolID:          d.Get("default_pool_id").(*string),
	}
	_, err := client.Listener.Update(context.Background(), d.Id(), &lur)
	if err != nil {
		return fmt.Errorf("Error when updating listener: %v")
	}
	return resourceBizFlyCloudLoadBalancerListenerRead(d, meta)
}

func resourceBizFlyCloudLoadBalancerListenerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	err := client.Listener.Delete(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error when deleting listener: %v", err)
	}
	return nil
}
