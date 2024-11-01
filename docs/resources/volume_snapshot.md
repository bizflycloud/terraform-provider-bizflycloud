---
subcategory: Cloud Server
page_title: "Bizfly Cloud: bizflycloud_volume_snapshot"
description: |-
    Provides a Bizfly Cloud Volume Snapshot resource. This can be used to create and delete volume snapshot.
---

# Resource: bizflycloud_volume_snapshot

Provides a Bizfly Cloud Volume Snapshot resource. This can be used to create,
and delete volume snapshot.

## Example Usage

```hcl
# Create a new snapshot
resource "bizflycloud_volume_snapshot" "snapshot1" {
    name = "snapshot1"
    volume_id = "cc3cce08-f514-4186-8e19-dbc38fb8f6d0"
}
```

## Argument Reference

The following arguments are supported:

-   `name` - (Required) The name of the volume snapshot.
-   `volume_id` - (Required) The ID of volume will be take snapshot.

## Attributes Reference

The following attributes are exported:

-   `id` - The ID of the volume snapshot
-   `name`- The name of the volume snapshot
-   `size` - The size of volume volume

## Import

Bizfly Cloud volume snapshot resource can be imported using the snapshot id in the Bizfly manage dashboard

```
$ terraform import bizflycloud_volume_snapshot.snapshot1 snapshot-id
```
