---
layout: "bizflycloud"
page_title: "Bizfly Cloud: bizflycloud_volume_attachment
sidebar_current: "docs-bizflycloud-resource-volume-attachment"
description: |-
  Provides a Bizfly Cloud Volume Attachment resource. This can be used to create and delete volume attachment to server.
---

# bizflycloud\_volume\_attachment

Provides a Bizfly Cloud Volume Attachment resource. This can be used to create and delete volume attachment to server.

## Example Usage

```hcl
# Create a new volume attachment
resource "bizflycloud_volume_attachment" "volume_attachment1" {
    server_id = bizflycloud_server.tf_server1.id
    volume_id = bizflycloud_volume.volume1.id
}
```

## Argument Reference

The following arguments are supported:

* `server_id` - (Required) The ID of the server.
* `volume_id` - (Required) The ID of the volume.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Volume ID
* `server_id` - The ID of the server
* `volume_id` - The ID of the volume