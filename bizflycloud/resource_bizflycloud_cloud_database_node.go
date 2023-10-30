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

func resourceBizflyCloudCloudDatabaseNodeRead(d *schema.ResourceData, meta interface{}) error {
	if err := dataSourceBizflyCloudDatabaseNodeRead(d, meta); err != nil {
		return err
	}
	return nil
}

func resourceBizflyCloudCloudDatabaseNodeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	id := d.Id()
	if d.HasChange("volume_size") && d.Get("volume_size").(int) != 0 {
		// retry
		retry := maxRetry
		err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			task, err := client.CloudDatabase.Nodes().ResizeVolume(context.Background(), id, d.Get("volume_size").(int))

			if err != nil {
				retry--
				if retry > 0 {
					time.Sleep(timeSleep)
					return resource.RetryableError(fmt.Errorf("[ERROR] Resize volume of database node [%s] failed: %v. Retrying", id, err))
				}

				return resource.NonRetryableError(fmt.Errorf("[ERROR] Resize volume of database node [%s] failed: %v. Can't retry", id, err))
			}

			_ = d.Set("task_id", task.TaskID)

			err = resourceBizflyCloudCloudDatabaseNodeRead(d, meta)
			if err != nil {
				return resource.NonRetryableError(fmt.Errorf("[ERROR] Resize volume of database node %s failed: %s. Can't retry", id, err))
			}

			// wait for database node is active again
			_, err = waitForCloudDatabaseNodeUpdate(d, meta, "volume_size", d.Get("volume_size").(int))
			if err != nil {
				return resource.NonRetryableError(fmt.Errorf("[ERROR] Resize volume of database node %s with task id (%s) failed: %s. Can't retry", id, task.TaskID, err))
			}

			return nil
		})

		if err != nil {
			return fmt.Errorf("[ERROR] Resize volume of database node %s failed: %s. Can't retry", id, err)
		}
	}

	if d.HasChange("flavor_name") && d.Get("flavor_name").(string) != "" {
		// retry
		retry := maxRetry
		err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			task, err := client.CloudDatabase.Nodes().ResizeFlavor(context.Background(), id, d.Get("flavor_name").(string))

			if err != nil {
				retry--
				if retry > 0 {
					time.Sleep(timeSleep)
					return resource.RetryableError(fmt.Errorf("[ERROR] Resize flavor of database node [%s] failed: %v. Retrying", id, err))
				}

				return resource.NonRetryableError(fmt.Errorf("[ERROR] Resize flavor of database node [%s] failed: %v. Can't retry", id, err))
			}

			_ = d.Set("task_id", task.TaskID)

			err = resourceBizflyCloudCloudDatabaseNodeRead(d, meta)
			if err != nil {
				return resource.NonRetryableError(fmt.Errorf("[ERROR] Resize flavor of database node %s failed: %s. Can't retry", id, err))
			}

			// wait for database node is active again
			_, err = waitForCloudDatabaseNodeUpdate(d, meta, "flavor_name", d.Get("flavor_name").(string))
			if err != nil {
				return resource.NonRetryableError(fmt.Errorf("[ERROR] Resize flavor of database node %s with task id (%s) failed: %s. Can't retry", id, task.TaskID, err))
			}

			return nil
		})
		if err != nil {
			return fmt.Errorf("[ERROR] Resize volume of database node %s failed: %s. Can't retry", id, err)
		}
	}

	return resourceBizflyCloudCloudDatabaseNodeRead(d, meta)
}

func waitForCloudDatabaseNodeUpdate(d *schema.ResourceData, meta interface{}, key string, newValue interface{}) (interface{}, error) {
	log.Printf("[INFO] Waiting for cloud database node (%s) to be update", d.Get("name").(string))
	stateConf := &resource.StateChangeConf{
		Pending:        []string{"false", "RESIZE"},
		Target:         []string{"true", "ACTIVE", "HEALTHY"},
		Refresh:        updateCloudDatabaseNodeStateRefreshFunc(d, key, newValue, meta),
		Timeout:        d.Timeout(schema.TimeoutUpdate),
		Delay:          60 * time.Second,
		MinTimeout:     20 * time.Second,
		NotFoundChecks: 30,
	}
	return stateConf.WaitForState()
}

