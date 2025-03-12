---
page_title: "Bizfly Cloud: bizflycloud_container_registry"
subcategory: "Container Registry"
description: |-
    Provides a Bizfly Cloud Container Registry resource. This can be used to create, read, and delete container registries.
---

# bizflycloud_container_registry

Provides a Bizfly Cloud Container Registry resource. This can be used to create and manage container registries.

## Example Usage

```hcl
# Create a new private container registry
resource "bizflycloud_container_registry" "registry1" {
  name = "my-registry"
}

# Create a new public container registry
resource "bizflycloud_container_registry" "registry2" {
  name   = "my-public-registry"
  public = true
}
```

## Argument Reference

The following arguments are supported:

-   `name` - (Required, ForceNew) The name of the container registry. Changing this will create a new registry.
-   `public` - (Optional) Whether the registry is public or private. Default is `false` (private). This can be updated without recreating the registry.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

-   `id` - The ID of the container registry (same as name)
-   `created_at` - The creation timestamp of the container registry

## Import

Container Registry can be imported using the `name`, e.g.

```bash
$ terraform import bizflycloud_container_registry.registry1 my-registry
```
