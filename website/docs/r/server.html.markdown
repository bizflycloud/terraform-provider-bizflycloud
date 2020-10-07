---
layout: "bizflycloud"
page_title: "BizFly Cloud: bizflycloud_server"
sidebar_current: "docs-bizflycloud-resource-server"
description: |-
  Provides a BizFly Cloud Server resource. This can be used to create, modify, and delete Servers. Servers also support provisioning.
---

# bizflycloud\_server

Provides a BizFly Cloud Server resource. This can be used to create,
modify, and delete Server. Servers also support
[provisioning](/docs/provisioners/index.html).

## Example Usage

```hcl
# Create a new Web Server
resource "bizflycloud_server" "web" {
    name = "sapd-tf-server"
    flavor_name = "4c_2g"
    ssh_key = "sapd1"
    os_type = "image"
    os_id = "5f218529-ce32-4cb6-8557-920b16307d35"
    category = "premium"
    availability_zone = "HN1"
    root_disk_type = "HDD"
    root_disk_size = 20
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