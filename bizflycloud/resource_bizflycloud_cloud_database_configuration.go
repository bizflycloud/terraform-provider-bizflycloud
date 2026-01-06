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
	"time"

	"github.com/bizflycloud/gobizfly"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBizflyCloudDatabaseConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceBizflyCloudCloudDatabaseConfigurationCreate,
		Delete: resourceBizflyCloudCloudDatabaseConfigurationDelete,
		Read:   resourceBizflyCloudCloudDatabaseConfigurationRead,
		Update: resourceBizflyCloudCloudDatabaseConfigurationUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: resourceCloudDatabaseConfigurationSchema(),
	}
}

func resourceBizflyCloudCloudDatabaseConfigurationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	datastore := readResourceCloudDatabaseDatastore(d)

	cfc := &gobizfly.CloudDatabaseConfigurationCreate{
		Name: d.Get("name").(string),
		Datastore: gobizfly.CloudDatabaseDatastore{
			ID:        datastore["id"],
			VersionID: datastore["version_id"],
		},
		Parameters: readArrayParameters(d.Get("parameter").(*schema.Set)),
	}

	// retry
	retry := maxRetry
	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		configuration, err := client.CloudDatabase.Configurations().Create(context.Background(), cfc)

		if err != nil {
			retry--
			if retry > 0 {
				time.Sleep(timeSleep)
				return resource.RetryableError(fmt.Errorf("[ERROR] create cloud database configuration %s failed: %s. Retrying", d.Get("name"), err))
			}

			return resource.NonRetryableError(fmt.Errorf("[ERROR] create cloud database configuration %s failed: %s. Can't retry", d.Get("name"), err))
		}

		log.Printf("[DEBUG] creating cloud database configuration %s", configuration.Name)
		d.SetId(configuration.ID)

		// wait for cloud database Configuration to become active
		_, err = waitForCloudDatabaseConfigurationCreate(d, meta)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("[ERROR] create cloud database configuration (%s) failed: %s. Can't retry", d.Get("name").(string), err))
		}

		err = resourceBizflyCloudCloudDatabaseConfigurationRead(d, meta)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("[ERROR] create cloud database Configuration (%s) failed: %s. Can't retry", d.Get("name").(string), err))
		}

		return nil
	})
}

func resourceBizflyCloudCloudDatabaseConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()

	if v, ok := d.GetOk("id"); ok {
		d.SetId(v.(string))
	}

	configurationID := d.Id()

	log.Printf("[DEBUG] Reading database Configuration: %s", configurationID)
	configuration, err := client.CloudDatabase.Configurations().Get(context.Background(), configurationID)

	log.Printf("[DEBUG] Checking for error: %s", err)
	if err != nil {
		return fmt.Errorf("error describing database Configuration: %w", err)
	}

	log.Printf("[DEBUG] Found database Configuration: %s", configurationID)
	log.Printf("[DEBUG] bizflycloud_cloud_database_Configuration - Single database Configuration found: %s", configuration.Name)

	d.SetId(configuration.ID)
	_ = d.Set("created_at", configuration.CreatedAt)
	_ = d.Set("name", configuration.Name)
	_ = d.Set("nodes", configuration.Nodes)
	if err := d.Set("datastore", FlattenStruct(configuration.Datastore)); err != nil {
		return fmt.Errorf("error setting datastore: %w", err)
	}

	return nil
}

func resourceBizflyCloudCloudDatabaseConfigurationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	datastore := readResourceCloudDatabaseDatastore(d)
	id := d.Id()

	if d.HasChange("parameter") {
		cfu := &gobizfly.CloudDatabaseConfigurationUpdate{
			Datastore: gobizfly.CloudDatabaseDatastore{
				ID:        datastore["id"],
				VersionID: datastore["version_id"],
			},
			Parameters: readArrayParameters(d.Get("parameter").(*schema.Set)),
		}

		// retry
		retry := maxRetry
		err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			task, err := client.CloudDatabase.Configurations().Update(context.Background(), id, cfu)

			if err != nil {
				retry--
				if retry > 0 {
					time.Sleep(timeSleep)
					return resource.RetryableError(fmt.Errorf("[ERROR] update cloud database configuration [%s] failed: %v. Retrying", id, err))
				}

				return resource.NonRetryableError(fmt.Errorf("[ERROR] update cloud database configuration [%s] failed: %v. Can't retry", id, err))
			}

			_ = d.Set("task_id", task.TaskID)

			err = resourceBizflyCloudCloudDatabaseConfigurationRead(d, meta)
			if err != nil {
				return resource.NonRetryableError(fmt.Errorf("[ERROR] update cloud database configuration %s failed: %s. Can't retry", id, err))
			}
			return nil
		})

		if err != nil {
			return fmt.Errorf("[ERROR] update cloud database configuration %s failed: %s. Can't retry", id, err)
		}
	}
	return nil
}

func resourceBizflyCloudCloudDatabaseConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	id := d.Id()

	// First, detach configuration from all attached nodes
	log.Printf("[DEBUG] Checking for nodes attached to configuration %s before deletion", id)
	configuration, err := client.CloudDatabase.Configurations().Get(context.Background(), id)
	if err == nil && len(configuration.Nodes) > 0 {
		log.Printf("[WARN] Configuration %s is attached to %d nodes, attempting to detach before deletion", id, len(configuration.Nodes))

		for _, node := range configuration.Nodes {
			log.Printf("[DEBUG] Detaching configuration %s from node %s (%s)", id, node.ID, node.Name)
			_, detachErr := client.CloudDatabase.Configurations().Detach(context.Background(), node.ID, id, false)
			if detachErr != nil {
				// Skip if resource not found (already detached or deleted)
				if errors.Is(detachErr, gobizfly.ErrNotFound) {
					log.Printf("[DEBUG] Configuration %s or node %s not found, skipping detach", id, node.ID)
					continue
				}
				// Just log warning for other errors, don't fail - node might be in transition state
				log.Printf("[WARN] Failed to detach configuration %s from node %s: %v. Continuing with deletion anyway.", id, node.ID, detachErr)
			} else {
				log.Printf("[DEBUG] Successfully detached configuration %s from node %s", id, node.ID)
			}
		}
	}

	// retry
	retry := maxRetry
	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		_, err := client.CloudDatabase.Configurations().Delete(context.Background(), id)

		if err != nil {
			retry--
			if retry > 0 {
				time.Sleep(timeSleep)
				return resource.RetryableError(fmt.Errorf("[ERROR] delete cloud database configuration %s failed: %v. Retrying", id, err))
			}

			return resource.NonRetryableError(fmt.Errorf("[ERROR] delete cloud database configuration %s failed: %v. Can't retry", id, err))
		}

		// wait for cloud database Configuration to delete
		_, err = waitForCloudDatabaseConfigurationDelete(d, meta)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("[ERROR] delete cloud database configuration %s failed: %v. Can't retry", id, err))
		}

		return nil
	})
}

func waitForCloudDatabaseConfigurationCreate(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	log.Printf("[INFO] Waiting for cloud database configuration (%s) to be ready", d.Get("name").(string))
	stateConf := &resource.StateChangeConf{
		Pending:        []string{"false"},
		Target:         []string{"true"},
		Refresh:        newCloudDatabaseConfigurationStateRefreshFunc(d, meta),
		Timeout:        d.Timeout(schema.TimeoutCreate),
		Delay:          10 * time.Second,
		MinTimeout:     20 * time.Second,
		NotFoundChecks: 20,
	}
	return stateConf.WaitForState()
}

func waitForCloudDatabaseConfigurationDelete(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	log.Printf("[INFO] Waiting for cloud database configuration (%s) to be delete", d.Get("name").(string))
	stateConf := &resource.StateChangeConf{
		Pending:        []string{"false"},
		Target:         []string{"true"},
		Refresh:        deleteCloudDatabaseConfigurationStateRefreshFunc(d, meta),
		Timeout:        d.Timeout(schema.TimeoutCreate),
		Delay:          10 * time.Second,
		MinTimeout:     20 * time.Second,
		NotFoundChecks: 20,
	}
	return stateConf.WaitForState()
}

func newCloudDatabaseConfigurationStateRefreshFunc(d *schema.ResourceData, meta interface{}) resource.StateRefreshFunc {
	client := meta.(*CombinedConfig).gobizflyClient()

	return func() (interface{}, string, error) {
		err := resourceBizflyCloudCloudDatabaseConfigurationRead(d, meta)
		if err != nil {
			return nil, "", err
		}
		conf, err := client.CloudDatabase.Configurations().Get(context.Background(), d.Id())
		if err != nil {
			return nil, "", fmt.Errorf("[ERROR] Retrieving cloud database configuration %s error: %v", d.Id(), err)
		}
		return &conf, "true", nil
	}

}

func deleteCloudDatabaseConfigurationStateRefreshFunc(d *schema.ResourceData, meta interface{}) resource.StateRefreshFunc {
	client := meta.(*CombinedConfig).gobizflyClient()

	return func() (interface{}, string, error) {
		conf, err := client.CloudDatabase.Configurations().Get(context.Background(), d.Id())

		if errors.Is(err, gobizfly.ErrNotFound) {
			return &conf, "true", nil
		}

		return &conf, "", nil

	}
}

func readArrayParameters(params *schema.Set) map[string]interface{} {
	results := make(map[string]interface{})
	for _, param := range params.List() {
		_param := param.(map[string]interface{})
		key := _param["name"].(string)
		val := _param["value"]

		if key == "" || val == "" {
			continue
		}

		if convertedValue, err := strconv.ParseBool(val.(string)); err == nil {
			val = convertedValue
		} else if convertedValue, err := strconv.Atoi(val.(string)); err == nil {
			val = convertedValue
		} else if convertedValue, err := strconv.ParseFloat(val.(string), 32); err == nil {
			val = convertedValue
		}

		results[key] = val
	}
	return results
}
