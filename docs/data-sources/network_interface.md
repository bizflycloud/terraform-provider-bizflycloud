---
subcategory: Cloud Server
page_title: "Bizfly Cloud: bizflycloud_network_interface"
description: |-
    Provides a Bizfly Cloud Network Interface datasource. This can be used to create, modify, and delete Network Interface.
---

# Data Source: bizflycloud_network_interface

Get ìnformation about Bizfly Cloud Network Interface resource.

## Example Usage

```hcl
# Get information of Network Interface
data "bizflycloud_network_interface" "lan_ip_2" {
  ip_address = "10.27.214.158"
}

```

## Argument Reference

The following arguments are supported:

-   `ip_address` - (Required) The IP address of network interface.

## Attributes Reference

The following attributes are exported:

-   `id` - The ID of network interface.
-   `name`- The name of network interface.
-   `network_id` - The Network ID of network interface.
-   `attached_server` - The attached server of network interface.
-   `fixed_ip` - The fixed IP of network interface.
-   `mac_address` - The media access control address of network interface.
-   `admin_state_up` - The admin state up the network interface: true or false.
-   `status` - The status of network interface.
-   `device_id` - The device ID of network interface.
-   `port_security_enabled` - The port security enabled of network interface.
-   `fixed_ips` - The fixed ips of network interface.
    -   `subnet_id` - The subnet ID of network interface.
    -   `ip_address` - The IP address of network interface.
-   `security_groups` - List ID of security groups.
-   `create_at` - The created time of network interface.
-   `updated_at` - The updated time of network interface.
