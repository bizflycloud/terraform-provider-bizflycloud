// This file is part of terraform-provider-bizflycloud
//
// Copyright (C) 2026  Bizfly Cloud
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
	"time"

	"github.com/bizflycloud/gobizfly"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBizflyCloudKafka() *schema.Resource {
	return &schema.Resource{
		Create: resourceBizflyCloudKafkaCreate,
		Read:   resourceBizflyCloudKafkaRead,
		Update: resourceBizflyCloudKafkaUpdate,
		Delete: resourceBizflyCloudKafkaDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
		},
		SchemaVersion: 1,
		Schema: resourceKafkaSchema(),
	}
}

func resourceBizflyCloudKafkaCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()

	req := &gobizfly.KafkaInitClusterRequest{
		ClusterName:      d.Get("name").(string),
		VersionID:        d.Get("version_id").(string),
		Nodes:            d.Get("nodes").(int),
		Flavor:           d.Get("flavor").(string),
		VolumeSize:       d.Get("volume_size").(int),
		VPCNetworkID:     d.Get("vpc_network_id").(string),
		AvailabilityZone: d.Get("availability_zone").(string),
		PublicAccess:     d.Get("public_access").(bool),
	}

	_, err := client.Kafka.Create(context.Background(), req)
	if err != nil {
		return fmt.Errorf("Error creating kafka cluster: %w", err)
	}

	// Wait for cluster to be created and become Active
	// The wait function will set the cluster ID when found
	_, waitErr := waitForKafkaClusterProvision(d, meta)
	if waitErr != nil {
		return fmt.Errorf("error waiting for kafka cluster creation: %w", waitErr)
	}

	return resourceBizflyCloudKafkaRead(d, meta)
}

func resourceBizflyCloudKafkaRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	if d.Id() == "" {
		return nil
	}

	cluster, err := client.Kafka.Get(context.Background(), d.Id())
	if err != nil {
		if errors.Is(err, gobizfly.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving kafka cluster %s: %w", d.Id(), err)
	}

	_ = d.Set("name", cluster.Name)
	_ = d.Set("version_id", cluster.KafkaVersion)
	_ = d.Set("nodes", len(cluster.Nodes))
	_ = d.Set("flavor", cluster.Flavor)
	_ = d.Set("volume_size", cluster.VolumeSize)
	_ = d.Set("status", cluster.Status)
	_ = d.Set("created_at", cluster.CreatedAt)
	_ = d.Set("availability_zone", cluster.AvailabilityZone)
	// ClusterResponse returned by Get does not include VPCNetworkID field
	_ = d.Set("vpc_network_id", "")
	_ = d.Set("public_access", cluster.PublicAccess)
	_ = d.Set("obs_dashboard_url", cluster.OBS.DashboardURL)

	return nil
}

func resourceBizflyCloudKafkaUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	clusterID := d.Id()
	if clusterID == "" {
		return fmt.Errorf("missing kafka cluster id for update")
	}

	// Handle flavor change (resize by flavor)
	if d.HasChange("flavor") {
		_, newVal := d.GetChange("flavor")
		payload := &gobizfly.KafkaResizeClusterRequest{
			Type:   "flavor",
			Flavor: newVal.(string),
		}
		_, err := client.Kafka.Resize(context.Background(), clusterID, payload)
		if err != nil {
			return fmt.Errorf("error resizing kafka cluster flavor: %w", err)
		}
	}

	// Handle volume size change
	if d.HasChange("volume_size") {
		_, newVal := d.GetChange("volume_size")
		payload := &gobizfly.KafkaResizeClusterRequest{
			Type:       "volume",
			VolumeSize: newVal.(int),
		}
		_, err := client.Kafka.Resize(context.Background(), clusterID, payload)
		if err != nil {
			return fmt.Errorf("error resizing kafka cluster volume: %w", err)
		}
	}

	// Handle nodes change (only support adding nodes)
	if d.HasChange("nodes") {
		oldVal, newVal := d.GetChange("nodes")
		oldNodes := oldVal.(int)
		newNodes := newVal.(int)
		if newNodes > oldNodes {
			addReq := &gobizfly.KafkaAddNodeRequest{Nodes: newNodes - oldNodes, Type: "increase"}
			_, err := client.Kafka.AddNode(context.Background(), clusterID, addReq)
			if err != nil {
				return fmt.Errorf("error adding node(s) to kafka cluster: %w", err)
			}
		} else if newNodes < oldNodes {
			return fmt.Errorf("decreasing kafka nodes is not supported via terraform update")
		}
	}

	// Wait for any updates to finish
	_, waitErr := waitForKafkaClusterProvision(d, meta)
	if waitErr != nil {
		return fmt.Errorf("error waiting for kafka cluster update: %w", waitErr)
	}

	return resourceBizflyCloudKafkaRead(d, meta)
}

func resourceBizflyCloudKafkaDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	if d.Id() == "" {
		return nil
	}
	_, err := client.Kafka.Delete(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error deleting kafka cluster %s: %w", d.Id(), err)
	}
	// Wait for deletion
	err = waitForKafkaClusterDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("error waiting for kafka cluster deletion: %w", err)
	}
	return nil
}

func waitForKafkaClusterProvision(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	name := d.Get("name").(string)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Creating", "Pending", "Provisioning", "PENDING_PROVISION", "PROVISIONING", "Resizing", "RESIZING"},
		Target:     []string{"Active", "Failed", "Error", "PROVISIONED", "ACTIVE"},
		Refresh:    newKafkaStatusRefreshFunc(d, name, meta),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	return stateConf.WaitForState()
}

func newKafkaStatusRefreshFunc(d *schema.ResourceData, name string, meta interface{}) resource.StateRefreshFunc {
	client := meta.(*CombinedConfig).gobizflyClient()
	return func() (interface{}, string, error) {
		clusters, err := client.Kafka.List(context.Background(), &gobizfly.KafkaClusterListOptions{Name: name})
		if err != nil {
			return nil, "", err
		}
		for _, c := range clusters {
			if c.Name == name {
				// Set the ID when cluster is found
				if d.Id() == "" {
					d.SetId(c.ID)
				}
				return c, c.Status, nil
			}
		}
		// Not found yet
		return nil, "", nil
	}
}

func waitForKafkaClusterDeleted(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Destroying", "Deleting", "DESTROYING", "DELETING"},
		Target:     []string{},
		Refresh:    newKafkaStatusRefreshFunc(d, name, meta),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, err := stateConf.WaitForState()
	return err
}