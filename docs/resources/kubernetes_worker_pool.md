---
subcategory: Cloud Kubernetes Engine
page_title: "Bizfly Cloud: bizflycloud_kubernetes_worker_pool"
description: |-
    Provide a Bizfly Cloud Worker Pool of Kubernetes Engine resource. This can be used to create, modify, and delete Worker Pool.
---

# Resource: bizflycloud_kubernetes_worker_pool

Provide a Bizfly Cloud Worker Pool of Kubernetes Engine resource. This can be used to create, modify, and delete Worker Pool.

## Example Usage

```hcl
# Get version of the kubernetes
data "bizflycloud_kubernetes_version" "tf_k8s_version" {
  version = "v1.29.13"
}

# Get package of the kubernetes
data "bizflycloud_kubernetes_package" "tf_k8s_package" {
  provision_type = "standard"
  name = "STANDARD-1"
}

# Get VPC network
data "bizflycloud_vpc_network" "tf_vpc" {
  cidr = "10.20.2.0/24"
}

resource "bizflycloud_kubernetes" "tf_cluster" {
  name           = "cluster-name"
  version        = data.bizflycloud_kubernetes_version.tf_k8s_version.id
  vpc_network_id = data.bizflycloud_vpc_network.tf_vpc.id
  tags           = ["tag-name"]
  package_id     = data.bizflycloud_kubernetes_package.tf_k8s_package.id

  worker_pool {
    availability_zone  = "HN1"
    billing_plan       = "on_demand"
    desired_size       = 1
    enable_autoscaling = true
    flavor             = "nix.2c_2g"
    labels             = {
        "label-key" = "label-value"
    }
    max_size           = 3
    min_size           = 1
    name               = "pool-69645"
    network_plan       = "free_datatransfer"
    profile_type       = "premium"
    tags               = [
        "pool_tag"
    ]
    volume_size        = 40
    volume_type        = "PREMIUM-HDD1"

    taints {
        effect = "NoSchedule"
        key    = "taint-key"
        value  = "taint-value"
    }
  }
}


resource "bizflycloud_kubernetes_worker_pool" "tf_k8s_pool" {
  cluster_id         = resource.bizflycloud_kubernetes.tf_cluster.id
  availability_zone  = "HN1"
  billing_plan       = "on_demand"
  desired_size       = 1
  enable_autoscaling = true
  flavor             = "nix.2c_2g"
  labels             = {
    "label-key" = "label-value"
  }
  max_size           = 3
  min_size           = 1
  name               = "pool-12345"
  network_plan       = "free_datatransfer"
  profile_type       = "premium"
  tags               = [
    "pool_tag"
  ]
  volume_size        = 40
  volume_type        = "PREMIUM-HDD1"

  taints {
    effect = "NoSchedule"
    key    = "taint-key"
    value  = "taint-value"
  }
} 
```

## Argument Reference

The following arguments are supported:
-   `cluster_id` - (Required) The unique ID of the Kubernetes cluster to which this worker pool belongs.
-   `name` - (Required) The name of the worker pool, used to differentiate between multiple pools within the same cluster.
-   `flavor` - (Required) The specification (flavor) for the worker nodes in this pool.
-   `profile_type` - (Required) The profile type defining the characteristics (category) of the worker pool.
-   `volume_type` - (Required) The type of storage volume attached to the worker nodes.
-   `volume_size` - (Required) The size of the attached storage volume (in GB).
-   `availability_zone` - (Required) The specific availability zone in which the worker nodes are deployed.
-   `enable_autoscaling` - (Optional) Determines whether autoscaling is enabled for this worker pool (`true` or `false`). Default: `false`.
-   `min_size` - (Optional) The minimum number of nodes allowed in the worker pool (applicable when autoscaling is enabled).
-   `max_size` - (Optional) The maximum number of nodes allowed in the worker pool (applicable when autoscaling is enabled).
-   `tags` - (Optional) Custom metadata tags assigned to the worker pool.
-   `labels` - (Optional) Key-value pairs assigned to the worker nodes for identification and grouping.
-   `taints` - (Optional) Scheduling constraints applied to the worker nodes to control which workloads can be scheduled on them.
    -   `effect` - (Required) Defines how the taint affects pod scheduling (`NoSchedule`, `PreferNoSchedule`, or `NoExecute`).
    -   `key` - (Required) The taint key.
    -   `value` - (Optional) The taint value.
-   `desired_size` - (Required) The desired number of nodes in the worker pool.
-   `network_plan` - (Optional) The networking plan for the worker pool. Possible values: `free_datatransfer` (Data transfer is free) | `free_bandwidth` (Bandwidth usage is free). Default: `free_datatransfer`.
-   `billing_plan` - (Optional) The billing model for the worker pool. Possible values: `saving_plan` (Cost-efficient pricing) | `on_demand` (Pay-as-you-go pricing). Default: `on_demand`.

## Attributes Reference

When a worker pool is created, the following attributes are returned:

-   `id` - The unique identifier of the worker pool.
-   `cluster_id` - The ID of the cluster this worker pool belongs to.
-   `name` - The name of the worker pool.
-   `flavor` - The specification (flavor) of the worker nodes.
-   `profile_type` - The profile type (category) of the worker pool.
-   `volume_type` - The type of storage volume.
-   `volume_size` - The size of the storage volume (in GB).
-   `availability_zone` - The deployment zone of the worker nodes.
-   `enable_autoscaling` - Whether autoscaling is enabled.
-   `min_size` - The minimum number of nodes in the pool.
-   `max_size` - The maximum number of nodes in the pool.
-   `tags` - The assigned metadata tags.
-   `labels` - Key-value labels applied to the nodes.
-   `taints` - Applied scheduling constraints.
    -   `effect` - The effect of the taint.
    -   `key` - The taint key.
    -   `value` - The taint value.
-   `network_plan` - The selected network plan.
-   `billing_plan` - The selected billing model.
-   `provision_status` - The current status of the worker pool.

## Importing a Worker Pool

A Kubernetes Worker Pool on Bizfly Cloud can be imported using the worker pool ID.

```
$ terraform import bizflycloud_kubernetes_worker_pool.tf_k8s_pool pool-id
```
