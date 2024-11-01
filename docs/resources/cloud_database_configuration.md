---
subcategory: Cloud Database
page_title: "Bizfly Cloud: bizflycloud_cloud_database_configuration"
description: |-
    Provide a Bizfly Cloud Configuration Group resource. This can be used to create, modify, and delete.
---

# Resource: bizflycloud_cloud_database_configuration

Provides a Bizfly Cloud Database Instances resource. This can be used to create, modify, and delete.

## Example

Create a configuration group

```hcl
# Create a configuration group
resource "bizflycloud_cloud_database_configuration" "terraform_appendOnly" {
  name = "terraform_appendOnly"
  datastore = {
    id         = "4b0c246a-2152-4eef-98de-d9c8e9d7550b",
    name       = "5.0.14"
    type       = "Redis",
    version_id = "8b3b46fa-1141-46ab-8d05-ce2720a0dcb2",
  }

  parameter {
    name  = "appendonly"
    value = "false"
  }

  parameter {
    name  = "maxclients"
    value = "1000"
  }
}
```

### Argument Reference

The following arguments are supported:

-   `name` - (Required) Name of this database instance
-   `datastore` - (Required) The datastore of database instance
-   `parameter` - (Required) The define for an option
    -   `name` - the option's name
    -   `value` - the option's value

### Atrributes Reference

The following attributes are exported:

-   `nodes` - The list nodes - that being apply this configuration group
