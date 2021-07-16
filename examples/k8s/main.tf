terraform {
  required_providers {
    bizflycloud = {
      source = "bizflycloud/bizflycloud"
    }
  }
}

provider "bizflycloud" {
  auth_method = "password"
  region_name = "HN"
  email       = "svtt.tungds@vccloud.vn"
  password    = "NEq.c151{[Mu"
}

data "bizflycloud_kubernetes_version" "test_k8s_version" {
  version = "v1.17.9"
}

resource "bizflycloud_kubernetes" "test_k8s_cluster" {
  name           = "tung491-test-k8s_25"
  version        = data.bizflycloud_kubernetes_version.test_k8s_version.id
  vpc_network_id = "145bed1f-a7f7-4f88-ab3d-ce2fc95a4e71"
  tags           = ["tags", "123"]
  worker_pools {
    name               = "pool-69643"
    flavor             = "nix.2c_4g"
    profile_type       = "premium"
    volume_type        = "PREMIUM-HDD1"
    volume_size        = 40
    availability_zone  = "HN1"
    desired_size       = 1
    enable_autoscaling = true
    min_size           = 1
    max_size           = 3
    tags               = ["pool_tag"]
  }
  worker_pools {
    name               = "pool-9696"
    flavor             = "nix.2c_2g"
    profile_type       = "premium"
    volume_type        = "PREMIUM-HDD1"
    volume_size        = 40
    availability_zone  = "HN1"
    desired_size       = 1
    enable_autoscaling = true
    min_size           = 1
    max_size           = 3
    tags               = ["pool_tag"]
  }
}