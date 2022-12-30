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
	"strconv"
	"strings"
	"time"

	"github.com/bizflycloud/gobizfly"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBizFlyCloudDatabaseInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceBizFlyCloudCloudDatabaseInstanceCreate,
		Read:   resourceBizFlyCloudCloudDatabaseInstanceRead,
		Update: resourceBizFlyCloudCloudDatabaseInstanceUpdate,
		Delete: resourceBizFlyCloudCloudDatabaseInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(80 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},
		Schema: resourceCloudDatabaseInstanceSchema(),
	}
}

func resourceBizFlyCloudCloudDatabaseInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()

	insc := &gobizfly.CloudDatabaseInstanceCreate{
		Name:         d.Get("name").(string),
		InstanceType: d.Get("instance_type").(string),
		FlavorName:   d.Get("flavor_name").(string),
		VolumeSize:   d.Get("volume_size").(int),
		Datastore: gobizfly.CloudDatabaseDatastore{
			Type:      d.Get("datastore_type").(string),
			VersionID: d.Get("datastore_version_id").(string),
		},
		PublicAccess:     d.Get("public_access").(bool),
		AvailabilityZone: d.Get("availability_zone").(string),
		Networks:         makeArrayNetworkIDFromArray(d.Get("network_ids").(*schema.Set).List()),
		AutoScaling: &gobizfly.CloudDatabaseAutoScaling{
			Enable: d.Get("autoscaling_enable").(bool),
			Volume: gobizfly.CloudDatabaseAutoScalingVolume{
				Threshold: d.Get("autoscaling_volume_threshold").(int),
				Limited:   d.Get("autoscaling_volume_limited").(int),
			},
		},
		BackupID: d.Get("backup_id").(string),
	}

	// retry
	retry := maxRetry
	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		instance, err := client.CloudDatabase.Instances().Create(context.Background(), insc)

		if err != nil {
			retry -= 1
			if retry > 0 {
				time.Sleep(timeSleep)
				return resource.RetryableError(fmt.Errorf("[ERROR] create cloud database instance %s failed: %s. Retrying", d.Get("name"), err))
			}

			return resource.NonRetryableError(fmt.Errorf("[ERROR] create cloud database instance %s failed: %s. Can't retry", d.Get("name"), err))
		}

		log.Printf("[DEBUG] creating cloud database instance %s", instance.Name)

		d.SetId(instance.ID)
		_ = d.Set("task_id", instance.TaskID)

		// wait for cloud database instance to become active
		_, err = waitForCloudDatabaseInstanceCreate(d, meta)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("[ERROR] create cloud database instance (%s) failed: %s. Can't retry", d.Get("name").(string), err))
		}

		err = resourceBizFlyCloudCloudDatabaseInstanceRead(d, meta)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("[ERROR] read cloud database instance (%s) failed: %s. Can't retry", d.Get("name").(string), err))
		}

		return nil
	})
}

func resourceBizFlyCloudCloudDatabaseInstanceRead(d *schema.ResourceData, meta interface{}) error {
	if err := dataSourceBizFlyCloudDatabaseInstanceRead(d, meta); err != nil {
		return err
	}
	return nil
}

func resourceBizFlyCloudCloudDatabaseInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	id := d.Id()
	if d.HasChange("volume_size") {
		// retry
		retry := maxRetry
		err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			task, err := client.CloudDatabase.Instances().ResizeVolume(context.Background(), id, d.Get("volume_size").(int))

			if err != nil {
				retry -= 1
				if retry > 0 {
					time.Sleep(timeSleep)
					return resource.RetryableError(fmt.Errorf("[ERROR] Resize volume of database instance [%s] error: %v. Retrying", id, err))
				}

				return resource.NonRetryableError(fmt.Errorf("[ERROR] Resize volume of database instance [%s] error: %v. Can't retry", id, err))
			}

			_ = d.Set("task_id", task.TaskID)

			err = resourceBizFlyCloudCloudDatabaseInstanceRead(d, meta)
			if err != nil {
				return resource.NonRetryableError(fmt.Errorf("[ERROR] Resize volume of database instance %s failed: %s. Can't retry", id, err))
			}

			// wait for database instance is active again
			_, err = waitForCloudDatabaseInstanceUpdate(d, meta, "volume_size", d.Get("volume_size").(int))
			if err != nil {
				return resource.NonRetryableError(fmt.Errorf("[ERROR] Resize volume of database instance %s with task id (%s) error: %s. Can't retry", id, task.TaskID, err))
			}

			return nil
		})

		if err != nil {
			return fmt.Errorf("[ERROR] Resize volume of database instance %s failed: %s", id, err)
		}
	}

	if d.HasChange("flavor_name") {
		// retry
		retry := maxRetry
		err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			task, err := client.CloudDatabase.Instances().ResizeFlavor(context.Background(), id, d.Get("flavor_name").(string))

			if err != nil {
				retry -= 1
				if retry > 0 {
					time.Sleep(timeSleep)
					return resource.RetryableError(fmt.Errorf("[ERROR] Resize flavor of database instance [%s] failed: %v. Retrying", id, err))
				}

				return resource.NonRetryableError(fmt.Errorf("[ERROR] Resize flavor of database instance [%s] failed: %v. Can't retry", id, err))
			}

			_ = d.Set("task_id", task.TaskID)

			err = resourceBizFlyCloudCloudDatabaseInstanceRead(d, meta)
			if err != nil {
				return resource.NonRetryableError(fmt.Errorf("[ERROR] Resize flavor of database instance %s failed: %s. Can't retry", id, err))
			}

			// wait for database instance is active again
			_, err = waitForCloudDatabaseInstanceUpdate(d, meta, "flavor_name", d.Get("flavor_name").(string))
			if err != nil {
				return resource.NonRetryableError(fmt.Errorf("[ERROR] Resize flavor of database instance %s with task (%s) failed: %s. Can't retry", id, task.TaskID, err))
			}
			return nil
		})

		if err != nil {
			return fmt.Errorf("[ERROR] Resize volume of database instance %s failed: %s", id, err)
		}
	}

	return resourceBizFlyCloudCloudDatabaseInstanceRead(d, meta)
}

func resourceBizFlyCloudCloudDatabaseInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	id := d.Id()

	// retry
	retry := maxRetry
	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		task, err := client.CloudDatabase.Instances().Delete(context.Background(), id, &gobizfly.CloudDatabaseDelete{})

		if err != nil {
			retry -= 1
			if retry > 0 {
				time.Sleep(timeSleep)
				return resource.RetryableError(fmt.Errorf("[ERROR] delete cloud database instance %s failed: %v. Retrying", id, err))
			}

			return resource.NonRetryableError(fmt.Errorf("[ERROR] delete cloud database instance %s failed: %v. Can't retry", id, err))
		}
		_ = d.Set("task_id", task.TaskID)

		// wait for cloud database instance to delete
		_, err = waitForCloudDatabaseInstanceDelete(d, meta)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("[ERROR] delete cloud database instance %s with task %s failed: %v. Can't retry", id, task.TaskID, err))
		}

		return nil
	})
}

func waitForCloudDatabaseInstanceCreate(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	log.Printf("[INFO] Waiting for cloud database instance (%s) to be ready", d.Get("name").(string))
	stateConf := &resource.StateChangeConf{
		Pending:        []string{"false", "BUILD"},
		Target:         []string{"true", "ACTIVE", "HEALTHY"},
		Refresh:        newCloudDatabaseInstanceStateRefreshFunc(d, meta),
		Timeout:        d.Timeout(schema.TimeoutCreate),
		Delay:          60 * time.Second,
		MinTimeout:     20 * time.Second,
		NotFoundChecks: 50,
	}
	return stateConf.WaitForState()
}

func waitForCloudDatabaseInstanceUpdate(d *schema.ResourceData, meta interface{}, key string, newValue interface{}) (interface{}, error) {
	log.Printf("[INFO] Waiting for cloud database instance (%s) to be update", d.Get("name").(string))
	stateConf := &resource.StateChangeConf{
		Pending:        []string{"false", "RESIZE", "RESTART_REQUIRED", "REBOOTING"},
		Target:         []string{"true", "ACTIVE", "HEALTHY"},
		Refresh:        updateCloudDatabaseInstanceStateRefreshFunc(d, key, newValue, meta),
		Timeout:        d.Timeout(schema.TimeoutUpdate),
		Delay:          60 * time.Second,
		MinTimeout:     10 * time.Second,
		NotFoundChecks: 30,
	}
	return stateConf.WaitForState()
}

