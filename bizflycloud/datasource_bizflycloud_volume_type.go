package bizflycloud

import (
	"context"
	"errors"
	"fmt"
	"github.com/bizflycloud/gobizfly"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

func datasourceBizflyCloudVolumeTypes() *schema.Resource {
	return &schema.Resource{
		Read:   datasourceBizflyCloudVolumeTypesRead,
		Schema: dataSourceVolumeTypeSchema(),
	}
}

func datasourceBizflyCloudVolumeTypesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	volumeTypes, err := client.Volume.ListVolumeTypes(context.Background(), &gobizfly.ListVolumeTypesOptions{})
	if err != nil {
		return err
	}
	name := d.Get("name").(string)
	category := d.Get("category").(string)
	err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
		for _, volumeType := range volumeTypes {
			if volumeType.Name == name && volumeType.Category == category {
				d.SetId(volumeType.Type)
				d.Set("name", volumeType.Name)
				d.Set("category", volumeType.Category)
				d.Set("type", volumeType.Type)
				d.Set("availability_zones", volumeType.AvailabilityZones)
				return nil
			}
		}
		return resource.RetryableError(err)
	})
	if !d.IsNewResource() && errors.Is(err, gobizfly.ErrNotFound) {
		log.Printf("[WARN] Volume Type %s is not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("error read Volume Type %s: %w", d.Id(), err)
	}
	return nil
}
