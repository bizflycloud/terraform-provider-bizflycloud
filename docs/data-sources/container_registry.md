---
page_title: "BizFly Cloud: bizflycloud_container_registry"
subcategory: "Container Registry"
description: |-
    Get information about a container registry.
---

# bizflycloud_container_registry

Use this data source to get information about an existing container registry

## Example Usage

```hcl
data "bizflycloud_container_registry" "registry" {
  name = "my-registry"
}
```

## Argument Reference

The following arguments are supported:

-   `name` - (Required) The name of the container registry

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

-   `id` - The ID of the container registry (same as name)
-   `public` - Whether the registry is public or private
-   `created_at` - The creation timestamp of the container registry
