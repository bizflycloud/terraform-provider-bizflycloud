---
layout: "bizflycloud"
page_title: "Bizfly Cloud: bizflycloud_vpc_network"
sidebar_current: "docs-bizflycloud-resource-vpc-network"
description: |-
  Provides a Bizfly Cloud VPC Network resource. This can be used to create, modify, and delete VPC Networks.
---

# bizflycloud\_vpc\_network

Provides a Bizfly Cloud VPC Network resource. This can be used to create,
modify, and delete VPC Network.

## Example Usage

```hcl
# Create a new VPC Network
resource "bizflycloud_vpc_network" "vpc_network" {
    name = var.vpc_network_name
    description = "test vpc network"
    cidr = "10.108.16.0/20"
    is_default = false
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of VPC Network.
* `description` - (Optional) The description of VPC Network.
* `cidr` - (Optional) CIDR Block: IPv4 or IPv6 CIDR. 
* `is_default` - (Optional) The default of VPC Network: true or false.

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
