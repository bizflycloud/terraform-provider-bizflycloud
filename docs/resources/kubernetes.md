---
subcategory: Cloud Kubernetes Engine
page_title: "Bizfly Cloud: bizflycloud_kubernetes"
description: |-
    Provide a Bizfly Cloud Kubernetes Engine resource. This can be used to create, modify, and delete Clusters.
---

# Resource: bizflycloud_kubernetes

Provides a Bizfly Cloud Kubernetes Engine resource. This can be used to create, modify, and delete Cluster.

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

```

## Argument Reference

The following arguments are supported:

-   `name` - (Required) The unique name of the Kubernetes cluster.
-   `version` - (Required) The version ID of the Kubernetes cluster.
-   `package_id` - (Required) The ID of the package that defines the cluster’s resource allocation and features.
-   `tags` - (Optional) A list of custom metadata tags for the cluster.
-   `vpc_network_id` - (Required) The ID of the Virtual Private Cloud (VPC) network where the cluster is deployed.
-   `auto_upgrade` - (Optional) Enables automatic Kubernetes version upgrades for the cluster. Ensures the cluster remains up-to-date with security patches and new features. Values: `true` (enabled) | `false` (disabled). Default: `false`.
-   `local_dns` - (Optional) Enables a local DNS service for cluster name resolution. Improves internal DNS performance and reliability. Values: `true` (enabled) | `false` (disabled). Default: `false`.
-   `cni_plugin` - (Optional) Specifies the Container Network Interface (CNI) plugin used for networking. Possible values: `kube-router` (A lightweight CNI plugin that provides network policies and BGP routing) | `cilium` (A more advanced CNI plugin with security features and network observability). Default: `kube-router`.
-   `enabled_upgrade_version` - (Optional) Allows upgrading the cluster version. If enabled, the cluster can be manually upgraded to a new Kubernetes version. Values: `true` (enabled) | `false` (disabled). Default: `false`.
-   `worker_pool` - (Required) A *worker pool* defines a set of worker nodes that handle workloads within the cluster. Additional worker pools may be added to the cluster using the **bizflycloud_kubernetes_worker_pool** resource. The following arguments may be specified:
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
    -   `network_plan` - (Optional) The networking plan for the worker pool. Possible values: `free_datatransfer` (default): Data transfer is free | `free_bandwidth`: Bandwidth usage is free.
    -   `billing_plan` - (Optional) The billing model for the worker pool. Possible values:
`saving_plan` (Cost-efficient pricing) | `on_demand` (Pay-as-you-go pricing). Default: `on_demand`
    
    **Note**: Update when the following fields change: `min_size`, `max_size`, `desired_size`, `labels` or `taints`. Replace when the following fields change: `name`, `flavor`, `profile_type`, `volume_type`, `volume_size`, `availability_zone`, `network_plan` or `billing_plan`.

## Attributes Reference

The following attributes are exported:

-   `id` - The unique identifier assigned to the Kubernetes cluster.
-   `name` - The name of the Kubernetes cluster.
-   `version` - The Kubernetes version ID running on the cluster.
-   `package_id` - The package ID defining the cluster’s resource allocation and configurations.
-   `create_at` - The timestamp indicating when the cluster was created.
-   `created_by` - The identifier of the user or system that created the cluster.
-   `auto_upgrade` - Indicates whether automatic upgrades are enabled for the cluster.
-   `local_dns` - Specifies whether local DNS resolution is enabled for internal services.
-   `cni_plugin` - The Container Network Interface (CNI) plugin used for networking.
-   `tags` - A list of metadata tags assigned to the cluster.
-   `vpc_network_id` - The ID of the Virtual Private Cloud (VPC) where the cluster is deployed.
-   `enabled_upgrade_version` - Indicates whether upgrading the cluster version is enabled.
-   `is_latest` - Specifies whether the cluster is running the latest available Kubernetes version.
-   `current_version` - Displays the current running version of Kubernetes in the cluster.
-   `next_version` - The next available Kubernetes version for upgrading the cluster.
-   `worker_pool` - A worker pools define the compute resources used to run workloads within the cluster.
    -   `id` - The unique identifier of the worker pool.
    -   `name` - The name assigned to the worker pool.
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

## Importing a cluster

Bizfly Cloud kubernetes resource can be imported using the cluster id

```
$ terraform import bizflycloud_kubernetes.tf_cluster cluster-id
```
