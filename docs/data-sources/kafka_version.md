---
subcategory: Cloud Kafka
page_title: "Bizfly Cloud: bizflycloud_kafka"
description: |-
    Provides a Bizfly Cloud Kafka resource. This can be used to create, modify, and delete Cluster Kafka.
---

# Data Source: bizflycloud_kafka_version

Get available version Kafka about Bizfly Cloud Kafka. The version_id use to create cluster.

## Example Usage

```hcl
# Get available kafka version from datasource
data "bizflycloud_kafka_version" "available_version" {
}

```

## Attributes Reference

The following attributes are exported:

-   `id` - The ID of the Version
-   `name`- The name of the Version
-   `code` - The code version
