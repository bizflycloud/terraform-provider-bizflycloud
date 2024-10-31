---
subcategory: Cloud Database
layout: "bizflycloud"
page_title: "Bizfly Cloud: bizflycloud_cloud_database_backup_schedule"
description: - Provide a Bizfly Cloud Database Backup Schedule resource. This can be used to create, modify, and delete.
---

# Resource: bizflycloud_cloud_database_backup_schedule

Provides a Bizfly Cloud Database Backup Schedule resource. This can be used to create, modify, and delete.

## Example

Create a backup schedule for a node with node_id

```hcl
# Create a backup schedule
resource "bizflycloud_cloud_database_backup_schedule" "terraform_backup_schedule" {
  limit_backup    = 1
  name            = "terraform_backup_schedule"
  node_id         = "edc43593-003f-475c-9179-8e4c9799ca03"
  cron_expression = "45 * * * *"
}
```


### Argument Reference

The following arguments are supported:

* `name` - (Required) The name of this cloud database backup schedule
* `node_id` - (Required) The ID of cloud database node want do backup
* `limit_backup` - (Required) The number of backup that being keep when create the next backup
* `cron_expression` - (Required) The cron pattern describe about time to do create the backup

### Atrributes Reference

Haven't any attributes are exported.

