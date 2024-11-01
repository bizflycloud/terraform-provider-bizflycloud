---
subcategory: AutoScaling
page_title: "Bizfly Cloud: bizflycloud_autoscaling_deletion_policy"
description: |-
    Provide a Bizfly Cloud AutoScaling Deletion Policy resource. This can be used to create, modify, and delete.
---

Provides a Bizfly Cloud AutoScaling Deletion Policy resource. This can be used to create, modify, and delete.

# Resource: bizflycloud_autoscaling_deletion_policy

## Example

```hcl
# Update criteria in deletion policy
resource "bizflycloud_autoscaling_deletion_policy" "deletion_policy" {
  cluster_id = bizflycloud_autoscaling_group.hutao.id
  criteria   = "YOUNGEST_FIRST"
}
```

### Argument Reference

The following arguments are supported:

-   `cluster_id` - (Required) The ID of AutoScaling Group
-   `criteria` - (Required) The criteria used in selecting node candidates for deletion. Included:
    -   `OLDEST_FIRST` - always select node(s) which were created earlier than other nodes.
    -   `YOUNGEST_FIRST` - always select node(s) which were created recently instead of those created earlier.
    -   `OLDEST_PROFILE_FIRST` - compare the profile used by each individual nodes and select the node(s) whose profile(s) were created earlier than others.
    -   `RANDOM` - randomly select node(s) from the cluster for deletion. This is the default criteria if omitted.
