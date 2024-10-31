---
subcategory: Cloud Server
layout: "bizflycloud"
page_title: "Bizfly Cloud: bizflycloud_vpc_network"
description: |-
  Provides a Bizfly Cloud VPC Network resource. This can be used to create, modify, and delete VPC Networks.
---

# Data Source: bizflycloud_vpc_network

Get information about Bizfly Cloud VPC Network resource.

## Example Usage

```hcl
# Get information of VPC Network from datasource
data "bizflycloud_vpc_network" "vpc_network" {
   cidr = "10.20.9.0/24"
}
```

## Argument Reference

The following arguments are supported:

* `cidr` - (Required) The CIDR of VPC Network

## Attributes Reference

The following attributes are exported:

* `id` - The ID of VPC Network
* `name`- The name of VPC Network
* `description` - The description of VPC Network
* `cidr` - CIDR Block: IPv4 or IPv6 CIDR. 
* `status` - The status of the VPC Network
* `is_default` - The default of VPC Network: true or false.
* `availability_zones` - The availability zones of the VPC Network
* `mtu` - The maximum transmission unit of VPC Network.
* `subnets` - The subnets of VPC Network
  * `project_id` - The project id subnets of VPC Network.
  * `ip_version` - The IP version subnets of VPC Network.
  * `gateway_ip` - The IP gateway subnets of VPC Network.
  * `allocation_pools` - The allocation pools subnets of VPC Network.
* `create_at` - The created time.
* `updated_at` - The updated time.
