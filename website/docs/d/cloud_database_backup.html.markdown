---
subcategory: Cloud Database
layout: "bizflycloud"
page_title: "Bizfly Cloud: bizflycloud_cloud_database_backup"
description: |-
  Provides a Bizfly Cloud Database Backup. This can be used to get backup detail from backup id.
---

# Data Source: bizflycloud_cloud_database_backup

Get ìnformation about Bizfly Cloud Database Backups

## Example Usage

```hcl
# Get information of Cloud Database Backup from datasource
data "bizflycloud_cloud_database_backup" "backup_daily" {
  id = "7201f1dd-6185-46a2-b488-87f2a4fce70e"
}
```

## Argument Reference

The following arguments are supported:

* `id` - The ID of this cloud database backup

## Attributes Reference

The following attributes are exported:

* `created` - The time that init this cloud database backup
* `datastore` - The datastore of backup
* `instance_id` - The database instance ID that init this cloud database backup
* `name` - The name of this cloud database backup
* `node_id` - The node ID that init this cloud database backup
* `parent_id` - The backup ID - that is parent of this cloud database backup when do backup incremental
* `size` - The size of this cloud database backup
* `status` - The status of this cloud database backup. Have some status like: `COMPLETE`, `BUILD`, `ERROR`
* `type` - The type of this cloud database backup. Have some type like: `manual`, `automatic`
* `updated` - The time that being created done (status is COMPLETE)