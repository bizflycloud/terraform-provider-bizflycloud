terraform {
  required_providers {
    bizflycloud = {
      source = "tung491/bizflycloud"
    }
  }
}


provider "bizflycloud" {
  auth_method = "password"
  region_name = "HN"
  email = "svtt.tungds@vccloud.vn"
  password = "NEq.c151{[Mu"
}

resource "bizflycloud_kubernetes" "test_k8s" {
  name = "tung491-test-k8s"
  version = "5f6425f3d0d3befd40e7a31f"
  auto_upgrade = false
  tags = [
    "string"]
  worker_pools = [
    {
      name = "pool"
      version = "v1.18.0"
      flavor = "8c_8g"
      profile_type = "premium"
      volume_type = "SSD"
      volume_size = 40,
      availability_zone = "HN1"
      desired_size = 1
      enable_autoscaling = true
      min_size = 1
      max_size = 3
      tags = [
        "string"]
    }
  ]
}
