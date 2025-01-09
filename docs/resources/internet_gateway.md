---
subcategory: Cloud Server
page_title: "Bizfly Cloud: bizflycloud_internet_gateway"
description: |-
    Provides a Bizfly Cloud Internet Gateway resource. This can be used to create,
modify, and delete Internet Gateway.
---

# Resource: bizflycloud_internet_gateway

Provides a Bizfly Cloud Internet Gateway resource. This can be used to create,
modify, and delete Internet Gateway.

## Example Usage

```hcl
# Create a new Internet Gateway for the VPC network
data "bizflycloud_vpc_network" "tf_vpc" {
    cidr = "10.20.2.0/24"
}

resource "bizflycloud_internet_gateway" "tf_igw" {
    name = "igw-name"
    description = "Internet gateway for tf_vpc"
    vpc_network_id = data.bizflycloud_vpc_network.tf_vpc.id
}
```

## Argument Reference

The following arguments are supported:

-   `name` - (Required) The name of internet gateway
-   `description` - (Optional) The description of internet gateway
-   `vpc_network_id` - (optional) The ID of VPC network. Attach the internet gateway to VPC network.

## Attributes Reference

The following attributes are exported:

-   `id` - The ID of the Internet Gateway
-   `name`- The name of the Internet Gateway
-   `description` - The description of the Internet Gateway
-   `vpc_network_id` - The ID of VPC network
-   `vpc_network_name` - The name of VPC network
-   `project_id` - The ID of project
-   `status` - The status of the Internet Gateway
-   `tags` - The tags of the Internet Gateway
-   `availability_zones` - The availability zones of the Internet Gateway
-   `created_at` - The created time of the Internet Gateway
-   `updated_at` - The updated time of the Internet Gateway

## Import

Bizfly Cloud Internet gateway resource can be imported using the id of Internet gateway in the Bizfly manage dashboard

```
$ terraform import bizflycloud_internet_gateway.tf_igw internet-gateway-id
```