func waitForCloudDatabaseNodeDelete(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	log.Printf("[INFO] Waiting for cloud database node (%s) to be delete", d.Get("name").(string))
	stateConf := &resource.StateChangeConf{
		Pending:        []string{"false", "SHUTDOWN"},
		Target:         []string{"true"},
		Refresh:        deleteCloudDatabaseNodeStateRefreshFunc(d, meta),
		Timeout:        d.Timeout(schema.TimeoutDelete),
		Delay:          30 * time.Second,
		MinTimeout:     20 * time.Second,
		NotFoundChecks: 30,
	}
	return stateConf.WaitForState()
}

func newCloudDatabaseNodeStateRefreshFunc(d *schema.ResourceData, meta interface{}) resource.StateRefreshFunc {
	client := meta.(*CombinedConfig).gobizflyClient()

	return func() (interface{}, string, error) {
		err := resourceBizflyCloudCloudDatabaseNodeRead(d, meta)
		if err != nil {
			return nil, "false", err
		}

		if attr, ok := d.GetOk("status"); ok {
			node, err := client.CloudDatabase.Nodes().Get(context.Background(), d.Id())
			if err != nil {
				return nil, "", fmt.Errorf("[ERROR] Retrieving cloud database Node %s error: %v", d.Id(), err)
			}
			switch attr := attr.(type) {
			case bool:
				return &node, strconv.FormatBool(attr), nil
			default:
				return &node, attr.(string), nil
			}
		}
		return nil, "false", nil
	}
}

func updateCloudDatabaseNodeStateRefreshFunc(d *schema.ResourceData, key string, newValue interface{}, meta interface{}) resource.StateRefreshFunc {
	client := meta.(*CombinedConfig).gobizflyClient()

	return func() (interface{}, string, error) {
		err := resourceBizflyCloudCloudDatabaseNodeRead(d, meta)
		if err != nil {
			return nil, "", err
		}

		node, err := client.CloudDatabase.Nodes().Get(context.Background(), d.Id())
		if err != nil {
			return nil, "", fmt.Errorf("[ERROR] Retrieving cloud database Node %s error: %v", d.Id(), err)
		}

		switch key {
		case "volume_size":
			if node.Volume.Size < newValue.(int) {
				log.Println("[DEBUG] Cloud database node is updating")
				return nil, "", nil
			}

		case "flavor_name":
			if !strings.Contains(node.Flavor, newValue.(string)) {
				log.Println("[DEBUG] Cloud database node is updating")
				return nil, "", nil
			}
		}

		if attr, ok := d.GetOk("status"); ok {
			node, err = client.CloudDatabase.Nodes().Get(context.Background(), d.Id())
			if err != nil {
				return nil, "", fmt.Errorf("[ERROR] Retrieving cloud database node %s error: %v", d.Id(), err)
			}
			switch attr := attr.(type) {
			case bool:
				return &node, strconv.FormatBool(attr), nil
			default:
				return &node, attr.(string), nil
			}
		}
		return nil, "", nil
	}
}

func deleteCloudDatabaseNodeStateRefreshFunc(d *schema.ResourceData, meta interface{}) resource.StateRefreshFunc {
	client := meta.(*CombinedConfig).gobizflyClient()

	return func() (interface{}, string, error) {
		ins, err := client.CloudDatabase.Nodes().Get(context.Background(), d.Id())

		if errors.Is(err, gobizfly.ErrNotFound) {
			return ins, "true", nil
		} else if err != nil {
			return nil, "", err
		}
		return ins, "false", nil

	}
}