func waitForCloudDatabaseInstanceDelete(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	log.Printf("[INFO] Waiting for cloud database instance (%s) to be delete", d.Get("name").(string))
	stateConf := &resource.StateChangeConf{
		Pending:        []string{"false", "SHUTDOWN"},
		Target:         []string{"true"},
		Refresh:        deleteCloudDatabaseInstanceStateRefreshFunc(d, meta),
		Timeout:        d.Timeout(schema.TimeoutDelete),
		Delay:          30 * time.Second,
		MinTimeout:     10 * time.Second,
		NotFoundChecks: 30,
	}
	return stateConf.WaitForState()
}

func newCloudDatabaseInstanceStateRefreshFunc(d *schema.ResourceData, meta interface{}) resource.StateRefreshFunc {
	client := meta.(*CombinedConfig).gobizflyClient()

	return func() (interface{}, string, error) {
		err := resourceBizFlyCloudCloudDatabaseInstanceRead(d, meta)
		if err != nil {
			return nil, "", err
		}

		if attr, ok := d.GetOk("status"); ok {
			ins, err := client.CloudDatabase.Instances().Get(context.Background(), d.Id())
			if err != nil {
				return nil, "", fmt.Errorf("[ERROR] Retrieving cloud database instance %s error: %v", d.Id(), err)
			}
			switch attr := attr.(type) {
			case bool:
				return &ins, strconv.FormatBool(attr), nil
			default:
				return &ins, attr.(string), nil
			}
		}
		return nil, "", nil
	}
}

func updateCloudDatabaseInstanceStateRefreshFunc(d *schema.ResourceData, key string, newValue interface{}, meta interface{}) resource.StateRefreshFunc {
	client := meta.(*CombinedConfig).gobizflyClient()

	return func() (interface{}, string, error) {
		err := resourceBizFlyCloudCloudDatabaseInstanceRead(d, meta)
		if err != nil {
			return nil, "", err
		}

		ins, err := client.CloudDatabase.Instances().Get(context.Background(), d.Id())
		if err != nil {
			return nil, "", fmt.Errorf("[ERROR] Retrieving cloud database instance %s error: %v", d.Id(), err)
		}

		switch key {
		case "volume_size":
			if ins.Volume.Size < newValue.(int) {
				log.Println("[DEBUG] Cloud database instance is updating")
				return nil, "", nil
			}

		case "flavor_name":
			for _, node := range ins.Nodes {
				if strings.Contains(node.Flavor, newValue.(string)) {
					log.Println("[DEBUG] Cloud database instance is updating")
					return nil, "", nil
				}
			}
		}

		if attr, ok := d.GetOk("status"); ok {
			ins, err := client.CloudDatabase.Instances().Get(context.Background(), d.Id())
			if err != nil {
				return nil, "", fmt.Errorf("[ERROR] Retrieving cloud database instance %s error: %v", d.Id(), err)
			}
			switch attr := attr.(type) {
			case bool:
				return &ins, strconv.FormatBool(attr), nil
			default:
				return &ins, attr.(string), nil
			}
		}
		return nil, "", nil
	}
}

func deleteCloudDatabaseInstanceStateRefreshFunc(d *schema.ResourceData, meta interface{}) resource.StateRefreshFunc {
	client := meta.(*CombinedConfig).gobizflyClient()

	return func() (interface{}, string, error) {
		ins, err := client.CloudDatabase.Instances().Get(context.Background(), d.Id())

		if errors.Is(err, gobizfly.ErrNotFound) {
			return ins, "true", nil
		} else if err != nil {
			return nil, "false", err
		}
		return ins, "false", nil
	}
}

func makeArrayNetworkIDFromArray(items []interface{}) []gobizfly.CloudDatabaseNetworks {
	stringDictArray := make([]gobizfly.CloudDatabaseNetworks, 0)
	for i := 0; i < len(items); i++ {
		networkID := items[i].(string)
		networkDict := gobizfly.CloudDatabaseNetworks{NetworkID: networkID}
		stringDictArray = append(stringDictArray, networkDict)
	}
	return stringDictArray
}
