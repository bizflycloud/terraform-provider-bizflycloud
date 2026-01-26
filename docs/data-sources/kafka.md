---
subcategory: Cloud Server
page_title: "Bizfly Cloud: bizflycloud_kafka"
description: |-
    Provides a Bizfly Cloud Kafka resource. This can be used to create, modify, and delete Cluster Kafka.
---

# Data Source: bizflycloud_kafka

Get information about Bizfly Cloud kafka.

## Example Usage

```hcl
# get cluster infomation from datasource
data "bizflycloud_kafka" "tf-kafka1" {
  id = "cfcec8cb-76c7-4e8c-8792-15ab0103a5bb"
}

```

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
