// This file is part of terraform-provider-bizflycloud
//
// Copyright (C) 2021  Bizfly Cloud
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>

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

func resourceBizflyCloudKubernetes() *schema.Resource {
	return &schema.Resource{
		Create: resourceBizflyClusterCreate,
		Read:   resourceBizflyCloudClusterRead,
		Delete: resourceBizflyCloudClusterDelete,
		Update: resourceBizflyCloudClusterUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"package_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"create_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"worker_pools": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Resource{Schema: workerPoolSchema()},
			},
			"worker_pools_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vpc_network_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"auto_upgrade": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"local_dns": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"cni_plugin": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      constants.KubernetesKubeRouter,
				ValidateFunc: validation.StringInSlice(constants.ValidCNIPlugins, false),
			},
			"enabled_upgrade_version": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"is_latest": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"current_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"next_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"package_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
		},
	}
}

func workerPoolSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Computed: true,
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
	}
}

func taintsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"effect": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(constants.ValidEffects, false),
		},
		"key": {
			Type:     schema.TypeString,
			Required: true,
		},
		"value": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
}

func resourceBizflyClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()

	// Build up creation options
	tags := make([]string, 0)
	for j := 0; j < len(d.Get("tags").([]interface{})); j++ {
		tagPattern := fmt.Sprintf("tags.%d", j)
		tags = append(tags, d.Get(tagPattern).(string))
	}

	log.Println("[DEBUG] creating cluster")
	ccrq := &gobizfly.ClusterCreateRequest{
		Name:         d.Get("name").(string),
		Version:      d.Get("version").(string),
		Package:      d.Get("package_id").(string),
		VPCNetworkID: d.Get("vpc_network_id").(string),
		AutoUpgrade:  d.Get("auto_upgrade").(bool),
		LocalDNS:     d.Get("local_dns").(bool),
		CNIPlugin:    d.Get("cni_plugin").(string),
		WorkerPools:  readWorkerPoolFromConfig(d),
		Tags:         tags,
	}
	log.Printf("[DEBUG] Create Cluster configuration: %#v\n", ccrq)
	cluster, err := client.KubernetesEngine.Create(context.Background(), ccrq)
	if err != nil {
		return fmt.Errorf("Error creating cluster: %v", err)
	}
	log.Println("[DEBUG] set id " + cluster.UID)
	d.SetId(cluster.UID)
	return resourceBizflyCloudClusterRead(d, meta)
}

func resourceBizflyCloudClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	log.Printf("[DEBUG] cluster ID %s", d.Id())
	cluster, err := client.KubernetesEngine.Get(context.Background(), d.Id())
	if err != nil {
		if errors.Is(err, gobizfly.ErrNotFound) {
			log.Printf("[WARN] Cluster %s is not found", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieved cluster: %v", err)
	}
	clusterID := cluster.UID
	upgradeVersion, err := client.KubernetesEngine.GetUpgradeClusterVersion(context.Background(), clusterID)
	if err != nil {
		log.Printf("Error get upragde cluster %s version failed: %+v", clusterID, err)
		upgradeVersion = &gobizfly.UpgradeClusterVersionResponse{}
	}
	_ = d.Set("name", cluster.Name)
	_ = d.Set("version", cluster.Version.ID)
	_ = d.Set("package_name", cluster.ClusterPackage.Name)
	_ = d.Set("vpc_network_id", cluster.VPCNetworkID)
	_ = d.Set("worker_pools_count", cluster.WorkerPoolsCount)
	_ = d.Set("create_at", cluster.CreatedAt)
	_ = d.Set("created_by", cluster.CreatedBy)
	_ = d.Set("auto_upgrade", cluster.AutoUpgrade)
	_ = d.Set("local_dns", cluster.LocalDNS)
	_ = d.Set("cni_plugin", cluster.CNIPlugin)
	_ = d.Set("is_latest", upgradeVersion.IsLatest)
	_ = d.Set("current_version", cluster.Version.K8SVersion)
	_ = d.Set("next_version", upgradeVersion.UpgradeTo)
	_ = d.Set("enabled_upgrade_version", false)

	workerPools := parseWorkerPools(cluster.WorkerPools)
	_ = d.Set("worker_pools", workerPools)
	return nil
}

func resourceBizflyCloudClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	err := client.KubernetesEngine.Delete(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error delete cluster: %v", err)
	}
	return nil
}

func resourceBizflyCloudClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	clusterID := d.Id()
	cluster, err := client.KubernetesEngine.Get(context.Background(), clusterID)
	if err != nil {
		return fmt.Errorf("Error update cluster: %v", err)
	}
	if d.HasChange("auto_upgrade") {
		_, new_auto_upgrade := d.GetChange("auto_upgrade")
		update_auto_upgrade := new_auto_upgrade.(bool)
		updateClusterPayload := gobizfly.UpdateClusterRequest{
			AutoUpgrade: &update_auto_upgrade,
		}
		log.Printf("[DEBUG] Update cluster payload: %+v", updateClusterPayload)
		_, err = client.KubernetesEngine.UpdateCluster(context.Background(), clusterID, &updateClusterPayload)
		if err != nil {
			return fmt.Errorf("Error update auto_upgrade: %+v", err)
		}
	}
	if d.HasChange("enabled_upgrade_version") {
		if d.Get("is_latest").(bool) {
			log.Printf("[DEBUG] Cluster version is latest.")
		} else {
			payload := gobizfly.UpgradeClusterVersionRequest{}
			err := client.KubernetesEngine.UpgradeClusterVersion(context.Background(), clusterID, &payload)
			if err != nil {
				return fmt.Errorf("Upragde cluster version error: %+v", err)
			}
		}
	}
	if d.HasChange("worker_pools") {
		newPools := readWorkerPoolFromConfig(d)
		addPools := make([]gobizfly.WorkerPool, 0)
		isOldPool := make(map[string]bool)
		isNewPool := make(map[string]bool)
		for _, pool := range cluster.WorkerPools {
			isOldPool[pool.Name] = true
		}
		newPoolMap := make(map[string]gobizfly.WorkerPool, len(newPools))
		for _, pool := range newPools {
			newPoolMap[pool.Name] = pool
			if !isOldPool[pool.Name] {
				addPools = append(addPools, pool)
			} else {
				isNewPool[pool.Name] = true
			}
		}

		for _, oldPool := range cluster.WorkerPools {
			if isNewPool[oldPool.Name] {
				// Check that the pool has any change
				newStatePool := newPoolMap[oldPool.Name]
				isUpdateLabels := !cmp.Equal(newStatePool.Labels, oldPool.Labels)
				isUpdateTaints := !cmp.Equal(newStatePool.Taints, oldPool.Taints)
				isUpdatePool := (oldPool.MaxSize != newStatePool.MaxSize) || (oldPool.MinSize != newStatePool.MinSize) ||
					(oldPool.DesiredSize != newStatePool.DesiredSize) || isUpdateLabels || isUpdateTaints

				if isUpdatePool {
					fmt.Printf("[DEBUG] Old pool state: %+v\nNew pool state: %+v", oldPool, newStatePool)
					updateRequest := &gobizfly.UpdateWorkerPoolRequest{
						DesiredSize: newStatePool.DesiredSize,
						MinSize:     newStatePool.MinSize,
						MaxSize:     newStatePool.MaxSize,
						Labels:      newStatePool.Labels,
						Taints:      newStatePool.Taints,
					}
					log.Printf("[DEBUG] update pool %+v to %+v", oldPool, updateRequest)
					err := client.KubernetesEngine.UpdateClusterWorkerPool(context.Background(), d.Id(),
						oldPool.UID, updateRequest)
					if err != nil {
						return fmt.Errorf("error update pool: %+v", err)
					}
					_, err = waitForPoolUpdate(d, oldPool.UID, meta)
					if err != nil {
						return fmt.Errorf("error waiting for pool update: %+v", err)
					}
				}
			}
		}
		log.Printf("[DEBUG] add Pools %+v", addPools)
		if len(addPools) > 0 {
			awrq := &gobizfly.AddWorkerPoolsRequest{
				WorkerPools: addPools,
			}
			_, err := client.KubernetesEngine.AddWorkerPools(context.Background(), d.Id(), awrq)
			if err != nil {
				return fmt.Errorf("Error add pool: %v", err)
			}
		}
		for _, pool := range cluster.WorkerPools {
			if isOldPool[pool.Name] && !isNewPool[pool.Name] {
				log.Printf("[DEBUG] remove pool %+v", pool)
				err := client.KubernetesEngine.DeleteClusterWorkerPool(context.Background(), d.Id(), pool.UID)
				if err != nil {
					return fmt.Errorf("Error delete pool: %v", err)
				}
			}
		}
	}
	return resourceBizflyCloudClusterRead(d, meta)
}

