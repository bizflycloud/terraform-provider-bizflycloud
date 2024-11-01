---
subcategory: Cloud Server
page_title: "Bizfly Cloud: bizflycloud_server_type"
description: |-
    Provides a Bizfly Cloud Server type datasource. This can be used to read server type.
---

# Data Source: bizflycloud_server_type

Get information about Bizfly Cloud server type resource.

## Example Usage

```hcl
# Get information of server type
data "bizflycloud_server_type" "example_server_type" {
  name = "dedicated"
}
```

## Argument Reference

The following arguments are supported:

-   `name` - The name of server type

## Attributes Reference

The following attributes are exported:

-   `id`- The ID of server types
-   `name`- The name of server type
-   `enabled` - The state of server type
-   `compute_class` - The compute class of server type
-   `priority` - The priority of server type
