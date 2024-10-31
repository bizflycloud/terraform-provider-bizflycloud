---
subcategory: Cloud Load Balancer
layout: "bizflycloud"
page_title: "Bizfly Cloud: bizflycloud_loadbalancer_listener"
description: |-
  Provides a Bizfly Cloud Listener of Load Balancer resource. This can be used to create, modify, and delete listeners of Load Balancer.
---

# Resource: bizflycloud_loadbalancer_listener

Provides a Bizfly Cloud Listener of Load Balancer resource. This can be used to create,
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
    listener_timeout = 5000
    server_timeout = 5000
    server_connect_timeout = 5000
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of listener
* `port` - (Required) The port for listener
* `description` - (Optional) The description for listener
* `protocol` - (Required) The protocol for listener: `HTTP`, `TCP`, `TERMINATED_HTTPS`, `UDP`
* `default_pool_id` - (Required) The default pool ID which are using for the listener
* `default_tls_ref` - (Optional) The TLS reference link for listener. The option is using when protocol is `TERMINATED_HTTPS`
* `load_balancer_id` - (Required) The ID of Load Balancer
* `listener_timeout` - (Optional) The listener timeout (Default: 5000)
* `server_timeout` - (Optional) The server timeout (Default: 5000)
* `server_connect_timeout` - (Optional) The server connect timeout (Default: 5000)

## Attributes Reference

The following attributes are exported:

* `name` - The name of listener
* `port` - The port for listener
* `description` - The description for listener
* `protocol` -  The protocol for listener: `HTTP`, `TCP`, `TERMINATED_HTTPS`, `UDP`
* `default_pool_id`  - The default pool ID which are using for the listener
* `default_tls_ref`  - The TLS reference link for listener. The option is using when protocol is `TERMINATED_HTTPS`
* `load_balancer_id`  - The ID of Load Balancer
* `listener_timeout` - The listener timeout 
* `server_timeout` - The server timeout
* `server_connect_timeout` - The server connect timeout
* `operating_status` - The operating status
* `provisioning_status` - The provisioning status
* `l7policy_ids` - The L7 policy IDs
* `created_at` - The created time of listener
* `updated_at` - The updated time of listener

## Import

Bizfly Cloud loadbalancer listener resource can be imported using the listener id in the Bizfly manage dashboard

```
$ terraform import bizflycloud_loadbalancer_listener.l1 listener-id
```