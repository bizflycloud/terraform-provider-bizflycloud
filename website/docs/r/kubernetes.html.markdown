---
layout: "bizflycloud"
page_title: "Bizfly Cloud: bizflycloud_kubernetes"
sidebar_current: "docs-bizflycloud-resource-kubernetes"
description: |-
  Provide a Bizfly Cloud Kubernetes Engine resource. This can be used to create, modify, and delete Clusters.
---

# bizflycloud\_kubernetes

Provides a Bizfly Cloud Kubernetes Engine resource. This can be used to create, modify, and delete Cluster.

## Example Usage

```hcl
resource "bizflycloud_kubernetes" "tf_create_k8s" {
  name            = "create-ducnv"
  auto_upgrade    = false
  cni_plugin      = "kube-router"
  local_dns       = false
  version         = "64c8709c3f881935b73b43f0"
  vpc_network_id  = "e15c7244-7f16-4af6-8a67-2a31f7af38f9"
  enabled_upgrade_version = false
  worker_pools {
      availability_zone  = "HN1"
      billing_plan       = "saving_plan"
      desired_size       = 1
      enable_autoscaling = false
      flavor             = "nix.2c_2g"
      labels             = {
          "ducnv" = "123"
      }
      max_size           = 1
      min_size           = 1
      name               = "pool-name"
      network_plan       = "free_bandwidth"
      profile_type       = "premium"
      tags               = ["tag-name"]
      volume_size        = 30
      volume_type        = "PREMIUM-SSD1"

      taints {
          effect = "NoSchedule"
          key    = "duc"
          value  = "123"
      }

      taints {
          effect = "PreferNoSchedule"
          key    = "ducnv"
      }
  }
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The cluster name.
* `version` - (Required) The Version id
* `tags` - (Optional) The tags of cluster
* `vpc_network_id` - (Required) The VPC network id
* `auto_upgrade` - (Optional) The auto upgrade (true/false). Default value is false.
* `local_dns` - (Optional) The local DNS (true/false). Default value is false.
* `cni_plugin` - (Optional) The CNI plugin (kube-router/cilium). Default value is kube-router.
* `enabled_upgrade_version` - (Optional) The enabled upgrade cluster version (true/false). Default value is false
* `worker_pools` - (Required) The worker pools of Cluster
  * `name` - (Required) The worker pool name
  * `flavor` - (Required) The flavor of pool
  * `profile_type` - (Required) The profile type of pool
  * `volume_type` - (Required) The volume type
  * `volume_size` - (Required) The volume size
  * `availability_zone` - (Required) The availability zone
  * `enable_autoscaling` - (Optional) Enable auto scaling or not
  * `min_size` - (Optional) The number of the minimum node
  * `max_size` - (Optional) The number of the maximum node
  * `tags` - (Optional) The tags of the pool
  * `labels` - (Optional) The labels
  * `taints` - (Optional) The taints
    * `effect` - (Required) The effect (NoSchedule/PreferNoSchedule/NoExecute).
    * `key` - (Required) The key
    * `value` - (Optional) The value
  * `desired_size` - (Required) The desired size
  * `network_plan` - (Optional) The network plan (free_datatransfer/free_bandwidth). Default value is free_datatransfer.
  * `billing_plan` - (Optional) The billing plan (saving_plan/on_demand). Default value is on_demand.
  

## Attributes Reference

The following attributes are exported:
* `id` - The cluster ID
* `name` - The cluster name.
* `version` - The Version id
* `create_at` - The created time
* `created_by` - The person creating cluster
* `auto_upgrade` - The auto upgrade
* `local_dns` - The local DNS
* `cni_plugin` - The CNI plugin
* `tags` - The tags of cluster
* `vpc_network_id` - The VPC network id
* `enabled_upgrade_version` - The enabled upgrade cluster version
* `is_latest` - The cluster version is latest
* `current_version` - The current version of cluster
* `next_version` - The next version for upgrade cluster version
* `worker_pools` - The worker pools of Cluster
  * `id` - The worker pool ID
  * `name` - The worker pool name
  * `flavor` - The flavor of pool
  * `profile_type` - The profile type of pool
  * `volume_type` - The volume type
  * `volume_size` - The volume size
  * `availability_zone` - The availability zone
  * `enable_autoscaling` - Enable auto scaling or not
  * `min_size` - The number of the minimum node
  * `max_size` - The number of the maximum node
  * `tags` - The tags of the pool
  * `labels` - The labels
  * `taints` - The taints
    * `effect` - The effect
    * `key` - The key
    * `value` - The value
  * `network_plan` - The network plan
  * `billing_plan` - The billing plan


## Import

Bizfly Cloud kubernetes resource can be imported using the cluster id

```
$ terraform import bizflycloud_kubernetes.tf_create_k8s cluster-id
```