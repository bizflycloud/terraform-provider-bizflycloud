package bizflycloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Provider returns a schema.Provider for BizFly Cloud.
func Provider() terraform.ResourceProvider {
	p := &schema.Provider{}
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
		TerraformVersion:    terraformVersion,
	}
	return config.Client()
}
