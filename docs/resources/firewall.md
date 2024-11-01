---
subcategory: Cloud Server
page_title: "Bizfly Cloud: bizflycloud_firewall"
description: |-
    Provides a Bizfly Cloud Firewall resource. This can be used to create, modify, and delete firewall.
---

# Resource: bizflycloud_firewall

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
    network_interfaces = [
        "${bizflycloud_network_interface.network_interface.id}"
    ]

}
```

## Argument Reference

The following arguments are supported:

-   `name` - (Required) Name of the firewall
-   `ingress` - (Optional) Can be specified multiple times for each ingress rule. Each ingress block supports fields documented below
-   `egress` - (Optional) Can be specified multiple times for each egress rule. Each egress block supports fields documented below.
-   `network_interfaces` - (Optional) Can be specified network interfaces use the firewall.

The `ingress` and `egress` block supports:

-   `cidr` - (Required) CIDR Block: IPv4 or IPv6 CIDR. Example: 0.0.0.0/24, ::/0
-   `port_range` - (Optional) Port range. Example: `80` or `8000-9000`
-   `protocol` - (Optional) Layer 4 protocol. Available: tcp or udp

## Attributes Reference

The following attributes are exported:

-   `id` - The ID of the firewall
-   `name`- The name of the firewall
-   `ingress` - (Optional) Can be specified multiple times for each ingress rule. Each ingress block supports fields documented below
    -   `cidr` - CIDR Block: IPv4 or IPv6 CIDR
    -   `port_range` - Port range. Example: `80` or `8000-9000`
-   `egress` - (Optional) Can be specified multiple times for each egress rule. Each egress block supports fields documented below.
    -   `cidr` - CIDR Block: IPv4 or IPv6 CIDR
    -   `port_range` - Port range. Example: `80` or `8000-9000`
-   `rules_count` - Number of rules of the firewall
-   `network_interface_count` - Number of network interface of the firewall
-   `network_interfaces` - (Optional) Can be specified network interfaces use the firewall.

## Import

Bizfly Cloud firewall resource can be imported using the firewall id in the Bizfly manage dashboard

```
$ terraform import bizflycloud_firewall.fw1 firewall-id
```
