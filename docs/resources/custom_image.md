---
subcategory: Cloud Server
layout: "bizflycloud"
page_title: "Bizfly Cloud: bizflycloud_custom_image"
description: |-
  Provides a Bizfly Cloud Custom Image resource. This can be used to create, and delete custom image.
---

# Resource: bizflycloud_custom_image

Provides a Bizfly Cloud Custom Image resource. This can be used to create,
and delete custom image.

## Example Usage

```hcl
# Create a new custom image
resource "bizflycloud_custom_image" "new_custom_image1" {
  name        = "new_custom_image_10"
  disk_format = "qcow2"
  image_url   = "https://releases.ubuntu.com/22.04.3/ubuntu-22.04.3-desktop-amd64.iso"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of Custom Image
* `disk_format` - (Required) The disk format of Custom Image
* `image_url` - (Required) The URL of image

## Attributes Reference

The following attributes are exported:

* `id` - The ID of Custom Image
* `name` - The name of Custom Image
* `size` - The size of Custom Image
* `disk_format` - The disk format of Custom Image
* `container_format` - The container format of Custom Image
* `billing_plan` - The billing plan of Custom Image
* `visibility` - The visibility of Custom Image
* `created_at` - The time when Custom Image was created
* `updated_at` - The time when Custom Image was updated
* `description` - The description of Custom Image


## Import

Bizfly Cloud custom image resource can be imported using the custom image id in the Bizfly manage dashboard

```
$ terraform import bizflycloud_custom_image.new_custom_image1 custom-image-id
```