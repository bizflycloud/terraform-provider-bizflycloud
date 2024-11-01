---
subcategory: Cloud Database
page_title: "Bizfly Cloud: bizflycloud_cloud_database_backup"
description: -|
    Provide a Bizfly Cloud Database Backup resource. This can be used to create, modify, and delete.
---

# Resource: bizflycloud_cloud_database_backup

Provides a Bizfly Cloud Database Backup resource. This can be used to create, modify, and delete.

## Example

Create a backup for an instance with instance_id

```hcl
# Create a new Cloud Database Backup with instance_id
resource "bizflycloud_cloud_database_backup" "terraform_backup" {
  name        = "terraform_backup"
  node_id     = ""
  instance_id = "a30ae0fd-aa80-4540-851e-f0a4e2de3b62"
}
```

Create a backup for a node with node_id

```hcl
# Create a new Cloud Database Backup with node_id
resource "bizflycloud_cloud_database_backup" "terraform_backup" {
  name        = "terraform_backup"
  instance_id     = ""
  node_id = "a30ae0fd-aa80-4540-851e-f0a4e2de3b62"
}
```

### Argument Reference

The following arguments are supported:

-   `name` - (Required) The name of Cloud Database Backup
-   `node_id` - (Required) The ID of cloud database node want do backup - set "" if want ignore
-   `instance_id` - (Required) The ID of cloud database instance want do backup - set "" if want ignore. If you set instance_id, the order of priority do backup will happen in secondary node, if not have secondary node, it being do in replica node, if not have replica node, it being do in primary node. This option being conflict with node_id

### Atrributes Reference

The following attributes are exported:

-   `created` - The time that init this cloud database backup
-   `datastore` - The datastore of backup
-   `id` - The ID of this cloud database backup
-   `instance_id` - The database instance ID that init this cloud database backup
-   `name` -The name of this cloud database backup
-   `node_id` - The node ID that init this cloud database backup
-   `parent_id` - The backup ID - that is parent of this cloud database backup when do backup incremental
-   `size` - The size of this cloud database backup
-   `status` - The status of this cloud database backup. Have some status like: `COMPLETE`, `BUILD`, `ERROR`
-   `type` - The type of this cloud database backup. Have some type like: `manual`, `automatic`
-   `updated` - The time that being created done (status is COMPLETE)
