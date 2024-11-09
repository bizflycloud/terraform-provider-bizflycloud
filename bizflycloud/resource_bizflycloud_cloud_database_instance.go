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

func resourceBizflyCloudDatabaseInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceBizflyCloudCloudDatabaseInstanceCreate,
		Delete: resourceBizflyCloudCloudDatabaseInstanceDelete,
		Read:   resourceBizflyCloudCloudDatabaseInstanceRead,
		Update: resourceBizflyCloudCloudDatabaseInstanceUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(80 * time.Minute),
		},
		Schema: resourceCloudDatabaseInstanceSchema(),
	}
}

func resourceBizflyCloudCloudDatabaseInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()

	autoscaling := readResourceCloudDatabaseAutoScaling(d)
	datastore := readResourceCloudDatabaseDatastore(d)

	insc := &gobizfly.CloudDatabaseInstanceCreate{
		Name:         d.Get("name").(string),
		InstanceType: d.Get("instance_type").(string),
		FlavorName:   d.Get("flavor_name").(string),
		VolumeSize:   d.Get("volume_size").(int),
		Datastore: gobizfly.CloudDatabaseDatastore{
			Type:      datastore["type"],
			VersionID: datastore["version_id"],
		},
		PublicAccess:     d.Get("public_access").(bool),
		AvailabilityZone: d.Get("availability_zone").(string),
		Networks:         makeArrayNetworkIDFromArray(d.Get("network_ids").(*schema.Set).List()),
		BackupID:         d.Get("backup_id").(string),
	}

	insc.AutoScaling = &gobizfly.CloudDatabaseAutoScaling{
		Enable: false,
		Volume: gobizfly.CloudDatabaseAutoScalingVolume{
			Limited:   autoscaling["volume_limited"],
			Threshold: autoscaling["volume_threshold"],
		},
	}

	if autoscaling["enable"] == 1 {
		insc.AutoScaling.Enable = true
	}

	if _, ok := d.GetOk("secondaries"); ok {
		secondaries := readResourceCloudDatabaseInstanceSecondary(d.Get("secondaries").(*schema.Set))
		// Other has replica known as secondary nodes
		if datastore["type"] == "MariaDB" || datastore["type"] == "MySQL" || datastore["type"] == "MongoDB" || datastore["type"] == "Postgres" || datastore["type"] == "Redis" {
			insc.Secondaries = &gobizfly.CloudDatabaseReplicaNodeCreate{
				Quantity:       secondaries["quantity"].(int),
				Configurations: gobizfly.CloudDatabaseReplicasConfiguration{AvailabilityZone: secondaries["availability_zone"].(string)},
			}
		} else {
			insc.Replicas = &gobizfly.CloudDatabaseReplicaNodeCreate{
				Quantity:       secondaries["quantity"].(int),
				Configurations: gobizfly.CloudDatabaseReplicasConfiguration{AvailabilityZone: secondaries["availability_zone"].(string)},
			}
		}
	}

	instance, err := client.CloudDatabase.Instances().Create(context.Background(), insc)
	if err != nil {
		return fmt.Errorf("[ERROR] create cloud database instance %s failed: %s", d.Get("name"), err)
	}
	log.Printf("[DEBUG] creating cloud database instance %s", instance.Name)

	d.SetId(instance.ID)
	_ = d.Set("task_id", instance.TaskID)

	// wait for cloud database instance to become active
	_, err = waitForCloudDatabaseInstanceCreate(d, meta)
	if err != nil {
		return fmt.Errorf("[ERROR] create cloud database instance (%s) failed: %s", d.Get("name").(string), err)
	}

	if _, ok := d.GetOk("init_databases"); ok {
		// Do create new databases
		initDatabases := d.Get("init_databases").([]interface{})
		newDatabases := []*gobizfly.CloudDatabaseDB{}

		for _, database := range initDatabases {
			newDatabases = append(newDatabases, &gobizfly.CloudDatabaseDB{Name: database.(string)})
		}

		if len(newDatabases) > 0 {
			err := client.CloudDatabase.Instances().CreateDatabases(context.Background(), instance.ID, newDatabases)

			if err != nil {
				return fmt.Errorf("[ERROR] Create new database for database instance [%s] failed: %v", instance.ID, err)
			}
		}
	}

	if _, ok := d.GetOk("users"); ok {
		// Do create new users
		newUsers := readDatabaseUsers(d.Get("users").(*schema.Set))

		err = client.CloudDatabase.Instances().CreateUsers(context.Background(), instance.ID, newUsers)
		if err != nil {
			return fmt.Errorf("[ERROR] Create new user for database instance [%s] failed: %v", instance.ID, err)
		}
	}

	if _, ok := d.GetOk("configuration_group"); ok {
		// Do attach configuration_group to all nodes
		cfg := make(map[string]string)
		for k, v := range d.Get("configuration_group").(map[string]interface{}) {
			cfg[k] = fmt.Sprintf("%v", v)
		}

		ins, err := client.CloudDatabase.Instances().Get(context.Background(), instance.ID)
		if err != nil {
			return fmt.Errorf("[ERROR] Attach configuration group for database instance [%s] failed: %v", instance.ID, err)
		}

		for _, node := range ins.Nodes {
			_, _ = client.CloudDatabase.Configurations().Attach(context.Background(), node.ID, cfg["id"], true)
		}

		if cfg["apply_immediately"] == "true" {
			for _, node := range ins.Nodes {
				_, _ = client.CloudDatabase.Nodes().Restart(context.Background(), node.ID)
			}
		}
	}

	return resourceBizflyCloudCloudDatabaseInstanceRead(d, meta)
}

