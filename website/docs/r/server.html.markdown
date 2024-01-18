---
layout: "bizflycloud"
page_title: "Bizfly Cloud: bizflycloud_server"
sidebar_current: "docs-bizflycloud-resource-server"
description: |-
Provides a Bizfly Cloud Server resource. This can be used to create, modify, and delete Servers. Servers also support
provisioning.
---

# bizflycloud\_server

Provides a Bizfly Cloud Server resource. This can be used to create,
modify, and delete Server. Servers also support
[provisioning](/docs/provisioners/index.html).

## Example Usage

```hcl
# Create a new Web Server
resource "bizflycloud_server" "tf_server1" {
  name                   = "tf_server_1"
  flavor_name            = "2c_2g"
  ssh_key                = data.bizflycloud_ssh_key.ssh_key.name
  os_type                = "image"
  os_id                  = "d646476d-850c-423e-b02c-6b86aeda3717"
  category               = "premium"
  availability_zone      = "HN1"
  root_disk_volume_type  = data.bizflycloud_volume_type.example_volume_type.type
  root_disk_size         = 20
  network_plan           = "free_bandwidth"
  billing_plan           = "on_demand"
  vpc_network_ids        = [data.bizflycloud_vpc_network.vpc_network.id, data.bizflycloud_vpc_network.vpc_network_1.id]
  state                  = "running"
  default_public_ipv4 {
    enabled = true
    firewall_ids = [
      "0ca6eaad-2941-4141-a537-d6623322ed8c"
    ]
  }
  default_public_ipv6 {
    enabled = true
    firewall_ids = [
      "56cf1f22-d8cb-41c3-948a-cc03582f0adc"
    ]
  }
}
```

## Argument Reference

The following arguments are supported:

* `os_type` - (Required) The type for create server root disk: image, snapshot, rootdisk
* `os_id` - (Required) The ID of OS - image ID, snapshot ID or volume rootdisk ID
* `name` - (Required) The Server name.
* `flavor_name` - (Required) The flavor of your server. The format for flavor is xc_yg, x is number of CPU, and y is GB
  of RAM.
* `category` - (Required) The category of a server: basic, premium, enterprise
* `ssh_key` - (Optional) The name of SSH Key for the server
* `availability_zone` - (Required) The availability zone of the server. Example: HN1, HN2, HCM1
* `root_disk_type` - (Deprecated) The type of Root disk volume: SSD or HDD
* `root_disk_volume_type` - (Required) The type of root disk volume. Get from data source volume type
* `root_disk_size` - (Required) The size of Root disk volume.
* `volume_ids` - (Optional) A list of the attached block storage volumes
* `network_plan` - (Optional) The network plan for the server. The default value is free_datatransfer.
* `vpc_network_ids` - (Optional) A list of the VPC network IDs.
* `billing_plan` - (Optional) The billing plan applied for the server (saving_plan/on_demand). Default value is
  saving_plan
* `user_data` - (Optional) The user data to provide when launching the server.
* `state` - (Optional) The state of server (running/stopped). Default value is running
* `default_public_ipv4` - (Optional) The default public IPv4 WAN network interface of the server.
  - `firewall_ids` - (Optional) A list of the firewall IDs of the network interface.
  - `enabled` - (Optional) The enabled public IPv4 WAN (true/false). Default value is true.
* `default_private_ipv6` - (Optional) The default private IPv6 LAN network interface of the server.
  - `firewall_ids` - (Optional) A list of the firewall IDs of the network interface.
  - `enabled` - (Optional) The enabled private IPv6 WAN (true/false). Default value is true.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Server
* `name`- The name of the Server
* `flavor_name` - The flavor of the server
* `category` - The category of the server
* `status` - The status of the Server
* `root_disk_id` - The ID of Server root disk
* `root_disk_type` - The type of Server root disk
* `root_disk_size` - The size of Server root disk
* `availability_zone` - The availability zone of server
* `volume_ids` - A list of the attached block storage volumes
* `default_public_ipv4` - The default public IPv4 WAN network interface of the server.
  - `id` - The ID of the IPv4 WAN.
  - `firewall_ids` - A list of the firewall IDs of the network interface.
  - `enabled` - The enabled public IPv4 WAN.
  - `ip_address` - The IPv4 WAN address.
* `default_private_ipv6` - The default private IPv6 LAN network interface of the server.
  - `id` - The ID of the IPv6 WAN.
  - `firewall_ids` - A list of the firewall IDs of the network interface.
  - `enabled` - The enable private IPv6 WAN.
  - `ip_address` - The IPv6 WAN address.
* `network_interface_ids` - A list of the network interfaces
* `network_plan` - The network plan for the server. The default value is free_datatransfer.
* `vpc_network_ids` - A list of the VPC network IDs.
* `billing_plan` - The billing plan applied for the server
* `is_available` - The state that the server is available (not in a VM action)
* `locked` - Is the server locked state
* `state` - The state of server.

## Import

Bizfly Cloud Server resource can be imported using the server id in the Bizfly manage dashboard

```
$ terraform import bizflycloud_server.tf_server1 server-id
```