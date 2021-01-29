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
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Provider returns a schema.Provider for BizFly Cloud.
func Provider() terraform.ResourceProvider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The URL use for the BizFly Cloud API",
				DefaultFunc: schema.EnvDefaultFunc("BIZFLYCLOUD_API_ENDPOINT", "https://manage.bizflycloud.vn/api"),
			},
			"auth_method": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("BIZFLYCLOUD_AUTH_METHOD", "password"),
				Description: "Authentication method for BizFly Cloud API",
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
				Description: "BizFly Cloud Region Name. Default is HN",
				DefaultFunc: schema.EnvDefaultFunc("BIZFLYCLOUD_REGION_NAME", "HN"),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"bizflycloud_server":                resourceBizFlyCloudServer(),
			"bizflycloud_volume":                resourceBizFlyCloudVolume(),
			"bizflycloud_volume_snapshot":       resourceBizFlyCloudVolumeSnapshot(),
			"bizflycloud_ssh_key":               resourceBizFlyCloudSSHKey(),
			"bizflycloud_firewall":              resourceBizFlyCloudFirewall(),
			"bizflycloud_loadbalancer":          resourceBizFlyCloudLoadBalancer(),
			"bizflycloud_loadbalancer_listener": resourceBizFlyCloudLoadBalancerListener(),
			"bizflycloud_loadbalancer_pool":     resourceBizFlyCloudLoadBalancerPool(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"bizflycloud_image": datasourceBizFlyCloudImages(),
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
	}
	return config.Client()
}
