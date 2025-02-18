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

-   `name` - (Required) The cluster name.
-   `version` - (Required) The Version id
-   `package_id` - (Required) The Package id
-   `tags` - (Optional) The tags of cluster
-   `vpc_network_id` - (Required) The VPC network id
-   `auto_upgrade` - (Optional) The auto upgrade (true/false). Default value is false.
-   `local_dns` - (Optional) The local DNS (true/false). Default value is false.
-   `cni_plugin` - (Optional) The CNI plugin (kube-router/cilium). Default value is kube-router.
-   `enabled_upgrade_version` - (Optional) The enabled upgrade cluster version (true/false). Default value is false
-   `worker_pool` - (Required) A worker pool of the cluster.
    -   `name` - (Required) The worker pool name
    -   `flavor` - (Required) The flavor of pool
    -   `profile_type` - (Required) The profile type of pool
    -   `volume_type` - (Required) The volume type
    -   `volume_size` - (Required) The volume size
    -   `availability_zone` - (Required) The availability zone
    -   `enable_autoscaling` - (Optional) Enable auto scaling or not
    -   `min_size` - (Optional) The number of the minimum node
    -   `max_size` - (Optional) The number of the maximum node
    -   `tags` - (Optional) The tags of the pool
    -   `labels` - (Optional) The labels
    -   `taints` - (Optional) The taints
        -   `effect` - (Required) The effect (NoSchedule/PreferNoSchedule/NoExecute).
        -   `key` - (Required) The key
        -   `value` - (Optional) The value
    -   `desired_size` - (Required) The desired size
    -   `network_plan` - (Optional) The network plan (free_datatransfer/free_bandwidth). Default value is free_datatransfer.
    -   `billing_plan` - (Optional) The billing plan (saving_plan/on_demand). Default value is on_demand.

## Attributes Reference

The following attributes are exported:

-   `id` - The cluster ID
-   `name` - The cluster name.
-   `version` - The Version id
-   `package_id` - The Package id
-   `create_at` - The created time
-   `created_by` - The person creating cluster
-   `auto_upgrade` - The auto upgrade
-   `local_dns` - The local DNS
-   `cni_plugin` - The CNI plugin
-   `tags` - The tags of cluster
-   `vpc_network_id` - The VPC network id
-   `enabled_upgrade_version` - The enabled upgrade cluster version
-   `is_latest` - The cluster version is latest
-   `current_version` - The current version of cluster
-   `next_version` - The next version for upgrade cluster version
-   `worker_pool` - A worker pool of the cluster.
    -   `id` - The worker pool ID
    -   `name` - The worker pool name
    -   `flavor` - The flavor of pool
    -   `profile_type` - The profile type of pool
    -   `volume_type` - The volume type
    -   `volume_size` - The volume size
    -   `availability_zone` - The availability zone
    -   `enable_autoscaling` - Enable auto scaling or not
    -   `min_size` - The number of the minimum node
    -   `max_size` - The number of the maximum node
    -   `tags` - The tags of the pool
    -   `labels` - The labels
    -   `taints` - The taints
        -   `effect` - The effect
        -   `key` - The key
        -   `value` - The value
    -   `network_plan` - The network plan
    -   `billing_plan` - The billing plan

## Import

Bizfly Cloud kubernetes resource can be imported using the cluster id

```
$ terraform import bizflycloud_kubernetes.tf_cluster cluster-id
```
