---
subcategory: AutoScaling
page_title: "Bizfly Cloud: bizflycloud_autoscaling_scaleout_policy"
description: -|
    Provide a Bizfly Cloud AutoScaling Scale Out Policy resource. This can be used to create, modify, and delete.
---

Provides a Bizfly Cloud AutoScaling Scale Out Policy resource. This can be used to create, modify, and delete.

# Resource: bizflycloud_autoscaling_scaleout_policy

## Example

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

-   `cluster_id` - (Required) The ID of AutoScaling Group
-   `cooldown` - (Required) The time between two action continuous to do scale out/in
-   `metric_type` - (Required) The metric type of policy
-   `range_time` - (Required) The range time of policy when was reach to threshold value
-   `threshold` - (Required) The threshold value was using make decision to do scale out/in
-   `scale_size` - (Optional) The number member to do add/remove when to do scale out/in

### Atrributes Reference

The following attributes are exported:

-   `cluster_id` - The ID of AutoScaling Group
-   `cooldown` - The time between two action continuous to do scale out/in
-   `metric_type` - The metric type of policy included:
    -   `ram_used` - The metric is type percentage of ram used
    -   `cpu_used` - The metric is type percentage of cpu used
    -   `net_used` - The metric is type bandwidth of network
    -   `request_per_second` - The metric is type request per second of load balancer if it was configured
-   `range_time` - The range time of policy when was reach to threshold value
-   `threshold` - The threshold value was using make decision to do scale out/in
-   `scale_size` - The number member to do add/remove when to do scale out/in
