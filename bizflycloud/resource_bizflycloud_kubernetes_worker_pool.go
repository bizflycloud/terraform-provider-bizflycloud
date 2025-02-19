package bizflycloud

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/bizflycloud/gobizfly"
	"github.com/bizflycloud/terraform-provider-bizflycloud/constants"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
			"provision_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
		},
	}
}

func resourceBizflyCloudKubernetesWorkerPoolCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	clusterID := d.Get("cluster_id").(string)
	workerPool := parseWorkerPoolFromConfig(d, meta)
	payload := gobizfly.AddWorkerPoolsRequest{
		WorkerPools: []gobizfly.WorkerPool{
			workerPool.WorkerPool,
		},
	}
	log.Printf("[DEBUG] add worker pool payload: %+v", payload)
	addedWorkerPools, err := client.KubernetesEngine.AddWorkerPools(context.Background(), clusterID, &payload)
	if err != nil {
		return fmt.Errorf("[ERROR] add worker pools for cluster %v error: %v", clusterID, err)
	}
	poolID := addedWorkerPools[0].UID
	_, err = waitForWorkerPoolChange(d, poolID, meta)
	if err != nil {
		return fmt.Errorf("[ERROR] Wait for worker pool create error: %v", err)
	}
	d.SetId(poolID)
	return resourceBizflyCloudKubernetesWorkerPoolRead(d, meta)
}

func resourceBizflyCloudKubernetesWorkerPoolRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	workerPoolID := d.Id()
	log.Printf("[DEBUG] worker pool ID %s", workerPoolID)
	workerPool, err := client.KubernetesEngine.GetDetailWorkerPool(context.Background(), workerPoolID)
	if err != nil {
		if errors.Is(err, gobizfly.ErrNotFound) {
			log.Printf("[WARN] Worker pool %s is not found", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Get worker pool %v error: %v", workerPoolID, err)
	}
	_ = d.Set("id", workerPool.UID)
	_ = d.Set("cluster_id", workerPool.ShootID)
	_ = d.Set("name", workerPool.Name)
	_ = d.Set("flavor", workerPool.Flavor)
	_ = d.Set("profile_type", workerPool.ProfileType)
	_ = d.Set("volume_type", workerPool.VolumeType)
	_ = d.Set("volume_size", workerPool.VolumeSize)
	_ = d.Set("availability_zone", workerPool.AvailabilityZone)
	_ = d.Set("desired_size", workerPool.DesiredSize)
	_ = d.Set("enable_autoscaling", workerPool.EnableAutoScaling)
	_ = d.Set("min_size", workerPool.MinSize)
	_ = d.Set("max_size", workerPool.MaxSize)
	_ = d.Set("tags", workerPool.Tags)
	_ = d.Set("labels", workerPool.Labels)
	_ = d.Set("taints", parseWorkerPoolTaints(workerPool.Taints))
	_ = d.Set("network_plan", workerPool.NetworkPlan)
	_ = d.Set("billing_plan", workerPool.BillingPlan)
	_ = d.Set("provision_status", workerPool.ProvisionStatus)
	return nil
}

func resourceBizflyCloudKubernetesWorkerPoolUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	changedKeys := []string{"labels", "taints", "min_size", "max_size", "desired_size"}
	if d.HasChanges(changedKeys...) {
		clusterID := d.Get("cluster_id").(string)
		poolID := d.Id()
		newWorkerPool := parseWorkerPoolFromConfig(d, meta)
		oldWorkerPool, err := client.KubernetesEngine.GetClusterWorkerPool(context.Background(), clusterID, poolID)
		if err != nil {
			if errors.Is(err, gobizfly.ErrNotFound) {
				log.Printf("[WARN] Worker pool %s of cluster %v is not found", d.Id(), clusterID)
				d.SetId("")
				return nil
			}
			return fmt.Errorf("[ERROR] Get worker pool %v of cluster %v error: %v", poolID, clusterID, err)
		}
		// Check that the pool has any change
		isUpdateLabels := !cmp.Equal(newWorkerPool.Labels, oldWorkerPool.Labels)
		isUpdateTaints := !cmp.Equal(newWorkerPool.Taints, oldWorkerPool.Taints)
		isUpdateMaxSize := newWorkerPool.MaxSize != oldWorkerPool.MaxSize
		isUpdateMinSize := newWorkerPool.MinSize != oldWorkerPool.MinSize
		isUpdateDesiredSize := newWorkerPool.DesiredSize != oldWorkerPool.DesiredSize
		isUpdatePool := isUpdateLabels || isUpdateTaints || isUpdateMaxSize || isUpdateMinSize || isUpdateDesiredSize
		if isUpdatePool {
			fmt.Printf("[DEBUG] Old pool state: %+v\nNew pool state: %+v", oldWorkerPool, newWorkerPool)
			updateRequest := &gobizfly.UpdateWorkerPoolRequest{
				DesiredSize: newWorkerPool.DesiredSize,
				MinSize:     newWorkerPool.MinSize,
				MaxSize:     newWorkerPool.MaxSize,
				Labels:      newWorkerPool.Labels,
				Taints:      newWorkerPool.Taints,
			}
			log.Printf("[DEBUG] update pool %+v to %+v", poolID, updateRequest)
			err := client.KubernetesEngine.UpdateClusterWorkerPool(context.Background(), clusterID,
				poolID, updateRequest)
			if err != nil {
				return fmt.Errorf("error update pool: %+v", err)
			}
			_, err = waitForWorkerPoolChange(d, poolID, meta)
			if err != nil {
				return fmt.Errorf("[ERROR] Wait for worker pool update error: %v", err)
			}
		}
	}
	return resourceBizflyCloudKubernetesWorkerPoolRead(d, meta)
}

