---
layout: "bizflycloud"
subcategory: Cloud Database
page_title: "Bizfly Cloud: bizflycloud_cloud_database_node"
sidebar_current: "docs-bizflycloud-datasource-cloud-database"
description: |-
  Provides a Bizfly Cloud Database Nodes. This can be used to get database node detail from database node id.
---

# bizflycloud\_cloud\_database\_node

Get ìnformation about Bizfly Cloud Database Node

## Example Usage

```hcl
# Get information of Cloud Database Node from datasource
data "bizflycloud_cloud_database_node" "mongo604-primary" {
  id = "7201f1dd-6185-46a2-b488-87f2a4fce70e"
}
```

## Argument Reference

The following arguments are supported:

* `id` - ID of this database node
* `name` - Name of this database node (option)

## Attributes Reference

The following attributes are exported:

* `availability_zone` - The data center - that allocating this database node
* `created_at` - The time that init this database node
* `datastore` - The datastore of database node
* `dns` - The DNS of this database node
* `flavor` - The setting CPU/ RAM of this database node
* `instance_id` - The database instance ID - that this node is member
* `node_type` - The node type of this database node
* `operating_status` - The operating status of this database node
* `port_access` - The port to access this database node
* `private_addresses` - The list private addresses to access this database node
* `public_addresses` - The list public addresses to access this database node
* `region_name` - The region of data center - that allocating this database node
* `replica_of` - The ID of database node - that is source replica of this database node
* `replicas` - The list nodes - that are replica of this database node
* `role` - The role of this database node. Have some role like: `primary`, `secondary`, `replica`
* `status` - The status of this database node. Have some status like: `ACTIVE`, `RESIZE`, `ERROR`
* `volume` - The size of storage provisioning/ used of this database node, unit is `Gigabytes`