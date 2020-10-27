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
			"members": {
				Type:       schema.TypeSet,
				Optional:   true,
				Computed:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: loadbalancerMemberSchema(),
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
	poolDescription := d.Get("description").(string)
	pcr := gobizfly.PoolCreateRequest{
		Name:        &poolName,
		LBAlgorithm: d.Get("algorithm").(string),
		Protocol:    d.Get("protocol").(string),
		Description: &poolDescription,
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
			Required: true,
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
