---
subcategory: Cloud Kafka
page_title: "Bizfly Cloud: bizflycloud_kafka"
description: |-
    Provides a Bizfly Cloud Kafka resource. This can be used to create, modify, and delete Cluster Kafka.
---

# Data Source: bizflycloud_kafka_flavor

Get available flavor Kafka about Bizfly Cloud Kafka. The id use to create cluster.

## Example Usage

```hcl
# Get available kafka flavor from datasource
data "bizflycloud_kafka" "available_flavor" {
}

```


## Attributes Reference

The following attributes are exported:

-   `id` - The ID of the Flavor
-   `name`- The name of the Flavor
-   `code` - The code of the flavor
-   `vCPU` - Number cpu of the Flavor (vCore)
-   `RAM` - Number Memmory of the Flavor (Mb)
