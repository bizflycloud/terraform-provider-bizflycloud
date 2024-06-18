---
subcategory: Cloud Server
layout: "bizflycloud"
page_title: "Bizfly Cloud: bizflycloud_volume"
description: |-
  Provides a Bizfly Cloud Volume resource. This can be used to create, modify, and delete volumes.
---

# Resource: bizflycloud_volume

Provides a Bizfly Cloud Volume resource. This can be used to create,
modify, and delete volume.
## Example Usage

```hcl
# Create a new volume
resource "bizflycloud_volume" "volume1" {
    name = "volume1"
    size = 20
    type = "HDD"
    category = "premium"
    availability_zone = "HN1"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the volume.
* `size` - (Required) The size of the volume.
* `type` - (Required) The type of the volume: HDD or SSD.
* `category` - (Required) - The category of the volume: basic, premium, enterprise or dedicated.
* `availability_zone` - (Required) - The availability zone of the volume.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Volume
* `name`- The name of the volume
* `category` - The category of the volume
* `status` - The status of the volume
* `availability_zone` - The availability zone of volume
* `type` - The volume type[resource_bizflycloud_network_interface.go](..%2F..%2F..%2Fbizflycloud%2Fresource_bizflycloud_network_interface.go)
* `size` - The size of volume


## Import

Bizfly Cloud volume resource can be imported using the volume id in the Bizfly manage dashboard

```
$ terraform import bizflycloud_volume.volume1 volume-id
```