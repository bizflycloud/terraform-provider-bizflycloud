package bizflycloud

import (
	"context"
	"fmt"

	"github.com/bizflycloud/gobizfly"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBizFlyCloudVolumeSnapshot() *schema.Resource {
	return &schema.Resource{
		Create:        resourceBizFlyCloudVolumeSnapshotCreate,
		Read:          resourceBizFlyCloudVolumeSnapshotRead,
		Delete:        resourceBizFlyCloudVolumeSnapshotDelete,
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"volume_id": {
				Type:     schema.TypeString,
				Required: true,
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
		},
	}
}

func resourceBizFlyCloudVolumeSnapshotCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	scr := gobizfly.SnapshotCreateRequest{
		Name:     d.Get("name").(string),
		VolumeId: d.Get("volume_id").(string),
		Force:    true,
	}
	snapshot, err := client.Snapshot.Create(context.Background(), &scr)
	if err != nil {
		return fmt.Errorf("Error creating snapshot: %v", err)
	}
	d.SetId(snapshot.Id)
	d.Set("volume_id", snapshot.VolumeId)
	return resourceBizFlyCloudVolumeSnapshotRead(d, meta)
}

func resourceBizFlyCloudVolumeSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	snapshot, err := client.Snapshot.Get(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving snapshot %s: %v", d.Id(), err)
	}
	_ = d.Set("name", snapshot.Name)
	_ = d.Set("size", snapshot.Size)
	_ = d.Set("status", snapshot.Status)
	_ = d.Set("volume_id", snapshot.VolumeId)
	_ = d.Set("created_at", snapshot.CreateAt)
	_ = d.Set("updated_at", snapshot.UpdatedAt)
	return nil
}

func resourceBizFlyCloudVolumeSnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	err := client.Snapshot.Delete(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error deleting volume snapshot %s: %v", d.Id(), err)
	}
	return nil
}
