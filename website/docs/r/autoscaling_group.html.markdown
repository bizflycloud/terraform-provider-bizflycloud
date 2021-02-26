---
layout: "bizflycloud"
page_title: "BizFly Cloud: bizflycloud_autoscaling_group"
sidebar_current: "docs-bizflycloud-resource-autoscaling-group"
description: - Provide a BizFly Cloud AutoScaling Group resource. This can be used to create, modify, and delete.
---

# bizflycloud\_autoscaling\_group

Provides a BizFly Cloud AutoScaling Group resource. This can be used to create, modify, and delete.

## Example
```hcl
# Create a new AutoScaling Group
resource "bizflycloud_autoscaling_group" "coreAPI" {
  name                    = "coreAPI"
  launch_configuration_id = bizflycloud_autoscaling_launch_configuration.basic-centos-terrafrom.id
  max_size                = 2
  min_size                = 1
  desired_capacity        = 1
  load_balancers {
    load_balancer_id  = "f659d36b-6c0d-48da-a92c-65c21e847491"
    server_group_id   = "52370ba8-dae9-41ff-b416-e24b5579e1fe"
    server_group_port = 443
  }
}
```

### Argument Reference

The following arguments are supported:

* `name` -(Required) The name of AutoScaling Group
* `desired_capacity` - (Required) The desired capacity value of AutoScaling Group
* `launch_configuration_id` - (Required) The ID of Launch Configuration which be used to create AutoScaling Group
* `load_balancers` - (Optional) The information of load balancer was using by AutoScaling Group included:
    - `load_balancer_id`: (Required) The ID of Load Balancer
    - `server_group_id`: (Required) The ID of Server Group in Load Balancer above
    - `server_group_port`: (Required) The port number was public by Server Group above
* `max_size` - (Required) The max size value of AutoScaling Group which value was limited when do scale out
* `min_size` - (Required) The min size value of AutoScaling Group which value was limited when do scale in


### Atrributes Reference

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