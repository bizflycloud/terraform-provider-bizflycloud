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
    members {
        name = "member2"
        address = "10.20.165.30"
        protocol_port = 80
        weight = 1
    }
    members {
        name = "member2"
        address = "10.20.165.40"
        protocol_port = 80
        weight = 1
    }
    health_monitor {
        name = "hm1"
        type = "TCP"
        timeout = 100
    }
    persistent {
        type = "APP_COOKIE"
        cookie_name = "TEST"
    }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of pool
* `description` - (Optional) The description for pool
* `protocol` - (Required) The protocol for pool: `HTTP`, `TCP`, `TERMINATED_HTTPS`
* `load_balancer_id` - (Required) The ID of Load Balancer
* `algorithm` - (Required) The algorithm to balance the server in pool. Supported algorithm: `ROUND_ROBIN`, `SOURCE_IP`, `LEAST_CONNECTIONS`
* `members` - (Optional) A member block as documented below
* `health_monitor` - (Optional) A health monitor block as documented below
* `persistent` - (Optional) Setup session persistent for pool. Session Persistent block as documented below.

Members (`members`) support the following:

* `name` - (Required) Name of member
* `address` - (Required)  Address of member
* `weight` - (Optional) Weight of member. Default value is 1.
* `protocol_port` - (Required) Port for member
* `backup` - (Optional) Member is backup or not.

Health Monitor (`health_monitor`) support the following:

* `name` - (Required) Name of health monitor
* `type` - (Required) Type of health monitor. Support: `TCP`, `HTTP`
* `timeout` - (Optional) Health Check timeout. Default is 3 (second)
* `max_retries` - (Optional) Health Check max retries. Default is 3 (second).
* `max_retries_down` - (Optional) Health Check max retries down. Default is 3 (second).
* `delay` - (Optional) Delay in second before checking. Default is 3 (second)
* `http_method` - (Optional) HTTP method when using `HTTP` health check type. Default is `GET`
* `url_path` - (Optional) HTTP URL path when using `HTTP` health check type. Default is `/`
* `expected_code` - (Optional) HTTP expected codes when using `HTTP` health check type. Default is `200-409`. You can specify one status code (`200`), list of status code (`200,201`) or range of status code (`200-400`)

Session Persistent (`persistent`) support the following: 
* `type` - (Required) Type of session persistent. Supported: `SOURCE_IP`, `HTTP_COOKIE` and `APP_COOKIE`
* `cookie_name` - (Optional) The name of the cookie if persistence mode is set appropriately. Required if `type` = `APP_COOKIE`.

## Attributes Reference

The following attributes are exported:

* `name` - The name of listener
* `port` - The port for listener
* `protocol` -  The protocol for listener: `HTTP`, `TCP`, `TERMINATED_HTTPS`
* `default_pool_id`  - The default pool ID which are using for the listener
* `default_tls_ref`  - The TLS reference link for listener. The option is using when protocol is `TERMINATED_HTTPS`
* `load_balancer_id`  - The ID of Load Balancer
