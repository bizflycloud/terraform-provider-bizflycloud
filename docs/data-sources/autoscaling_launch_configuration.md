---
subcategory: AutoScaling
layout: "bizflycloud"
page_title: "Bizfly Cloud: bizflycloud_launch_configuration"
description: |-
  Provides a Bizfly Cloud AutoScaling Launch Configuration. This can be used to create, modify, AutoScaling group
---

# Data Source: bizflycloud_launch_configuration

Get ìnformation about Bizfly Cloud AutoScaling Launch Configuration. The launch configuration include about information can be use to boot a cloud server in AutoScaling Group

## Example Usage

```hcl
# Get information of AutoScaling Launch Configuration from datasource
data "bizflycloud_autoscaling_launch_configuration" "basic-centos-4c-4g" {
    id = "1025269d-9eca-4b4d-bc4f-6edc6ae9f0c6"
    name = "basic-centos-4c-4g"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The ID of Launch Configuration

## Attributes Reference

The following attributes are exported:

* `id` - The ID of Launch Configuration
* `name` - The name of Launch Configuration
* `availability_zone` - The availability zone where a cloud server to be allocated
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
* `flavor` - The flavor of cloud server. The format for flavor is `xc_yg`, `x` is number of CPU, and `y` is GB of RAM.
* `instance_type` - The type of a server: `basic`, `premium`, `enterprise` or `dedicated`
* `network_plan` - The network plan using for server: `free_datatransfer` or `free_bandwidth`
* `networks` - The custom network interface with security groups and choos vpc networks
    - `network_id` - The network ID using create a interface for server with security groups (firewall)
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
