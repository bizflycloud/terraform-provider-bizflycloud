package bizflycloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

func resourceBizflyCloudVolumeAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceBizflyCloudVolumeAttachmentCreate,
		Read:   resourceBizflyCloudVolumeAttachmentRead,
		Delete: resourceBizflyCloudVolumeAttachmentDelete,
		Schema: resourceBizflyCloudVolumeAttachmentSchema(),
	}
}

func resourceBizflyCloudVolumeAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	volumeID := d.Get("volume_id").(string)
	serverID := d.Get("server_id").(string)
	log.Printf("[INFO] Attaching volume %s to server %s", volumeID, serverID)
	_, err := client.Volume.Attach(context.Background(), volumeID, serverID)
	if err != nil {
		log.Printf("[ERROR] Error attaching volume %s to server %s: %v", volumeID, serverID, err)
		return err
	}
	d.SetId(volumeID)
	return resourceBizflyCloudVolumeAttachmentRead(d, meta)
}

func resourceBizflyCloudVolumeAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	volumeID := d.Id()
	volume, err := client.Volume.Get(context.Background(), volumeID)
	if err != nil {
		log.Printf("[ERROR] Error reading volume %s: %v", volumeID, err)
		return err
	}
	if len(volume.Attachments) == 0 {
		log.Printf("[ERROR] Volume %s is not attached to any server", volumeID)
		return nil
	} else {
		_ = d.Set("server_id", volume.Attachments[0].ServerID)
	}
	_ = d.Set("volume_id", volumeID)
	return nil
}

func resourceBizflyCloudVolumeAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	volumeID := d.Id()
	serverID := d.Get("server_id").(string)
	_, err := client.Volume.Detach(context.Background(), serverID, volumeID)
	if err != nil {
		log.Printf("[ERROR] Error detaching volume %s from server %s: %v", volumeID, serverID, err)
		return err
	}
	return nil
}
