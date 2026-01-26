---
subcategory: Cloud Kafka
page_title: "Bizfly Cloud: bizflycloud_kafka"
description: |-
    Provides a Bizfly Cloud Kafka resource. This can be used to create, modify, and delete Cluster Kafka.
---

# Resource: bizflycloud_kafka

Provides a Bizfly Cloud Kafka resource. This can be used to create,
modify, and delete Cluster Kafka. Servers also support
[provisioning](/docs/provisioners/index.html).

## Example Usage

```hcl
# Create a new Cluster Kafka
resource "bizflycloud_kafka" "example" {
  name              = "tf-kafka"
  version_id        = "10d9c0ee-5b3b-4a9e-9f21-4416fbb9e8ef"
  nodes             = 1
  flavor            = "2c_4g"
  volume_size       = 10
  availability_zone = "HN1"
  vpc_network_id    = "your-vpc-network-id"
  public_access     = true
}

data "bizflycloud_kafka" "example" {
  name = bizflycloud_kafka.example.name
}

output "kafka_create_cluster_info" {
  value = {
    name = data.bizflycloud_kafka.example.name
    version_id = data.bizflycloud_kafka.example.version_id
    nodes  = data.bizflycloud_kafka.example.nodes
    volume_size = data.bizflycloud_kafka.example.volume_size
    flavor  = data.bizflycloud_kafka.example.flavor
    availability_zone = data.bizflycloud_kafka.example.availability_zone
    public_access = data.bizflycloud_kafka.example.public_access
    status = data.bizflycloud_kafka.example.status
  }
}
```

## Argument Reference

The following arguments are supported:

-   `version_id` - (Required) The ID of Kafka cluster version
-   `name` - (Required) The Server name.
-   `flavor` - (Required) The flavor of node for your cluster. The format for flavor is xc_yg, x is number of CPU, and y is GB
    of RAM.
-   `availability_zone` - (Required) The availability zone of the server. Example: HN1, HN2, HCM1
-   `volume_size` - (Required) The size of Data disk volume.
-   `nodes` - (Required) Number node of your cluster.
-   `vpc_network_id` - (Required) Your VPC Network ID.
-   `public_access` - (Required) Cluster can access from internet? (true/false). Default value is false

## Attributes Reference

The following attributes are exported:

-   `id` - The ID of the Cluster
-   `name`- The name of the Cluster
-   `flavor` - The flavor of the Cluster
-   `status` - The status of the Cluster
-   `volume_size` - The volume size each node of Cluster
-   `availability_zone` - The availability zone of Cluster
-   `nodes` - The number node of Cluster
-   `public_access` - Cluster can access from internet?

## Import

Bizfly Cloud Server resource can be imported using the v id in the Bizfly manage dashboard

```
$ terraform import bizflycloud_kafka.tf_kafka cluster-id
```
