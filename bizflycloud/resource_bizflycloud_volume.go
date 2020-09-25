package bizflycloud

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/bizflycloud/gobizfly"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBizFlyCloudVolume() *schema.Resource {
	return &schema.Resource{
		Create:        resourceBizFlyCloudVolumeCreate,
		Read:          resourceBizFlyCloudVolumeRead,
		Update:        resourceBizFlyCloudVolumeUpdate,
		Delete:        resourceBizFlyCloudVolumeDelete,
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"category": {
				Type:     schema.TypeString,
				Required: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceBizFlyCloudVolumeCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	vcr := &gobizfly.VolumeCreateRequest{
		Name:             d.Get("name").(string),
		Size:             d.Get("size").(int),
		VolumeType:       d.Get("type").(string),
		VolumeCategory:   d.Get("category").(string),
		AvailabilityZone: d.Get("availability_zone").(string),
	}
	log.Printf("[DEBUG] Create Volume configuration #{vcr}")
	volume, err := client.Volume.Create(context.Background(), vcr)
	if err != nil {
		fmt.Errorf("Error creating volume: #{err}")
		return nil
	}
	d.SetId(volume.ID)
	err = resourceBizFlyCloudVolumeRead(d, meta)
	if err != nil {
		fmt.Printf("Error retrieving volume: ${err}")
	}
	return nil
}

func resourceBizFlyCloudVolumeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	volume, err := client.Volume.Get(context.Background(), d.Id())
	if err != nil {
		if errors.Is(err, gobizfly.ErrNotFound) {
			log.Printf("[WARN] VOlume (%s) is not found", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving volume: #{err}")
	}
	d.Set("name", volume.Name)
	d.Set("size", volume.Size)
	d.Set("status", volume.Status)
	d.Set("type", volume.VolumeType)
	d.Set("category", volume.Category)
	d.Set("created_at", volume.CreatedAt)
	d.Set("availability_zone", volume.AvailabilityZone)
	d.Set("user_id", volume.UserID)
	d.Set("project_id", volume.ProjectID)
	return nil
}

func resourceBizFlyCloudVolumeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	if d.HasChange("size") {
		// resize volume
		// TODO check state of extending task
		_, err := client.Volume.ExtendVolume(context.Background(), d.Id(), d.Get("size").(int))
		if err != nil {
			if errors.Is(err, gobizfly.ErrNotFound) {
				log.Printf("[WARN] Volume is not found %s", d.Id())
				d.SetId("")
				return nil
			}
		}
	}
	return nil
}

func resourceBizFlyCloudVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	err := client.Volume.Delete(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error delete volume: #{err}")
	}
	return nil
}
