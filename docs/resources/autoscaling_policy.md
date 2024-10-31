---
subcategory: AutoScaling
layout: "bizflycloud"
page_title: "Bizfly Cloud: bizflycloud_autoscaling_policy"
description: - Provide a Bizfly Cloud AutoScaling Policy resource. This can be used to create, modify, and delete.
---

Provides a Bizfly Cloud AutoScaling Policy resource. This can be used to create, modify, and delete.

# Resource: bizflycloud_autoscaling_scalein_policy

## Example
```hcl
# Create a new AutoScaling ScaleIn Policy
resource "bizflycloud_autoscaling_scalein_policy" "name" {
  cluster_id  = bizflycloud_autoscaling_group.hutao.id
  metric_type = "ram_used"
  threshold   = 10
  range_time  = 600
  cooldown    = 600
}

```

# bizflycloud_autoscaling_scaleout_policy

```hcl
### Create a new AutoScaling Scaleout Policy
resource "bizflycloud_autoscaling_scaleout_policy" "name" {
  cluster_id  = bizflycloud_autoscaling_group.hutao.id
  metric_type = "ram_used"
  threshold   = 90
  range_time  = 600
  cooldown    = 600
}
```

### Argument Reference

The following arguments are supported:
* `cluster_id` - (Required) The ID of AutoScaling Group
* `cooldown` - (Required) The time between two action continuous to do scale out/in
* `metric_type` - (Required) The metric type of policy
* `range_time` - (Required) The range time of policy when was reach to threshold value
* `threshold` - (Required) The threshold value was using make decision to do scale out/in
* `scale_size` - (Optional) The number member to do add/remove when to do scale out/in

### Atrributes Reference

The following attributes are exported:

* `cluster_id` - The ID of AutoScaling Group
* `cooldown` - The time between two action continuous to do scale out/in
* `metric_type` - The metric type of policy included:
    - `ram_used` - The metric is type percentage of ram used
    - `cpu_used` - The metric is type percentage of cpu used
    - `net_used` - The metric is type bandwidth of network
    - `request_per_second` - The metric is type request per second of load balancer if it was configured
* `range_time` - The range time of policy when was reach to threshold value
* `threshold` - The threshold value was using make decision to do scale out/in
* `scale_size` - The number member to do add/remove when to do scale out/in



# bizflycloud_autoscaling_deletion_policy

```hcl
# Update criteria in deletion policy
resource "bizflycloud_autoscaling_deletion_policy" "deletion_policy" {
  cluster_id = bizflycloud_autoscaling_group.hutao.id
  criteria   = "YOUNGEST_FIRST"
}
```

### Argument Reference

The following arguments are supported:
* `cluster_id` - (Required) The ID of AutoScaling Group
* `criteria` - (Required) The criteria used in selecting node candidates for deletion. Included:
    - `OLDEST_FIRST` - always select node(s) which were created earlier than other nodes.
    - `YOUNGEST_FIRST` - always select node(s) which were created recently instead of those created earlier.
    - `OLDEST_PROFILE_FIRST` - compare the profile used by each individual nodes and select the node(s) whose profile(s) were created earlier than others.
    - `RANDOM` - randomly select node(s) from the cluster for deletion. This is the default criteria if omitted.
