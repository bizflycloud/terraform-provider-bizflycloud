---
layout: "bizflycloud"
page_title: "BizFly Cloud: bizflycloud_kubernetes"
sidebar_current: "docs-bizflycloud-resource-kubernetes"
description: - Provide a BizFly Cloud Kubernetes Engine resource. This can be used to create, modify, and delete
Clusters.
---

# bizflycloud\_kubernetes

Provides a BizFly Cloud Kubernetes Engine resource. This can be used to create, modify, and delete Cluster.

## Example
 ### Create a new cluster
```json
{
  "resource": {
    "bizflycloud_kubernetes": {
      "test_k8s": {
        "name": "tung491-test-k8s_23",
        "version": "5f7d3a91d857155ad4993a32",
        "vpc_network_id": "145bed1f-a7f7-4f88-ab3d-ce2fc95a4e71",
        "tags": [
          "tag"
        ],
        "worker_pools": [
          {
            "name": "pool-j9ngr550sssss",
            "flavor": "nix.2c_2g",
            "profile_type": "premium",
            "volume_type": "PREMIUM-HDD1",
            "volume_size": 40,
            "availability_zone": "HN1",
            "desired_size": 1,
            "enable_autoscaling": true,
            "min_size": 1,
            "max_size": 3,
            "tags": [
              "ssss"
            ]
          },
          {
            "name": "pool-j9ngr55ssss3",
            "flavor": "nix.2c_2g",
            "profile_type": "premium",
            "volume_type": "PREMIUM-HDD1",
            "volume_size": 40,
            "availability_zone": "HN1",
            "desired_size": 1,
            "enable_autoscaling": true,
            "min_size": 1,
            "max_size": 3,
            "tags": [
              "ssss"
            ]
          }
        ]
      }
    }
  }
}
```

### Argument Reference

The following arguments are supported:

* `name` - (Required) The cluster name.
* `version` - (Required) The Version id
* `tags` - (Optional) The tags of cluster
* `vpc_network_id` - (Required) The VPC network id
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
  

### Atrributes Reference

The following attributes are exported:

* `name` - The cluster name.
* `version` - The Version id
* `create_at` - The created time
* `created_by` - The person creating cluster
* `tags` - The tags of cluster
* `vpc_network_id` - The VPC network id
* `worker_pools` - The worker pools of Cluster
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
  