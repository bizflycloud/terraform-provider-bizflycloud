package bizflycloud

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bizflycloud/gobizfly"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBizflyCloudContainerRegistry() *schema.Resource {
	return &schema.Resource{
		Create: resourceBizflyCloudContainerRegistryCreate,
		Read:   resourceBizflyCloudContainerRegistryRead,
		Update: resourceBizflyCloudContainerRegistryUpdate,
		Delete: resourceBizflyCloudContainerRegistryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"public": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
	}
}

func resourceBizflyCloudContainerRegistryCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*gobizfly.Client)
	name := d.Get("name").(string)
	public := d.Get("public").(bool)

	payload := &gobizfly.CreateRepositoryPayload{
		Name:   name,
		Public: public,
	}

	err := client.ContainerRegistry.Create(context.Background(), payload)
	if err != nil {
		return fmt.Errorf("Error creating container registry: %v", err)
	}

	d.SetId(name)

	return resourceBizflyCloudContainerRegistryRead(d, m)
}

func resourceBizflyCloudContainerRegistryRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*gobizfly.Client)
	name := d.Id()

	opts := &gobizfly.ListOptions{}
	registries, err := client.ContainerRegistry.List(context.Background(), opts)
	if err != nil {
		return fmt.Errorf("Error retrieving container registries: %v", err)
	}

	for _, registry := range registries {
		if registry.Name == name {
			if err := d.Set("name", registry.Name); err != nil {
				return fmt.Errorf("Error setting name: %v", err)
			}
			if err := d.Set("public", registry.Public); err != nil {
				return fmt.Errorf("Error setting public: %v", err)
			}
			if err := d.Set("created_at", registry.CreatedAt); err != nil {
				return fmt.Errorf("Error setting created_at: %v", err)
			}
			return nil
		}
	}

	// If registry is not found, remove it from state
	d.SetId("")
	return nil
}

func resourceBizflyCloudContainerRegistryUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*gobizfly.Client)

	if d.HasChange("public") {
		name := d.Id()
		public := d.Get("public").(bool)

		payload := &gobizfly.EditRepositoryPayload{
			Public: public,
		}

		err := client.ContainerRegistry.EditRepo(context.Background(), name, payload)
		if err != nil {
			return fmt.Errorf("Error updating container registry visibility: %v", err)
		}
	}

	return resourceBizflyCloudContainerRegistryRead(d, m)
}

func resourceBizflyCloudContainerRegistryDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*gobizfly.Client)

	log.Printf("[DEBUG] Deleting container registry: %s", d.Id())
	err := client.ContainerRegistry.Delete(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error deleting container registry: %v", err)
	}

	return nil
}
