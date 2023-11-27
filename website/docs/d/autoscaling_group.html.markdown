---
layout: "bizflycloud"
page_title: "Bizfly Cloud: bizflycloud_autoscaling_group"
sidebar_current: "docs-bizflycloud-datasource-autoscaling-group"
description: |-
  Provides a Bizfly Cloud AutoScaling Group
---

# bizflycloud\_autoscaling\_group

Get information about Bizfly Cloud AutoScaling Group.

## Example Usage

```hcl
# Get information of an autoscaling group
data "bizflycloud_autoscaling_group" "coreAPI" {
  id   = "d819844a-e200-47e5-b32b-f49663406cd0"
  name = "coreAPI"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The ID of AutoScaling Group

## Attributes Reference

The following attributes are exported:

* `id` - The ID of AutoScaling Group
* `name` -The name of AutoScaling Group
* `desired_capacity` - The desired capacity value of AutoScaling Group
* `launch_configuration_id` - The ID of Launch Configuration which be used to create AutoScaling Group
* `launch_configuration_name` - The name of Launch Configuration which be used to create AutoScaling Group
* `load_balancers` - The information of load balancer was using by AutoScaling Group included:
    - `load_balancer_id`: The ID of Load Balancer
    - `server_group_id`: The ID of Server Group in Load Balancer above
    - `server_group_port`: The port number was public by Server Group above
* `max_size` - The max size value of AutoScaling Group which value was limited when do scale out
* `min_size` - The min size value of AutoScaling Group which value was limited when do scale in
* `node_ids` - List ID of members in AutoScaling Group
* `status` - The status of AutoScaling Group
* `scale_in_info` - List of AutoScaling Policy to do scale in included:
    - `cooldown` - The time between two action continuous to do scale in
    - `metric_type` - The metric type of policy included:
        - `ram_used` - The metric is type percentage of ram used
        - `cpu_used` - The metric is type percentage of cpu used
        - `net_used` - The metric is type bandwidth of network
        - `request_per_second` - The metric is type request per second of load balancer if it was configured
    - `range_time` - The range time of policy when was reach to threshold value
    - `threshold` - The threshold value was using make decision to do scale in
    - `scale_size` - The number member to do remove when to do scale in
* `scale_out_info` - List of AutoScaling Policy to do scale out
    - `cooldown` - The time between two action continuous to do scale out
    - `metric_type` - The metric type of policy included:
        - `ram_used` - The metric is type percentage of ram used
        - `cpu_used` - The metric is type percentage of cpu used
        - `net_used` - The metric is type bandwidth of network
        - `request_per_second` - The metric is type request per second of load balancer if it was configured
    - `range_time` - The range time of policy when was reach to threshold value
    - `threshold` - The threshold value was using make decision to do scale out
    - `scale_size` - The number member to do remove when to do scale out