func resourceBizflyCloudCloudDatabaseInstanceRead(d *schema.ResourceData, meta interface{}) error {
	if err := dataSourceBizflyCloudDatabaseInstanceRead(d, meta); err != nil {
		return err
	}
	return nil
}

func resourceBizflyCloudCloudDatabaseInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	id := d.Id()
	// We need keep current changed values in here because:
	// when do update autoscaling, d schema being set current values
	datastore := readResourceCloudDatabaseDatastore(d)
	instanceType := d.Get("instance_type").(string)
	newVolumeSize := d.Get("volume_size").(int)

	if d.HasChange("autoscaling") {
		autoscaling := readResourceCloudDatabaseAutoScaling(d)
		_ = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			das := &gobizfly.CloudDatabaseAutoScaling{
				Enable: false,
				Volume: gobizfly.CloudDatabaseAutoScalingVolume{
					Limited:   autoscaling["volume_limited"],
					Threshold: autoscaling["volume_threshold"],
				}}

			if autoscaling["enable"] == 1 {
				das.Enable = true
				_, _ = client.CloudDatabase.AutoScalings().Update(context.Background(), id, das)
			} else {
				_, _ = client.CloudDatabase.AutoScalings().Delete(context.Background(), id)
			}

			err := resourceBizflyCloudCloudDatabaseInstanceRead(d, meta)
			if err != nil {
				return resource.NonRetryableError(fmt.Errorf("[ERROR] Update autoscaling volume of database instance %s failed: %s. Can't retry", id, err))
			}

			return nil
		})
	}

	if d.HasChange("volume_size") {
		// retry
		instance, _ := client.CloudDatabase.Instances().Get(context.Background(), id)

		if newVolumeSize < instance.Volume.Size {
			return fmt.Errorf("[ERROR] New volume_size must be greater than %v", instance.Volume.Size)
		}

		retry := maxRetry
		err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			task, err := client.CloudDatabase.Instances().ResizeVolume(context.Background(), id, gobizfly.CloudDatabaseDatastore{
				Type:      datastore["type"],
				VersionID: datastore["version_id"],
			}, instanceType, newVolumeSize)

			if err != nil {
				retry--
				if retry > 0 {
					time.Sleep(timeSleep)
					return resource.RetryableError(fmt.Errorf("[ERROR] Resize volume of database instance [%s] error: %v. Retrying", id, err))
				}

				return resource.NonRetryableError(fmt.Errorf("[ERROR] Resize volume of database instance [%s] error: %v. Can't retry", id, err))
			}

			_ = d.Set("task_id", task.TaskID)

			err = resourceBizflyCloudCloudDatabaseInstanceRead(d, meta)
			if err != nil {
				return resource.NonRetryableError(fmt.Errorf("[ERROR] Resize volume of database instance %s failed: %s. Can't retry", id, err))
			}

			// wait for database instance is active again
			_, err = waitForCloudDatabaseInstanceUpdate(d, meta)
			if err != nil {
				return resource.NonRetryableError(fmt.Errorf("[ERROR] Resize volume of database instance %s with task id (%s) error: %s. Can't retry", id, task.TaskID, err))
			}

			return nil
		})

		if err != nil {
			return fmt.Errorf("[ERROR] Resize volume of database instance %s failed: %s", id, err)
		}
	}

	if d.HasChange("flavor_name") || d.HasChange("instance_type") {
		// retry
		retry := maxRetry

		err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			task, err := client.CloudDatabase.Instances().ResizeFlavor(
				context.Background(), id, gobizfly.CloudDatabaseDatastore{
					Type:      datastore["type"],
					VersionID: datastore["version_id"],
				}, instanceType, d.Get("flavor_name").(string))

			if err != nil {
				retry--
				if retry > 0 {
					time.Sleep(timeSleep)
					return resource.RetryableError(fmt.Errorf("[ERROR] Resize flavor of database instance [%s] failed: %v. Retrying", id, err))
				}

				return resource.NonRetryableError(fmt.Errorf("[ERROR] Resize flavor of database instance [%s] failed: %v. Can't retry", id, err))
			}

			_ = d.Set("task_id", task.TaskID)

			err = resourceBizflyCloudCloudDatabaseInstanceRead(d, meta)
			if err != nil {
				return resource.NonRetryableError(fmt.Errorf("[ERROR] Resize flavor of database instance %s failed: %s. Can't retry", id, err))
			}

			// wait for database instance is active again
			_, err = waitForCloudDatabaseInstanceUpdate(d, meta)
			if err != nil {
				return resource.NonRetryableError(fmt.Errorf("[ERROR] Resize flavor of database instance %s with task (%s) failed: %s. Can't retry", id, task.TaskID, err))
			}
			return nil
		})

		if err != nil {
			return fmt.Errorf("[ERROR] Resize flavor of database instance %s failed: %s", id, err)
		}
	}

	if d.HasChange("init_databases") {
		// We just create new databases without do delete any databases
		// Because delete database is an action really dangerous
		databases, _ := readCurrentDatabases(client, id)
		initDatabases := d.Get("init_databases").([]interface{})
		newDatabases := []*gobizfly.CloudDatabaseDB{}

		for _, database := range initDatabases {
			if len(databases) > 0 {
				_, avail := gobizfly.SliceContains(databases, database)
				if avail {
					continue
				}
			}

			newDatabases = append(newDatabases, &gobizfly.CloudDatabaseDB{Name: database.(string)})
		}

		err := client.CloudDatabase.Instances().CreateDatabases(context.Background(), id, newDatabases)

		if err != nil {
			return fmt.Errorf("[ERROR] Create new database for database instance [%s] failed: %v", id, err)
		}
	}

	if d.HasChange("users") {
		// To handle edge case change password:
		// First, we will do delete current users
		// After, we will do create new users
		newUsers := readDatabaseUsers(d.Get("users").(*schema.Set))

		err := client.CloudDatabase.Instances().DeleteUsers(context.Background(), id, newUsers)
		if err != nil {
			return fmt.Errorf("[ERROR] Create new user for database instance [%s] failed: %v", id, err)
		}

		err = client.CloudDatabase.Instances().CreateUsers(context.Background(), id, newUsers)
		if err != nil {
			return fmt.Errorf("[ERROR] Create new user for database instance [%s] failed: %v", id, err)
		}
	}

	return resourceBizflyCloudCloudDatabaseInstanceRead(d, meta)
}

func resourceBizflyCloudCloudDatabaseInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	id := d.Id()

	// retry
	retry := maxRetry
	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		task, err := client.CloudDatabase.Instances().Delete(context.Background(), id, &gobizfly.CloudDatabaseDelete{})

		if err != nil {
			retry--
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

func waitForCloudDatabaseInstanceUpdate(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	log.Printf("[INFO] Waiting for cloud database instance (%s) to be update", d.Get("name").(string))
	stateConf := &resource.StateChangeConf{
		Pending:        []string{"false", "RESIZE", "RESTART_REQUIRED", "REBOOTING"},
		Target:         []string{"true", "ACTIVE", "HEALTHY"},
		Refresh:        newCloudDatabaseInstanceStateRefreshFunc(d, meta),
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
		err := resourceBizflyCloudCloudDatabaseInstanceRead(d, meta)
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
		err := resourceBizflyCloudCloudDatabaseInstanceRead(d, meta)
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

func deleteCloudDatabaseInstanceStateRefreshFunc(d *schema.ResourceData, meta interface{}) resource.StateRefreshFunc {
	client := meta.(*CombinedConfig).gobizflyClient()

	return func() (interface{}, string, error) {
		// Check stask status
		taskID := d.Get("task_id").(string)
		task, err := client.CloudDatabase.Tasks().Get(context.Background(), taskID)
		if err != nil {
			return nil, "false", err
		}

		if task.Ready {
			return nil, "false", fmt.Errorf("instance is deleting")
		}

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

func readResourceCloudDatabaseDatastore(d *schema.ResourceData) map[string]string {
	datastore := make(map[string]string)
	for k, v := range d.Get("datastore").(map[string]interface{}) {
		datastore[k] = fmt.Sprintf("%v", v)
	}
	return datastore
}

func readResourceCloudDatabaseAutoScaling(d *schema.ResourceData) map[string]int {
	autoscaling := make(map[string]int)
	for k, v := range d.Get("autoscaling").(map[string]interface{}) {
		autoscaling[k] = v.(int)
	}
	return autoscaling
}

func readResourceCloudDatabaseInstanceSecondary(secondaries *schema.Set) map[string]interface{} {
	results := make(map[string]interface{})
	for _, _s := range secondaries.List() {
		secondary := _s.(map[string]interface{})

		results["quantity"] = secondary["quantity"].(int)
		results["availability_zone"] = secondary["availability_zone"].(string)
	}

	return results
}

func readCurrentDatabases(client *gobizfly.Client, instanceID string) ([]interface{}, error) {
	databases, err := client.CloudDatabase.Instances().ListDatabases(context.Background(), instanceID)
	if err != nil {
		return nil, err
	}

	var results []interface{}
	for _, item := range databases {
		results = append(results, item.Name)
	}

	return results, nil
}

func readCurrentUsers(client *gobizfly.Client, instanceID string) ([]interface{}, error) {
	users, err := client.CloudDatabase.Instances().ListUsers(context.Background(), instanceID)
	if err != nil {
		return nil, err
	}

	var results []interface{}
	for _, item := range users {
		results = append(results, item.Name)
	}

	return results, nil
}

func readDatabaseUsers(users *schema.Set) []*gobizfly.CloudDatabaseUser {
	var results []*gobizfly.CloudDatabaseUser

	for _, item := range users.List() {
		user := item.(map[string]interface{})
		cdb := make([]gobizfly.CloudDatabaseDB, 0)

		cdu := gobizfly.CloudDatabaseUser{
			Databases: []gobizfly.CloudDatabaseDB{
				gobizfly.CloudDatabaseDB{Name: "yanfei"},
				gobizfly.CloudDatabaseDB{Name: "genshin"},
			},
			Name:     user["username"].(string),
			Password: user["password"].(string),
		}

		if user["host"].(string) != "" {
			cdu.Host = user["host"].(string)
		}

		databases := user["databases"].([]interface{})
		if len(databases) > 0 {
			for _, db := range databases {
				cdb = append(cdb, gobizfly.CloudDatabaseDB{Name: db.(string)})
			}
			cdu.Databases = cdb
		}

		results = append(results, &cdu)
	}

	return results
}
