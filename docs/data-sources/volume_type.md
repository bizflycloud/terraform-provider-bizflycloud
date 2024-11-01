---
subcategory: Cloud Server
page_title: "Bizfly Cloud: bizflycloud_volume_type"
description: |-
    Provides a Bizfly Cloud volume type datasource. This can be used to read volume type.
---

# Data Source: bizflycloud_volume_type

Get information about Bizfly Cloud volume type resource.

## Example Usage

```hcl
# Get information of volume type
data "bizflycloud_volume_type" "example_volume_type" {
  name = "HDD"
  category = "premium"
}
```

## Argument Reference

The following arguments are supported:

-   `name` - The name of server type
-   `category` - The category of volume type

## Attributes Reference

The following attributes are exported:

-   `name`- The name of server type
-   `category`- The category of volume type
-   `type`- The type of volume type
-   `availability_zones`- The availability zones of volume type
