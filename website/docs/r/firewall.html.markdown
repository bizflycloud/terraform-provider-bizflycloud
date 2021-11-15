---
layout: "bizflycloud"
page_title: "Bizfly Cloud: bizflycloud_firewall"
sidebar_current: "docs-bizflycloud-resource-firewall"
description: |-
  Provides a Bizfly Cloud Firewall resource. This can be used to create, modify, and delete firewall.
---

# bizflycloud\_firewall

Provides a Bizfly Cloud Firewal resource. This can be used to create,
modify, and delete firewall.

## Example Usage

```hcl
# Create a new firewall and attach to a server
resource "bizflycloud_firewall" "fw1" {
    name = "sapd-firewall-2"
    ingress {
        cidr = "0.0.0.0/0"
        port_range = "80"
        protocol = "tcp"
    }
    ingress {
        cidr = "192.168.1.0/24"
        port_range = "8000"
        protocol = "udp"
    }
    egress {
        cidr = "0.0.0.0/0"
        port_range = "80"
        protocol = "tcp"
    }
    target_server_ids = [
        "${bizflycloud_server.sapd-ubuntu-20.id}"
    ]
    network_interfaces = [
        "${bizflycloud_network_interface.network_interface.id}"
    ]

}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the firewall
* `target_server_ids` - (Optional) - List ID of the server which will be applied the firewall
* `network_interfaces` - (Optional) - List ID of the network interface which will be applied the firewall
* `ingress` - (Optional) Can be specified multiple times for each ingress rule. Each ingress block supports fields documented below
* `egress` - (Optional) Can be specified multiple times for each egress rule. Each egress block supports fields documented below.

The `ingress` and `egress` block supprts:

* `cidr` - (Required) CIDR Block: IPv4 or IPv6 CIDR
* `port_range` - (Required) Port range. Example: `80` or `8000-9000`
* `protocol` - (Required) Layer 4 protocol.  Available: tcp or udp

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the firewall
* `name`- The name of the firewall
* `target_server_ids` - List ID of the server which will be applied the firewall
* `servers_count` - Number of server are applied the firewall
* `rules_count` - Number of rules of the firewall
* `network_interface_count` - Number of network interface of the firewall