func resourceBizflyCloudKubernetesWorkerPoolDelete(d *schema.ResourceData, meta interface{}) error {
	workerPoolID := d.Id()
	clusterID := d.Get("cluster_id").(string)
	client := meta.(*CombinedConfig).gobizflyClient()
	if err := client.KubernetesEngine.DeleteClusterWorkerPool(context.Background(), clusterID, workerPoolID); err != nil {
		if errors.Is(err, gobizfly.ErrNotFound) {
			log.Printf("[WARN] Worker pool %s is not found", workerPoolID)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Delete worker pool %v of cluster %v error: %v", workerPoolID, clusterID, err)
	}
	return nil
}

func parseWorkerPoolFromConfig(d *schema.ResourceData, meta interface{}) gobizfly.ExtendedWorkerPool {
	tags := make([]string, 0)
	for _, tag := range d.Get("tags").([]interface{}) {
		tags = append(tags, tag.(string))
	}
	labels := readLabelsConfig(d, "")
	taints := readTaintsConfig(d, "")
	pool := gobizfly.ExtendedWorkerPool{
		UID: d.Id(),
		WorkerPool: gobizfly.WorkerPool{
			Name:              d.Get("name").(string),
			Flavor:            d.Get("flavor").(string),
			ProfileType:       d.Get("profile_type").(string),
			VolumeType:        d.Get("volume_type").(string),
			VolumeSize:        d.Get("volume_size").(int),
			AvailabilityZone:  d.Get("availability_zone").(string),
			DesiredSize:       d.Get("desired_size").(int),
			EnableAutoScaling: d.Get("enable_autoscaling").(bool),
			MinSize:           d.Get("min_size").(int),
			MaxSize:           d.Get("max_size").(int),
			NetworkPlan:       d.Get("network_plan").(string),
			BillingPlan:       d.Get("billing_plan").(string),
			Tags:              tags,
			Labels:            labels,
			Taints:            taints,
		},
	}
	return pool
}

func waitForWorkerPoolChange(d *schema.ResourceData, poolID string, meta interface{}) (interface{}, error) {
	log.Printf("[INFO] Waiting for pool updating %s", poolID)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING_PROVISION", "PROVISIONING", "PENDING_UPDATE", "UPDATING"},
		Target:     []string{"PROVISIONED"},
		Refresh:    workerPoolStatusRefreshFunc(d, poolID, meta),
		Timeout:    1200 * time.Second,
		Delay:      20 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	return stateConf.WaitForState()
}

func workerPoolStatusRefreshFunc(d *schema.ResourceData, poolID string, meta interface{}) resource.StateRefreshFunc {
	client := meta.(*CombinedConfig).gobizflyClient()
	return func() (interface{}, string, error) {
		pool, err := client.KubernetesEngine.GetDetailWorkerPool(context.Background(), poolID)
		if err != nil {
			return nil, "", err
		}
		return pool, pool.ProvisionStatus, nil
	}
}
