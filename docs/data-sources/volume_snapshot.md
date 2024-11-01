---
subcategory: Cloud Server
page_title: "Bizfly Cloud: bizflycloud_volume_snapshot"
description: |-
    Provides a Bizfly Cloud Volume Snapshot
---

# Data Source: bizflycloud_volume_snapshot

Get information about Bizfly Cloud Volume Snapshot.

## Example Usage

```hcl
# Get information of a volume snapshot
data "bizflycloud_volume_snapshot" "volume_snapshot" {
  id = "90ef825f-3b8e-4194-93d4-8764c02f0d66"
}
```

## Argument Reference

The following arguments are supported:

-   `id` - (Required) The ID of Volume Snapshot

## Attributes Reference

The following attributes are exported:

-   `id` - The ID of Volume Snapshot
-   `name` - The name of Volume Snapshot
-   `volume_id` - The ID of Volume which was created Volume Snapshot
-   `size` - The size of Volume Snapshot
-   `created_at` - The time when Volume Snapshot was created
-   `updated_at` - The time when Volume Snapshot was updated
-   `snapshot_type` - The type of Volume Snapshot
-   `type` - The type of Volume Snapshot
-   `availability_zone` - The name of Zone which was created Volume Snapshot
-   `region_name` - The name of Region which was created Volume Snapshot
