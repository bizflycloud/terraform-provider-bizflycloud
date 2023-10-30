---
layout: "bizflycloud"
subcategory: Cloud Database
page_title: "Bizfly Cloud: bizflycloud_cloud_database_instance"
sidebar_current: "docs-bizflycloud-datasource-cloud-database"
description: |-
  Provides a Bizfly Cloud Database Instances. This can be used to get database instance detail from database instance id.
---

# bizflycloud\_cloud\_database\_instance

Get ìnformation about Bizfly Cloud Database Instance

## Example Usage

```hcl
# Get information of Cloud Database Instance from datasource
data "bizflycloud_cloud_database_instance" "mongo604" {
  id = "7201f1dd-6185-46a2-b488-87f2a4fce70e"
}
```

## Argument Reference

The following arguments are supported:

* `id` - ID of this database instance
* `name` - Name of this database instance (option)

## Attributes Reference

The following attributes are exported:

* `autoscaling` - Information about autoscaling volume, ... for this database instance
* `created_at` - The time that init this database instance
* `datastore` - The datastore of database instance
* `dns` - The DNS of this database instance (This current just have with MongoDB)
* `instance_type` - The instance type of this database instance
* `nodes` - The list nodes - that are member of this database instance
* `public_access` - Whether database instance can be connect from public internet
* `status` - The status of this database instance. Have some status like: `ACTIVE`, `RESIZE`, `ERROR`, `BUILD`
* `volume_size` - The size of storage provisioning for this database instance, unit is `Gigabytes`