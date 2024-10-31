---
subcategory: AutoScaling
layout: "bizflycloud"
page_title: "Bizfly Cloud: bizflycloud_autoscaling_nodes"
description: |-
  Provides a Bizfly Cloud AutoScaling Nodes. This can be used to get server detail from physical id.
---

# Data Source: bizflycloud_node

Get ìnformation about Bizfly Cloud AutoScaling Nodes 

## Example Usage

```hcl
# Create a new node with Cluster ID get from datasource
data "bizflycloud_autoscaling_nodes" "ubuntu18" {
    cluster_id = ""
}

```

## Argument Reference

The following arguments are supported:

* `cluster_id` - The autoscaling ID

## Attributes Reference

The following attributes are exported:

* `name` - The Name of Node
* `id` - The ID of Node
* `profile_name` - The ProfileName of Node
* `profile_id` - The ProfileID of Node
* `physical_id` - The PhysicalID of Node
* `status` - The Status of Node
* `status_reason` - The StatusReason of Node