package bizflycloud

import (
	"context"
	"fmt"
	"github.com/bizflycloud/gobizfly"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

func resourceBizflyCloudScheduledVolumeBackup() *schema.Resource {
	return &schema.Resource{
		Create: resourceBizflyCloudScheduledVolumeBackupCreate,
		Read:   resourceBizflyCloudScheduledVolumeBackupRead,
		Update: resourceBizflyCloudScheduledVolumeBackupUpdate,
		Delete: resourceBizflyCloudScheduledVolumeBackupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"volume_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"scheduled_hour": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"frequency": {
				Type:     schema.TypeString,
				Required: true,
			},
			"size": {
				Type:     schema.TypeString,
				Required: true,
			},
			"next_run_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"billing_plan": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceBizflyCloudScheduledVolumeBackupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()

	log.Println("[DEBUG] Creating scheduled volume backup")
	brq := &gobizfly.CreateBackupPayload{
		ResourceID: d.Get("volume_id").(string),
		Frequency:  d.Get("frequency").(string),
		Size:       d.Get("size").(string),
		Hour:       d.Get("scheduled_hour").(int),
	}
	log.Printf("[DEBUG] Create scheduled volume backup payload: %#v\n", brq)
	backup, err := client.ScheduledVolumeBackup.Create(context.Background(), brq)
	if err != nil {
		return fmt.Errorf("[DEBUG] Error creating scheduled volume backup: %s", err)
	}
	log.Println("[DEBUG] set id " + backup.ID)
	d.SetId(backup.ID)
	return resourceBizflyCloudScheduledVolumeBackupRead(d, meta)
}

func resourceBizflyCloudScheduledVolumeBackupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	log.Printf("[DEBUG] Reading scheduled volume backup %s", d.Id())
	backup, err := client.ScheduledVolumeBackup.Get(context.Background(), d.Id())
	if err != nil {
		if err == gobizfly.ErrNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[DEBUG] Error reading scheduled volume backup: %s", err)
	}
	_ = d.Set("id", backup.ID)
	_ = d.Set("volume_id", backup.ResourceID)
	_ = d.Set("scheduled_hour", backup.ScheduledHour)
	_ = d.Set("frequency", backup.Options.Frequency)
	_ = d.Set("size", backup.Options.Size)
	_ = d.Set("next_run_at", backup.NextRunAt)
	_ = d.Set("created_at", backup.CreatedAt)
	_ = d.Set("resource_type", backup.ResourceType)
	_ = d.Set("tenant_id", backup.TenantID)
	_ = d.Set("type", backup.Type)
	_ = d.Set("billing_plan", backup.BillingPlan)
	return nil
}

func resourceBizflyCloudScheduledVolumeBackupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	err := client.ScheduledVolumeBackup.Delete(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("[DEBUG] Error deleting scheduled volume backup: %s", err)
	}
	return nil
}

func resourceBizflyCloudScheduledVolumeBackupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	log.Printf("[DEBUG] Updating scheduled volume backup %s", d.Id())
	brq := &gobizfly.UpdateBackupPayload{}
	if d.HasChange("scheduled_hour") {
		brq.Hour = d.Get("scheduled_hour").(int)
	}
	if d.HasChange("frequency") {
		brq.Frequency = d.Get("frequency").(string)
	}
	if d.HasChange("size") {
		brq.Size = d.Get("size").(string)
	}
	_, err := client.ScheduledVolumeBackup.Update(context.Background(), d.Id(), brq)
	if err != nil {
		return fmt.Errorf("[DEBUG] Error updating scheduled volume backup: %s", err)
	}
	return resourceBizflyCloudScheduledVolumeBackupRead(d, meta)
}
