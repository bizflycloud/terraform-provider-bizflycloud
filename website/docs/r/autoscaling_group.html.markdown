---
layout: "bizflycloud"
page_title: "Bizfly Cloud: bizflycloud_autoscaling_group"
sidebar_current: "docs-bizflycloud-resource-autoscaling-group"
description: - Provide a Bizfly Cloud AutoScaling Group resource. This can be used to create, modify, and delete.
---

# bizflycloud\_autoscaling\_group

Provides a Bizfly Cloud AutoScaling Group resource. This can be used to create, modify, and delete.

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
