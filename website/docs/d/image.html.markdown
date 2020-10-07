---
layout: "bizflycloud"
page_title: "BizFly Cloud: bizflycloud_image"
sidebar_current: "docs-bizflycloud-datasource-image"
description: |-
  Provides a BizFly Cloud Server resource. This can be used to create, modify, and delete Servers. Servers also support provisioning.
---

# bizflycloud\_image

Get ìnformation about BizFly Cloud OS Image. The image can be use to boot a cloud server or create a rootdisk volume from an OS image. 

## Example Usage

```hcl
# Create a new server with OS Image ID get from datasource
data "bizflycloud_image" "ubuntu18" {
    distribution = "ubuntu"
    version = "18.04 x64"
}

resource "bizflycloud_server" "sapd-ubuntu-20" {
    name = "sapd-tf-server-2"
    flavor_name = "4c_2g"
    ssh_key = "sapd1"
    os_type = "image"
    os_id = "${data.bizflycloud_image.ubuntu18.id}"
    category = "premium"
    availability_zone = "HN1"
    root_disk_type = "HDD"
    root_disk_size = 20
}
```

## Argument Reference

The following arguments are supported:

* `distribution` - The distribution of OS: available: Ubuntu, CentOS, Debian and Windows
* `version` - The version of OS. Example: `20.04 x64`, `8.0 x64`

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the OS image
* `distribution` - The Distribution of the OS Image
* `version` - The version of the OS Image