package bizflycloud

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bizflycloud/terraform-provider-bizflycloud/constants"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceBizflyCloudKubernetesWorkerPool() *schema.Resource {
	return &schema.Resource{
		Create: resourceBizflyCloudKubernetesWorkerPoolCreate,
		Read:   resourceBizflyCloudKubernetesWorkerPoolRead,
		Delete: resourceBizflyCloudKubernetesWorkerPoolDelete,
		Update: resourceBizflyCloudKubernetesWorkerPoolUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"flavor": {
				Type:     schema.TypeString,
				Required: true,
			},
			"profile_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"volume_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"volume_size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
			},
			"desired_size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"enable_autoscaling": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"min_size": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"max_size": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"taints": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: taintsSchema(),
				},
			},
			"network_plan": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      constants.FreeDatatransfer,
				ValidateFunc: validation.StringInSlice(constants.ValidNetworkPlans, false),
			},
			"billing_plan": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      constants.OnDemand,
				ValidateFunc: validation.StringInSlice(constants.ValidBillingPlans, false),
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
		},
	}
}

func resourceBizflyCloudKubernetesWorkerPoolCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceBizflyCloudKubernetesWorkerPoolRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	workerPoolID := d.Id()
	clusterID := d.Get("cluster_id").(string)
	log.Printf("[DEBUG] worker pool ID %s", workerPoolID)
	workerPool, err := client.KubernetesEngine.GetClusterWorkerPool(context.Background(), clusterID, workerPoolID)
	if err != nil {
		return fmt.Errorf("[ERROR] Get worker pool %v of cluster %v error: %v", workerPoolID, clusterID, err)
	}
	_ = d.Set("id", workerPool.UID)
	return nil
}

func resourceBizflyCloudKubernetesWorkerPoolUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceBizflyCloudKubernetesWorkerPoolDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
