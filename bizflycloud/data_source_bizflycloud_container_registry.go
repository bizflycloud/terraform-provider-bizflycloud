package bizflycloud

import (
	"context"
	"fmt"

	"github.com/bizflycloud/gobizfly"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceBizflyCloudContainerRegistry() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBizflyCloudContainerRegistryRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceBizflyCloudContainerRegistryRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*gobizfly.Client)
	name := d.Get("name").(string)

	opts := &gobizfly.ListOptions{}
	registries, err := client.ContainerRegistry.List(context.Background(), opts)
	if err != nil {
		return fmt.Errorf("Error retrieving container registries: %v", err)
	}

	for _, registry := range registries {
		if registry.Name == name {
			d.SetId(registry.Name)
			d.Set("name", registry.Name)
			d.Set("created_at", registry.CreatedAt)
			return nil
		}
	}

	return fmt.Errorf("No container registry found with name: %s", name)
} 