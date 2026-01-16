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
		return fmt.Errorf("error retrieving container registries: %v", err)
	}

	for _, registry := range registries {
		if registry.Name == name {
			d.SetId(registry.Name)
			if err := d.Set("name", registry.Name); err != nil {
				return fmt.Errorf("error setting name: %v", err)
			}
			if err := d.Set("public", registry.Public); err != nil {
				return fmt.Errorf("error setting public: %v", err)
			}
			if err := d.Set("created_at", registry.CreatedAt); err != nil {
				return fmt.Errorf("error setting created_at: %v", err)
			}
			return nil
		}
	}

	return fmt.Errorf("no container registry found with name: %s", name)
}
