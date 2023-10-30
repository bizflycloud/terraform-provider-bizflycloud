---
layout: "bizflycloud"
subcategory: Cloud Database
page_title: "Bizfly Cloud: bizflycloud_cloud_database_datastore"
sidebar_current: "docs-bizflycloud-datasource-cloud-database"
description: |-
  Provides a Bizfly Cloud Database Datastore. This can be used to get database node detail from database node id.
---

# bizflycloud\_cloud\_database\_datastore

Get ìnformation about Bizfly Cloud Database Datastore

## Example Usage

```hcl
# Get information of Cloud Database Datastore from datasource
data "bizflycloud_cloud_database_datastore" "redis5014" {
  type = "Redis"
  name = "5.0.14"
}
```

## Argument Reference

The following arguments are supported:

* `type` - The type of datastore
* `name` - Name of this database node (option)

## Attributes Reference

The following attributes are exported:

* `id` - The ID of this datastore
* `version_id` - The ID of this version that have name