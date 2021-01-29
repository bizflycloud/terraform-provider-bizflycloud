---
layout: "bizflycloud"
page_title: "BizFly Cloud: bizflycloud_loadbalancer"
sidebar_current: "docs-bizflycloud-resource-loadbalancer"
description: |-
  Provides a BizFly Cloud Load Balancer resource. This can be used to create, modify, and delete Load balancers.
---

# bizflycloud\_loadbalancer

Provides a BizFly Cloud Load Balancer resource. This can be used to create,
modify, and delete Load Balancer.

## Example Create Load Balancer with external network facing

```hcl
# Create a new Load Balancer with external network facing
resource "bizflycloud_loadbalancer" "lb1" {
    name = "sapd-tf-lb-1"
    type = "small"
    network_type = "external"
}
```

## Example Create Load Balancer with only internal network

```hcl
# Create a new Load Balancer with only internal network
resource "bizflycloud_loadbalancer" "lb2" {
    name = "bizfly-tf-lb-2"
    type = "medium"
    network_type = "internal"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of load balancer
* `network_type` - (Optional) - The type of network: `external` or `internal`. Default value is `external`
* `type` - (Optional) The type of load balancer: `small`, `medium` or `large`. Default is `medium`
## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Load Balancer
* `name`- The name of the Load Balancer
* `network_type` - The type of network 
* `type` - The type of Load Balancer
* `vip_address` - The VIP of Load Balancer
* `provisioning_status` - The provisioning status of Load Balancer
* `operating_status` - The operating status of Load Balancer
* `pools` - The list ID of pool belong to load balancer
* `listeners` - The list ID of listener belong to load balancer
