---
layout: "bizflycloud"
page_title: "Bizfly Cloud: bizflycloud_server"
sidebar_current: "docs-bizflycloud-resource-server"
description: |-
  Provides a Bizfly Cloud Server resource. This can be used to create, modify, and delete Servers. Servers also support provisioning.
---

# bizflycloud\_server

Provides a Bizfly Cloud Server resource. This can be used to create,
modify, and delete Server. Servers also support
[provisioning](/docs/provisioners/index.html).

## Example Usage

```hcl
# Create a new Web Server
resource "bizflycloud_server" "tf_server1" {
  name                   = "tf_server_2"
  flavor_name            = "2c_2g"
  ssh_key                = data.bizflycloud_ssh_key.ssh_key.name
  os_type                = "image"
  os_id                  = "5f218529-ce32-4cb6-8557-920b16307d35"
  category               = "premium"
  availability_zone      = "HN1"
  root_disk_type         = "HDD"
  root_disk_size         = 20
  network_plan           = "free_bandwidth"
  billing_plan           = "saving_plan"
  wan_network_interfaces = [data.bizflycloud_wan_ip.wan_ip.id, data.bizflycloud_wan_ip.wan_ip_2.id]
  network_interfaces = [data.bizflycloud_network_interface.lan_ip_1.id, data.bizflycloud_network_interface.lan_ip_2.id]
}
```

## Argument Reference

The following arguments are supported:

* `os_type` - (Required) The type for create server root disk: image, snapshot, rootdisk
* `os_id` - (Required) The ID of OS - image ID, snapshot ID or volume rootdisk ID 
* `name` - (Required) The Server name.
* `flavor_name` - (Required) The flavor of your server. The format for flavor is xc_yg, x is number of CPU, and y is GB of RAM. 
* `category` - (Required) The category of a server: basic, premium, enterprise
* `ssh_key` - (Optional) The name of SSH Key for the server
* `availability_zone` - (Required) The availability zone of the server. Example: HN1, HN2, HCM1
* `root_disk_type` - (Required) The type of Root disk volume: SSD or HDD
* `root_disk_size` - (Required) The size of Root disk volume.
* `volume_ids` - (Optional) A list of the attached block storage volumes
* `network_plan` - (Optional) The network plan for the server. The default value is free_datatransfer.
* `wan_network_interfaces` - (Optional) A list of the WAN IP IDs.
* `network_interfaces` - (Optional) A list of the LAN IP IDs.
* `vpc_network_ids` - (Optional) A list of the VPC network IDs.
* `billing_plan` - (Optional) The billing plan applied for the server (saving_plan/on_demand). Default value is saving_plan

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Server
* `name`- The name of the Server
* `flavor_name` - The flavor of the server
* `category` - The category of the server
* `status` - The status of the Server
* `root_disk_type` - The type of Server root disk
* `root_disk_size` - The size of Server root disk
* `availability_zone` - The availability zone of server
* `volume_ids` - A list of the attached block storage volumes
* `lan_ip` - Lan IP of the server
* `wan_ipv4` - A list of the WAN IP v4 of the server
* `wan_ipv6` - A list of the WAN IP v6 of the server
* `network_plan` - The network plan for the server. The default value is free_datatransfer.
* `wan_network_interfaces` - A list of the WAN IP IDs.
* `network_interfaces` - A list of the LAN IP IDs.
* `vpc_network_ids` - A list of the VPC network IDs.
* `billing_plan` - The billing plan applied for the server
* `zone_name` - The zone name of the server
* `is_available` - The state that the server is available (not in a VM action)
* `locked` - Is the server locked state


## Import

Bizfly Cloud SSH key resource can be imported using the SSH key name in the BizFly manage dashboard

```
$ terraform import bizflycloud_ssh_key.example ssh-key-1
```