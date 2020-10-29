---
layout: "bizflycloud"
page_title: "BizFly Cloud: bizflycloud_loadbalancer_listener"
sidebar_current: "docs-bizflycloud-resource-loadbalancer-listener"
description: |-
  Provides a BizFly Cloud Listener of Load Balancer resource. This can be used to create, modify, and delete listeners of Load Balancer.
---

# bizflycloud\_loadbalancer_listener

Provides a BizFly Cloud Listener of Load Balancer resource. This can be used to create,
modify, and delete listeners of Load Balancer.

## Example Create Listener for Load Balancer 

```hcl
# Create a new Listener for Load Balancer
resource "bizflycloud_loadbalancer_listener" "l1" {
    name = "bizfly-listener-tf-1"
    port = 80
    protocol = "HTTP"
    load_balancer_id = "${bizflycloud_loadbalancer.lb1.id}"
    default_pool_id = "${bizflycloud_loadbalancer_pool.pool1.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of listener
* `port` - (Required) The port for listener
* `description` - (Optional) The description for listener
* `protocol` - (Required) The protocol for listener: `HTTP`, `TCP`, `TERMINATED_HTTPS`
* `default_pool_id` - (Required) The default pool ID which are using for the listener
* `default_tls_ref` - (Optional) The TLS reference link for listener. The option is using when protocol is `TERMINATED_HTTPS`
* `load_balancer_id` - (Required) The ID of Load Balancer

## Attributes Reference

The following attributes are exported:

* `name` - The name of listener
* `port` - The port for listener
* `protocol` -  The protocol for listener: `HTTP`, `TCP`, `TERMINATED_HTTPS`
* `default_pool_id`  - The default pool ID which are using for the listener
* `default_tls_ref`  - The TLS reference link for listener. The option is using when protocol is `TERMINATED_HTTPS`
* `load_balancer_id`  - The ID of Load Balancer
