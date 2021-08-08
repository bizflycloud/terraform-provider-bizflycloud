---
layout: "bizflycloud"
page_title: "BizFly Cloud: bizflycloud_server"
sidebar_current: "docs-bizflycloud-datasource-server"
description: |-
  Provides a BizFly Cloud Server resource. This can be used to create, modify, and delete Servers. Servers also support provisioning.
---

# bizflycloud\_server

Get ìnformation about BizFly Cloud OS Server. The server can be use to boot a cloud server or create a rootdisk volume from an OS server. 

## Example Usage

```hcl
# Create a new server with OS Server ID get from datasource
data "bizflycloud_server" "ubuntu18" {
    id = ""
}

```

## Argument Reference

The following arguments are supported:

* `id` - The server ID

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