func readWorkerPoolFromConfig(l *schema.ResourceData) []gobizfly.WorkerPool {
	pools := make([]gobizfly.WorkerPool, 0)
	for i := 0; i < len(l.Get("worker_pools").([]interface{})); i++ {
		pattern := fmt.Sprintf("worker_pools.%d.", i)
		tags := make([]string, 0)
		for j := 0; j < len(l.Get("tags").([]interface{})); j++ {
			tagPattern := pattern + fmt.Sprintf("tags.%d", j)
			tags = append(tags, l.Get(tagPattern).(string))
		}
		labels := readLabelsConfig(l, pattern)
		taints := readTaintsConfig(l, pattern)
		pool := gobizfly.WorkerPool{
			Name:              l.Get(pattern + "name").(string),
			Flavor:            l.Get(pattern + "flavor").(string),
			ProfileType:       l.Get(pattern + "profile_type").(string),
			VolumeType:        l.Get(pattern + "volume_type").(string),
			VolumeSize:        l.Get(pattern + "volume_size").(int),
			AvailabilityZone:  l.Get(pattern + "availability_zone").(string),
			DesiredSize:       l.Get(pattern + "desired_size").(int),
			EnableAutoScaling: l.Get(pattern + "enable_autoscaling").(bool),
			MinSize:           l.Get(pattern + "min_size").(int),
			MaxSize:           l.Get(pattern + "max_size").(int),
			NetworkPlan:       l.Get(pattern + "network_plan").(string),
			BillingPlan:       l.Get(pattern + "billing_plan").(string),
			Tags:              tags,
			Labels:            labels,
			Taints:            taints,
		}
		pools = append(pools, pool)
	}
	return pools
}

func readTaintsConfig(l *schema.ResourceData, pattern string) []gobizfly.Taint {
	taints := make([]gobizfly.Taint, 0)
	taintsData := l.Get(pattern + "taints")
	if taintsData == nil {
		return nil
	}

	for i := 0; i < len(taintsData.([]interface{})); i++ {
		taintPattern := fmt.Sprintf("%staints.%d.", pattern, i)
		taint := gobizfly.Taint{
			Effect: l.Get(taintPattern + "effect").(string),
			Key:    l.Get(taintPattern + "key").(string),
			Value:  l.Get(taintPattern + "value").(string),
		}
		taints = append(taints, taint)
	}
	return taints
}

func readLabelsConfig(l *schema.ResourceData, pattern string) map[string]string {
	labels := make(map[string]string)
	labelsData := l.Get(pattern + "labels")
	if labelsData == nil {
		return nil
	}

	for key, val := range labelsData.(map[string]interface{}) {
		labels[key] = fmt.Sprintf("%v", val)
	}
	return labels
}

func waitForPoolUpdate(d *schema.ResourceData, poolID string, meta interface{}) (interface{}, error) {
	log.Printf("[INFO] Waiting for pool updating %s", poolID)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING_PROVISION", "PROVISIONING", "PENDING_UPDATE", "UPDATING"},
		Target:     []string{"PROVISIONED"},
		Refresh:    newPoolStatusRefreshFunc(d, poolID, meta),
		Timeout:    1200 * time.Second,
		Delay:      20 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	return stateConf.WaitForState()
}

func newPoolStatusRefreshFunc(d *schema.ResourceData, poolID string, meta interface{}) resource.StateRefreshFunc {
	client := meta.(*CombinedConfig).gobizflyClient()
	return func() (interface{}, string, error) {
		pool, err := client.KubernetesEngine.GetClusterWorkerPool(context.Background(), d.Id(), poolID)
		if err != nil {
			return nil, "", err
		}
		return pool, pool.ProvisionStatus, nil
	}
}

func parseWorkerPoolTaints(taints []gobizfly.Taint) []map[string]interface{} {
	results := make([]map[string]interface{}, 0)
	for _, taint := range taints {
		result := map[string]interface{}{
			"effect": taint.Effect,
			"key":    taint.Key,
			"value":  taint.Value,
		}
		results = append(results, result)
	}
	return results
}

func parseWorkerPools(workerPools []gobizfly.ExtendedWorkerPool) []map[string]interface{} {
	results := make([]map[string]interface{}, 0)
	for _, workerPool := range workerPools {
		taints := parseWorkerPoolTaints(workerPool.Taints)
		result := map[string]interface{}{
			"id":                 workerPool.UID,
			"name":               workerPool.Name,
			"flavor":             workerPool.Flavor,
			"profile_type":       workerPool.ProfileType,
			"volume_type":        workerPool.VolumeType,
			"volume_size":        workerPool.VolumeSize,
			"availability_zone":  workerPool.AvailabilityZone,
			"desired_size":       workerPool.DesiredSize,
			"enable_autoscaling": workerPool.EnableAutoScaling,
			"min_size":           workerPool.MinSize,
			"max_size":           workerPool.MaxSize,
			"tags":               workerPool.Tags,
			"labels":             workerPool.Labels,
			"taints":             taints,
			"network_plan":       workerPool.NetworkPlan,
			"billing_plan":       workerPool.BillingPlan,
		}
		results = append(results, result)
	}
	return results
}
