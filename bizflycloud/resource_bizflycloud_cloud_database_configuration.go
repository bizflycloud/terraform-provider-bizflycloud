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

func resourceBizFlyCloudDatabaseConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceBizFlyCloudCloudDatabaseConfigurationCreate,
		Read:   resourceBizFlyCloudCloudDatabaseConfigurationRead,
		Update: resourceBizFlyCloudCloudDatabaseConfigurationUpdate,
		Delete: resourceBizFlyCloudCloudDatabaseConfigurationDelete,
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

func resourceBizFlyCloudCloudDatabaseConfigurationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()

	cfc := &gobizfly.CloudDatabaseConfigurationCreate{
		ConfigurationName:       d.Get("name").(string),
		ConfigurationParameters: readArrayParameters(d.Get("parameters").(map[string]interface{})),
		Datastore: gobizfly.CloudDatabaseDatastore{
			Type: d.Get("datastore_type").(string),
			Name: d.Get("datastore_version_name").(string),
			ID:   d.Get("datastore_version_id").(string),
		},
	}

	// retry
	retry := maxRetry
	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		configuration, err := client.CloudDatabase.Configurations().Create(context.Background(), cfc)

		if err != nil {
			retry -= 1
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

		err = resourceBizFlyCloudCloudDatabaseConfigurationRead(d, meta)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("[ERROR] create cloud database Configuration (%s) failed: %s. Can't retry", d.Get("name").(string), err))
		}

		return nil
	})
}

func resourceBizFlyCloudCloudDatabaseConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	if err := dataSourceBizFlyCloudDatabaseConfigurationRead(d, meta); err != nil {
		return err
	}
	return nil
}

func resourceBizFlyCloudCloudDatabaseConfigurationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	id := d.Id()
	if d.HasChange("parameters") {
		cfu := &gobizfly.CloudDatabaseConfigurationUpdate{
			ConfigurationParameters: readArrayParameters(d.Get("parameters").(map[string]interface{})),
		}

		// retry
		retry := maxRetry
		err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			task, err := client.CloudDatabase.Configurations().Update(context.Background(), id, cfu)

			if err != nil {
				retry -= 1
				if retry > 0 {
					time.Sleep(timeSleep)
					return resource.RetryableError(fmt.Errorf("[ERROR] update cloud database configuration [%s] failed: %v. Retrying", id, err))
				}

				return resource.NonRetryableError(fmt.Errorf("[ERROR] update cloud database configuration [%s] failed: %v. Can't retry", id, err))
			}

			_ = d.Set("task_id", task.TaskID)

			err = resourceBizFlyCloudCloudDatabaseConfigurationRead(d, meta)
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

func resourceBizFlyCloudCloudDatabaseConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	id := d.Id()

	// retry
	retry := maxRetry
	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		_, err := client.CloudDatabase.Configurations().Delete(context.Background(), id)

		if err != nil {
			retry -= 1
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
		err := resourceBizFlyCloudCloudDatabaseConfigurationRead(d, meta)
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

func readArrayParameters(paramMap map[string]interface{}) map[string]interface{} {
	parameterInterface := make(map[string]interface{})
	for key, val := range paramMap {
		if key == "" || val == "" {
			continue
		}

		if convertedValue, err := strconv.ParseBool(val.(string)); err == nil {
			val = convertedValue
		} else if convertedValue, err := strconv.ParseFloat(val.(string), 32); err == nil {
			val = convertedValue
		}
		parameterInterface[key] = val
	}
	return parameterInterface
}
