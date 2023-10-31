---
layout: "bizflycloud"
page_title: "Bizfly Cloud: bizflycloud_autoscaling_launch_configuration"
sidebar_current: "docs-bizflycloud-resource-autoscaling-launch-configuration"
description: - Provides a Bizfly Cloud AutoScaling Launch Configuration. This can be used to create, modify, AutoScaling Group
---

# bizflycloud\_autoscaling\_launch\_configuration

Provides a Bizfly Cloud AutoScaling Group resource. This can be used to create, modify, and delete.

## Example
```hcl
# Create a new AutoScaling Launch Configuration
resource "bizflycloud_autoscaling_launch_configuration" "basic-centos-terrafrom" {
  name              = "basic-centos-terrafrom"
  ssh_key           = "ministry"
  availability_zone = "HN1"
  flavor            = "1c_1g_basic"
  instance_type     = "basic"
  os {
    uuid        = "4cdbe57f-6ba1-4f40-a6fb-beb1ed974168"
    create_from = "image"

  }
  rootdisk {
    volume_type = "BASIC_SSD1"
    volume_size = 20
  }
  data_disks {
    volume_type = "BASIC_SSD1"
    volume_size = 20
  }
  data_disks {
    volume_type           = "BASIC_SSD1"
    volume_size           = 40
    delete_on_termination = true
  }
  user_data = "#!/bin/sh \n echo \"Hello World\" > /tmp/greeting.txt"
}
```

### Argument Reference

The following arguments are supported:

* `name` - (Required) The name of Launch Configuration
* `availability_zone` - (Required) The availability zone where a cloud server to be allocated
* `data_disks` - (Optional) The data disks using for cloud server
* `flavor` - (Required) The flavor of cloud server.
* `instance_type` - (Required) The type of a server: `basic`, `premium`, `enterprise` or `dedicated`
* `network_plan` - (Optional) The network plan using for server: `free_datatransfer` or `free_bandwidth`
* `networks` - (Optional) The custom network interface with security groups and choos vpc networks
    - `network_id` - (Required) The network ID using create a interface for server with security groups (firewall)
    - `security_groups` - (Required) List ID of security groups
* `os` - (Required) The information of OS
* `rootdisk` - (Required) The root disks using for cloud server
* `ssh_key` - (Required) The name of SSH Key using to be injected to cloud server
* `user_data` - (Optional) The script with text format to be injected to cloud server and run each when server start


### Atrributes Reference

The following attributes are exported:

* `id` - The ID of Launch Configuration
* `name` - The name of Launch Configuration
* `availability_zone` - The availability zone where a cloud server to be allocated. Included: `HN1` and `HN2` with region `HaNoi` or `HCM1` with region `HoChiMinh`
* `data_disks` - The data disks using for cloud server
    - `delete_on_termination` - Delete this disk when cloud server being deleted
    - `volume_size` - The size of data disk
    - `volume_type` - The type of data disk included:
        - `SSD1`
        - `HDD1`
        - `BASIC_SSD1`
        - `BASIC_HDD1`
        - `ENTERPRISE-HDD1`
        - `ENTERPRISE-SSD1`
        - `DEDICATED-SSD1`
        - `DEDICATED-HDD1`
* `flavor` - The flavor of your server. The format for flavor is `xc_yg`, `x` is number of CPU, and `y` is GB of RAM.
* `instance_type` - The type of a server: `basic`, `premium`, `enterprise` or `dedicated`
* `network_plan` - The network plan using for server: `free_datatransfer` or `free_bandwidth`
* `networks` - The custom network interface with security groups and choos vpc networks
    - `network_id` - The network ID using create a interface for server with security groups (firewall). Note: With public network will have ID is `wan`
    - `security_groups` - List ID of security groups
* `os` - The information of OS
    - `create_from` - Location to get resource to create rootdisk: `image` or `snapshot`
    - `error` - Error message when launch configuration was not available to use
    - `uuid` - The ID of resource
    - `os_name` - Information about OS
* `rootdisk` - The root disks using for cloud server
    - `delete_on_termination` - Delete this disk when cloud server being deleted
    - `volume_size` - The size of root disk
    - `volume_type` - The type of root disk included:
        - `SSD1`
        - `HDD1`
        - `BASIC_SSD1`
        - `BASIC_HDD1`
        - `ENTERPRISE-HDD1`
        - `ENTERPRISE-SSD1`
        - `DEDICATED-SSD1`
        - `DEDICATED-HDD1`
* `ssh_key` - The name of SSH Key using to be injected to cloud server
* `status` - Status of Launch Configuration
* `user_data` - The script with text format to be injected to cloud server and run each when server start
