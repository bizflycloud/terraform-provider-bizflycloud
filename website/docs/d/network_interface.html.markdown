---
layout: "bizflycloud"
page_title: "Bizfly Cloud: bizflycloud_network_interface"
sidebar_current: "docs-bizflycloud-datasource-network-interface"
description: |-
  Provides a Bizfly Cloud Network Interface datasource. This can be used to create, modify, and delete Network Interface.
---

# bizflycloud\_network\_interface

Get ìnformation about Bizfly Cloud Network Interface resource.

## Example Usage

```hcl
# Get information of Network Interface
data "bizflycloud_network_interface" "network_interface" {
  network_id = bizflycloud_network_interface.network_interface.network_id
}
```

## Argument Reference

The following arguments are supported:

* `network_id` - The Network ID of network interface.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of network interface.
* `name`- The name of network interface.
* `network_id` - The Network ID of network interface.
* `attached_server` - The attached server of network interface.
* `fixed_ip` - The fixed IP of network interface.
* `mac_address` - The media access control address of network interface.
* `admin_state_up` - The admin state up the network interface: true or false.
* `status` - The status of network interface.
* `device_id` - The device ID of network interface.
* `port_security_enabled` - The port security enabled of network interface.
* `fixed_ips` - The fixed ips of network interface.
  * `subnet_id` - The subnet ID of network interface.
  * `ip_address` - The IP address of network interface.
* `security_groups` - List ID of security groups.
* `create_at` - The created time of network interface.
* `updated_at` - The updated time of network interface.
