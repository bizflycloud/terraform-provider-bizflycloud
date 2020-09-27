package main

import (
	"github.com/bizflycloud/terraform-provider-bizflycloud/bizflycloud"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: bizflycloud.Provider})
}
