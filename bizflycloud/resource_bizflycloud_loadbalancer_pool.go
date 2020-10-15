package bizflycloud

import (
	"context"
	"fmt"

	"github.com/bizflycloud/gobizfly"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBizFlyCloudLoadBalancerPool() *schema.Resource {
	return &schema.Resource{
		Create: resourceBizFlyCloudLoadBalancerPoolCreate,
		Update: resourceBizFlyCloudLoadBalancerPoolUpdate,
		Read:   resourceBizFlyCloudLoadBalancerPoolRead,
		Delete: resourceBizFlyCloudLoadBalancerPoolDelete,
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
		},
	}
}
func resourceBizFlyCloudLoadBalancerPoolCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	pcr := gobizfly.PoolCreateRequest{
		Name:        d.Get("name").(*string),
		LBAlgorithm: d.Get("algorithm").(string),
		Protocol:    d.Get("protocol").(string),
		Description: d.Get("description").(*string),
	}
	pool, err := client.Pool.Create(context.Background(), d.Get("load_balancer_id").(string), &pcr)
	if err != nil {
		return fmt.Errorf("Error when creating load balancer pool: %v", err)
	}
	d.SetId(pool.ID)
	return resourceBizFlyCloudLoadBalancerPoolRead(d, meta)
}

func resourceBizFlyCloudLoadBalancerPoolUpdate(d *schema.ResourceData, meta interface{}) error {
	//client := meta.(*CombinedConfig).gobizflyClient()
	return nil
}
func resourceBizFlyCloudLoadBalancerPoolRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	pool, err := client.Pool.Get(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error when retrieving load balancer pool: %v", err)
	}
	d.Set("name", pool.Name)
	d.Set("algorithm", pool.LBAlgorithm)
	d.Set("description", pool.Description)
	d.Set("protocol", pool.Protocol)
	d.Set("load_balancer_id", pool.LoadBalancers[0].ID)
	return nil
}

func resourceBizFlyCloudLoadBalancerPoolDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	err := client.Pool.Delete(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error when deleting load balancer pool: %v", err)
	}
	return nil
}
