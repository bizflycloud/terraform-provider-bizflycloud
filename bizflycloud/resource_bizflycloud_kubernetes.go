// This file is part of gobizfly
//
// Copyright (C) 2020  BizFly Cloud
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
	"github.com/bizflycloud/gobizfly"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

func resourceBizFlyKubernetes() *schema.Resource {
	return &schema.Resource{
		Create:        resourceBizFlyClusterCreate,
		Read:          resourceBizFlyCloudClusterRead,
		Delete:        resourceBizFlyCloudClusterDelete,
		Update:        resourceBizFlyCloudClusterUpdate,
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
				Required: true,
			},
			"min_size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"max_size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"tags": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
			"auto_upgrade": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"enable_cloud": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"worker_pool_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"worker_pools": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Resource{Schema: workerPoolSchema()},
			},
		},
	}
}

func workerPoolSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"version": {
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
			Optional: true,
		},
		"volume_size": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"available_zone": {
			Type:     schema.TypeString,
			Required: true,
		},
		"desire_size": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"enable_autoscaling": {
			Type:     schema.TypeBool,
			Required: true,
		},
		"min_size": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"max_size": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"tags": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
	}
}
func resourceBizFlyClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()

	// Build up creation options
	ccrq := &gobizfly.ClusterCreateRequest{
		Name:        d.Get("name").(string),
		Version:     d.Get("version").(string),
		AutoUpgrade: d.Get("auto_upgrade").(bool),
		EnableCloud: d.Get("enable_cloud").(bool),
		Tags:        d.Get("tags").([]string),
		WorkerPools: d.Get("worker_pools").([]gobizfly.WorkerPool),
	}
	log.Printf("[DEBUG] Create Cluster configuration: %#v\n", ccrq)
	cluster, err := client.KubernetesEngine.Create(context.Background(), ccrq)
	if err != nil {
		fmt.Errorf("Error creating cluster: %v", err)
	}
	d.SetId(cluster.UID)
	err = resourceBizFlyCloudClusterRead(d, meta)
	if err != nil {
		fmt.Errorf("Error retrieving cluster: %v", err)
	}
	return nil
}

func resourceBizFlyCloudClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	cluster, err := client.KubernetesEngine.Get(context.Background(), d.Id())
	if err != nil {
		if errors.Is(err, gobizfly.ErrNotFound) {
			log.Printf("[WARN] Cluster %s is not found", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieved cluster: %v", err)
	}
	d.Set("cluster_id", cluster.UID)
	d.Set("name", cluster.Name)
	d.Set("version", cluster.Version)
	d.Set("status", cluster.ClusterStatus)
	d.Set("auto_upgrade", cluster.AutoUpgrade)
	d.Set("worker_pools_count", cluster.WorkerPoolsCount)
	d.Set("create_at", cluster.CreatedAt)
	d.Set("created_by", cluster.CreatedBy)
	d.Set("worker_pools", cluster.WorkerPools)
	return nil
}

func resourceBizFlyCloudClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	err := client.KubernetesEngine.Delete(context.Background(), d.Id())
	if err != nil {
		fmt.Errorf("Error delete cluster: %v", err)
	}
	return nil
}

func resourceBizFlyCloudClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	cluster, err := client.KubernetesEngine.Get(context.Background(), d.Get("cluster_id").(string))
	if err != nil {
		fmt.Errorf("Error update cluster: %v", err)
	}
	if d.HasChange("worker_pools") {
		newPools := d.Get("worker_pools").([]gobizfly.WorkerPool)
		addPools := make([]gobizfly.WorkerPool, 0)
		isOldPool := make(map[string]bool)
		isNewPool := make(map[string]bool)
		for _, pool := range cluster.WorkerPools {
			isOldPool[pool.Name] = true
		}
		for _, pool := range newPools {
			isNewPool[pool.Name] = true
			if _, ok := isOldPool[pool.Name]; ok {
				addPools = append(addPools, pool)
			}
		}
		awrq := &gobizfly.AddWorkerPoolsRequest{
			WorkerPools: addPools,
		}
		_, err := client.KubernetesEngine.AddWorkerPools(context.Background(), d.Id(), awrq)
		if err != nil {
			fmt.Errorf("Error add pool: %v", err)
		}

		for _, pool := range cluster.WorkerPools {
			if isOldPool[pool.Name] == true && isNewPool[pool.Name] == false {
				err := client.KubernetesEngine.DeleteClusterWorkerPool(context.Background(), d.Id(), pool.UID)
				if err != nil {
					fmt.Errorf("Error delete pool: %v", err)
				}
			}
		}
	}
	return resourceBizFlyCloudClusterRead(d, meta)
}
