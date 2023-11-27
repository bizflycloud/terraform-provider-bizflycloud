package bizflycloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceBizflyCloudVolumeSnapshot() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"volume_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
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
			"snapshot_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Read: dataSourceBizflyCloudVolumeSnapshotRead,
	}
}

func dataSourceBizflyCloudVolumeSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()

	snapshot, err := client.Snapshot.Get(context.Background(), d.Get("id").(string))
	if err != nil {
		return err
	}

	d.SetId(snapshot.Id)
	err = d.Set("name", snapshot.Name)
	if err != nil {
		return err
	}
	err = d.Set("volume_id", snapshot.VolumeId)
	if err != nil {
		return err
	}
	err = d.Set("size", snapshot.Size)
	if err != nil {
		return err
	}
	err = d.Set("created_at", snapshot.CreateAt)
	if err != nil {
		return err
	}
	err = d.Set("updated_at", snapshot.UpdatedAt)
	if err != nil {
		return err
	}
	err = d.Set("snapshot_type", snapshot.SnapshotType)
	if err != nil {
		return err
	}
	err = d.Set("type", snapshot.Type)
	if err != nil {
		return err
	}
	err = d.Set("availability_zone", snapshot.ZoneName)
	if err != nil {
		return err
	}
	err = d.Set("region_name", snapshot.RegionName)
	if err != nil {
		return err
	}
	return nil
}
