---
layout: "bizflycloud"
page_title: "Bizfly Cloud: bizflycloud_network_interface"
sidebar_current: "docs-bizflycloud-resource-network-interface"
description: |-
  Provides a Bizfly Cloud Network Interface resource. This can be used to create, modify, and delete Network Interface.
---

# bizflycloud\_network\_interface

Provides a Bizfly Cloud Network Interface resource. This can be used to create,
modify, and delete Network Interface.

## Example Usage

```hcl
# Create a new Network Interface
resource "bizflycloud_network_interface" "network_interface" {
  name = "test-name"
  network_id = "${bizflycloud_vpc_network.vpc_network.id}"
  fixed_ip = "10.108.16.5"
  firewall_ids = ["4b41c931-bf3d-443f-b311-df3817a3fbc0"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of network interface.
* `network_id` - (Required) The Network ID of network interface.
* `fixed_ip` - (Optional) The fixed IP of network interface.
* `firewall_ids` - (Optional) The list ID of security groups.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of network interface.
* `name`- The name of network interface.
* `network_id` - The Network ID of network interface.
* `fixed_ip` - The fixed IP of network interface.
* `mac_address` - The media access control address of network interface.
* `admin_state_up` - The admin state up the network interface: true or false.
* `status` - The status of network interface.
* `server_id` - The server ID of network interface.
* `port_security_enabled` - The port security enabled of network interface.
* `fixed_ips` - The fixed ips of network interface.
  * `subnet_id` - The subnet ID of network interface.
  * `ip_address` - The IP address of network interface.
* `firewall_ids` - List ID of security groups.
* `created_at` - The created time of network interface.
* `updated_at` - The updated time of network interface.


## Import

Bizfly Cloud network interface resource can be imported using the network interface id in the Bizfly manage dashboard

```
$ terraform import bizflycloud_network_interface.network_interface network-interface-id
```