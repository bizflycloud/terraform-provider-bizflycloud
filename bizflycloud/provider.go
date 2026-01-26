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
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Provider returns a schema.Provider for Bizfly Cloud.
func Provider() terraform.ResourceProvider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The URL use for the Bizfly Cloud API",
				DefaultFunc: schema.EnvDefaultFunc("BIZFLYCLOUD_API_ENDPOINT", "https://manage.bizflycloud.vn/api"),
			},
			"auth_method": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("BIZFLYCLOUD_AUTH_METHOD", "password"),
				Description: "Authentication method for Bizfly Cloud API",
			},
			"email": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Email to authenticate",
				DefaultFunc: schema.EnvDefaultFunc("BIZFLYCLOUD_EMAIL", nil),
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Password for email with auth_method password",
				DefaultFunc: schema.EnvDefaultFunc("BIZFLYCLOUD_PASSWORD", nil),
			},
			"application_credential_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Application credential ID for authenticate use application_credential",
				DefaultFunc: schema.EnvDefaultFunc("BIZFLYCLOUD_APPLICATION_CREDENTIAL_ID", nil),
			},
			"application_credential_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Application credential secret for authenticate use application_credential",
				DefaultFunc: schema.EnvDefaultFunc("BIZFLYCLOUD_APPLICATION_CREDENTIAL_SECRET", nil),
			},
			"region_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Bizfly Cloud Region Name. Default is HaNoi",
				DefaultFunc: schema.EnvDefaultFunc("BIZFLYCLOUD_REGION_NAME", "HaNoi"),
			},
			"project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Bizfly Cloud Project ID",
				DefaultFunc: schema.EnvDefaultFunc("BIZFLYCLOUD_PROJECT_ID", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"bizflycloud_server":                               resourceBizflyCloudServer(),
			"bizflycloud_volume":                               resourceBizflyCloudVolume(),
			"bizflycloud_volume_snapshot":                      resourceBizflyCloudVolumeSnapshot(),
			"bizflycloud_ssh_key":                              resourceBizflyCloudSSHKey(),
			"bizflycloud_firewall":                             resourceBizflyCloudFirewall(),
			"bizflycloud_loadbalancer":                         resourceBizflyCloudLoadBalancer(),
			"bizflycloud_loadbalancer_listener":                resourceBizflyCloudLoadBalancerListener(),
			"bizflycloud_loadbalancer_l7policy":                resourceBizflyCloudLoadBalancerL7Policy(),
			"bizflycloud_loadbalancer_pool":                    resourceBizflyCloudLoadBalancerPool(),
			"bizflycloud_simple_storage_bucket":                resourceBizflyCloudSimpleStorageBucket(),
			"bizflycloud_simple_storage_access_key":            resourceBizflyCloudSimpleStorageAccessKey(),
			"bizflycloud_simple_storage_bucket_acl":            resourceBizflyCloudSimpleStorageBucketAcl(),
			"bizflycloud_simple_storage_bucket_versioning":     resourceBizflyCloudSimpleStorageBucketVersioning(),
			"bizflycloud_simple_storage_bucket_cors":           resourceBizflyCloudSimpleStorageBucketCors(),
			"bizflycloud_simple_storage_bucket_website_config": resourceBizflyCloudSimpleStorageBucketWebsiteConfig(),
			"bizflycloud_autoscaling_group":                    resourceBizflyCloudAutoscalingGroup(),
			"bizflycloud_autoscaling_scalein_policy":           resourceBizflyCloudAutoscalingScaleInPolicy(),
			"bizflycloud_autoscaling_scaleout_policy":          resourceBizflyCloudAutoscalingScaleOutPolicy(),
			"bizflycloud_autoscaling_deletion_policy":          resourceBizflyCloudAutoscalingDeletionPolicy(),
			"bizflycloud_autoscaling_launch_configuration":     resourceBizflyCloudAutoscalingLaunchConfiguration(),
			"bizflycloud_kubernetes":                           resourceBizflyCloudKubernetes(),
			"bizflycloud_kubernetes_worker_pool":               resourceBizflyCloudKubernetesWorkerPool(),
			"bizflycloud_vpc_network":                          resourceBizflyCloudVPCNetwork(),
			"bizflycloud_network_interface":                    resourceBizflyCloudNetworkInterface(),
			"bizflycloud_dns":                                  resourceBizflyCloudDNS(),
			"bizflycloud_wan_ip":                               resourceBizflyCloudWanIP(),
			"bizflycloud_scheduled_volume_backup":              resourceBizflyCloudScheduledVolumeBackup(),
			"bizflycloud_cloud_database_backup":                resourceBizflyCloudDatabaseBackup(),
			"bizflycloud_cloud_database_backup_schedule":       resourceBizflyCloudDatabaseBackupSchedule(),
			"bizflycloud_cloud_database_configuration":         resourceBizflyCloudDatabaseConfiguration(),
			"bizflycloud_cloud_database_instance":              resourceBizflyCloudDatabaseInstance(),
			"bizflycloud_custom_image":                         resourceBizflyCloudCustomImage(),
			"bizflycloud_volume_attachment":                    resourceBizflyCloudVolumeAttachment(),
			"bizflycloud_cdn":                                  resourceBizflyCloudCDN(),
			"bizflycloud_internet_gateway":                     resourceInternetGateway(),
			"bizflycloud_container_registry":                   resourceBizflyCloudContainerRegistry(),
			"bizflycloud_kafka":                                resourceBizflyCloudKafka(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"bizflycloud_image":                            datasourceBizflyCloudImages(),
			"bizflycloud_autoscaling_group":                datasourceBizflyCloudAutoScalingGroup(),
			"bizflycloud_autoscaling_launch_configuration": datasourceBizflyCloudLaunchConfiguration(),
			"bizflycloud_vpc_network":                      dataSourceBizflyCloudVPCNetwork(),
			"bizflycloud_kubernetes_version":               datasourceBizflyCloudKubernetesControllerVersions(),
			"bizflycloud_kubernetes_package":               datasourceBizflyCloudKubernetesControllerPackage(),
			"bizflycloud_network_interface":                dataSourceBizflyCloudNetworkInterface(),
			"bizflycloud_server":                           datasourceBizflyCloudServers(),
			"bizflycloud_autoscaling_nodes":                datasourceBizflyCloudAutoscalingNodes(),
			"bizflycloud_ssh_key":                          dataSourceBizflyCloudSSHKey(),
			"bizflycloud_wan_ip":                           dataSourceBizflyCloudWanIP(),
			"bizflycloud_server_type":                      dataSourceBizflyCloudServerTypes(),
			"bizflycloud_volume_type":                      datasourceBizflyCloudVolumeTypes(),
			"bizflycloud_cloud_database_backup":            datasourceBizflyCloudDatabaseBackup(),
			"bizflycloud_cloud_database_datastore":         datasourceBizflyCloudDatabaseDatastore(),
			"bizflycloud_cloud_database_instance":          datasourceBizflyCloudDatabaseInstance(),
			"bizflycloud_cloud_database_node":              datasourceBizflyCloudDatabaseNode(),
			"bizflycloud_custom_image":                     dataSourceBizflyCloudCustomImage(),
			"bizflycloud_volume_snapshot":                  dataSourceBizflyCloudVolumeSnapshot(),
			"bizflycloud_container_registry":               dataSourceBizflyCloudContainerRegistry(),
			"bizflycloud_kafka":                            dataSourceBizflyCloudKafka(),
			"bizflycloud_kafka_version":                    dataSourceBizflyCloudKafkaVersion(),
			"bizflycloud_kafka_flavor":                     dataSourceBizflyCloudKafkaFlavor(),
		},
	}
	p.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		terraformVersion := p.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}
		return providerConfigure(d, terraformVersion)
	}
	return p
}

func providerConfigure(d *schema.ResourceData, terraformVersion string) (interface{}, error) {
	config := Config{
		APIEndpoint:         d.Get("api_endpoint").(string),
		AuthMethod:          d.Get("auth_method").(string),
		Email:               d.Get("email").(string),
		Password:            d.Get("password").(string),
		AppCredentialID:     d.Get("application_credential_id").(string),
		AppCredentialSecret: d.Get("application_credential_secret").(string),
		RegionName:          d.Get("region_name").(string),
		TerraformVersion:    terraformVersion,
		ProjectID:           d.Get("project_id").(string),
	}
	combinedClient, err := config.Client()
	if err != nil {
		return nil, err
	}
	return combinedClient, nil
}
