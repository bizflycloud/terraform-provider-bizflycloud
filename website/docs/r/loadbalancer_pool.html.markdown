---
layout: "bizflycloud"
page_title: "BizFly Cloud: bizflycloud_loadbalancer_pool"
sidebar_current: "docs-bizflycloud-resource-loadbalancer-pool"
description: |-
  Provides a BizFly Cloud Pool of Load Balancer resource. This can be used to create, modify, and delete pools of Load Balancer.
---

# bizflycloud\_loadbalancer_pool

Provides a BizFly Cloud Pool of Load Balancer resource. This can be used to create,
modify, and delete pools of Load Balancer.

## Example Create Pool for Load Balancer 

```hcl
# Create a new Pool for Load Balancer
resource "bizflycloud_loadbalancer_pool" "pool1" {
    name = "sapd-pool-tf-1"
    protocol = "HTTP"
    algorithm = "ROUND_ROBIN"
    load_balancer_id = "${bizflycloud_loadbalancer.lb1.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of pool
* `description` - (Optional) The description for pool
* `protocol` - (Required) The protocol for pool: `HTTP`, `TCP`, `TERMINATED_HTTPS`
* `load_balancer_id` - (Required) The ID of Load Balancer
* `algorithm` - (Required) The algorithm to balance the server in pool. Supported algorithm: `ROUND_ROBIN`, `SOURCE_IP`, `LEAST_CONNECTIONS`

## Attributes Reference

The following attributes are exported:

* `name` - The name of listener
* `port` - The port for listener
* `protocol` -  The protocol for listener: `HTTP`, `TCP`, `TERMINATED_HTTPS`
* `default_pool_id`  - The default pool ID which are using for the listener
* `default_tls_ref`  - The TLS reference link for listener. The option is using when protocol is `TERMINATED_HTTPS`
* `load_balancer_id`  - The ID of Load Balancer
