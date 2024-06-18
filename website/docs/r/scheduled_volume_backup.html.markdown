---
subcategory: Cloud Server
layout: "bizflycloud"
page_title: "Bizfly Cloud: bizflycloud_scheduled_volume_backup"
description: |-
Provides a Bizfly Cloud Scheduled Volume Backup resource. This can be used to create and delete, update scheduled volume backup
---

# Resource: bizflycloud_scheduled_volume_backup

Provides a Bizfly Cloud schedule volume backup resource. This can be used to create,
and delete ssh key.
## Example Usage

```hcl
resource "bizflycloud_scheduled_volume_backup" "backup_test" {
  volume_id = "11a2e71b-8701-47a0-b247-41843db17e54"
  frequency = "2880"
  size = "4"
  scheduled_hour = 4
}
```
        

## Argument Reference

The following arguments are supported:

* `volume_id` - (Required) The volume Id targeted for backup.
* `frequency` - (Required) The interval between backups in seconds.
* `size` - (Required) The number of snapshots to keep.
* `scheduled_hour` - (Optional) The hour of the day to start the backup. Default is 0


## Attributes Reference

The following attributes are exported:

* `volume_id` - The volume Id targeted for backup.
* `frequency` - The interval between backups in seconds.
* `size` - The number of snapshots to keep.
* `scheduled_hour` - The hour of the day to start the backup
* `next_run_at` - The next time the backup will run
* `created_at` - The time the backup was created
* `updated_at` - The last time the backup was updated
* `resource_type` - The type of the resource
* `tenant_id` - The tenant id
* `type` - The type of the resource
* `billing_plan` - The billing plan

## Import

Bizfly Cloud volume backup resource can be imported using the backup id in the Bizfly manage dashboard

```
$ terraform import bizflycloud_scheduled_volume_backup.backup_test backup-id
```