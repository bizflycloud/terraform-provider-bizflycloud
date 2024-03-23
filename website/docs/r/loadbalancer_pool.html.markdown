---
layout: "bizflycloud"
page_title: "Bizfly Cloud: bizflycloud_loadbalancer_pool"
sidebar_current: "docs-bizflycloud-resource-loadbalancer-pool"
description: |-
  Provides a Bizfly Cloud Pool of Load Balancer resource. This can be used to create, modify, and delete pools of Load Balancer.
---

# bizflycloud\_loadbalancer_pool

Provides a Bizfly Cloud Pool of Load Balancer resource. This can be used to create,
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
        max_retries = 3
        max_retries_down = 3
        delay = 5
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
* `protocol` - (Required) The protocol for pool: `HTTP`, `TCP`, `PROXY`, `UDP`
* `load_balancer_id` - (Required) The ID of Load Balancer
* `algorithm` - (Required) The algorithm to balance the server in pool. Supported algorithm: `ROUND_ROBIN`, `SOURCE_IP`, `LEAST_CONNECTIONS`
* `members` - (Optional) A member block as documented below
  - `name` - (Required) Name of member
  - `address` - (Required)  Address of member
  - `weight` - (Optional) Weight of member [1-256]. Default value is 1.
  - `protocol_port` - (Required) Port for member
  - `backup` - (Optional) Member is backup or not (Default is `false`): `true`, `false`
* `health_monitor` - (Optional) A health monitor block as documented below
  - `name` - (Optional) Name of health monitor (Default is `pool-monitor`)
  - `type` - (Required) Type of health monitor. Support: `HTTP`, `HTTPS`, `PING`, `SCTP`, `TCP`, `TLS-HELLO`, `UDP-CONNECT`
  - `timeout` - (Optional) Health Check timeout. Default is 5 (second)
  - `max_retries` - (Optional) Health Check max retries [1-10]. Default is 3 (second).
  - `max_retries_down` - (Optional) Health Check max retries down [1-10]. Default is 3 (second).
  - `delay` - (Optional) Delay in second before checking. Default is 5 (second)
  - `http_method` - (Optional) HTTP method when using `HTTP` health check type: `GET`, `POST`, `HEAD`, `PUT`, `DELETE`, `TRACE`, `OPTIONS`, `PATCH`, `CONNECT` 
  - `url_path` - (Optional) HTTP URL path when using `HTTP` health check type (Valid start with `/`).
  - `expected_code` - (Optional) HTTP expected codes when using `HTTP` health check type: `200`, `201`, `202`, `203`, `204`

* `persistent` - (Optional) Setup session persistent for pool. Session Persistent block as documented below.
  - `type` - (Required) Type of session persistent. Supported: `SOURCE_IP`, `HTTP_COOKIE` and `APP_COOKIE`
  - `cookie_name` - (Optional) The name of the cookie if persistence mode is set appropriately. Required if `type` = `APP_COOKIE`.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of pool
* `name` - The name of pool
* `algorithm` - The algorithm of pool
* `description` -  The description of pool
* `protocol`  - The protocol of pool
* `load_balancer_id`  - The ID of loadbalancer
* `members`  - The members of pool
  - `id` - The ID of member
  - `name` - The name of member
  - `weight` - The weight of member
  - `address` - The address of member
  - `protocol_port` - The protocol port of member
  - `backup` - The backup of member
  - `operating_status` - The operating status of member
  - `provisioning_status` - The provisioning status of member
  - `subnet_id` - The subnet id of member
  - `project_id` - The project id of member
  - `created_at` - The created time of member
  - `updated_at` - The updated time of member
* `health_monitor` - The health monitor of pool
  - `id` - The ID of health monitor
  - `name` - The name of health monitor
  - `type` - The type of health monitor
  - `timeout` - The timeout of health monitor
  - `max_retries` - The max retries of health monitor
  - `max_retries_down` - The max retries down of health monitor
  - `delay` - The delay of health monitor
  - `http_method` - The http method of health monitor
  - `url_path` - The url path of health monitor
  - `expected_code` - The expected code of health monitor
  - `operating_status` - The operating status of health monitor
  - `provisioning_status` - The provisioning status of health monitor
  - `created_at` - The created time of health monitor
  - `updated_at` - The updated time of health monitor
* `persistent` - The sticky session of pool
  - `type` - The type of sticky session
  - `cookie_name` - The cookie name when `type` is `APP_COOKIE`


## Import

Bizfly Cloud loadbalancer pool resource can be imported using the pool id in the Bizfly manage dashboard

```
$ terraform import bizflycloud_loadbalancer_pool.pool1 pool-id
```