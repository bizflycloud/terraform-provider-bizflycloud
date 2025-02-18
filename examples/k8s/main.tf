terraform {
  required_providers {
    bizflycloud = {
      source = "bizflycloud/bizflycloud"
    }
  }
}

variable "EMAIL" {
  type = string
}

variable "PASSWORD" {
  type = string
}

provider "bizflycloud" {
  auth_method = "password"
  region_name = "HaNoi"
  email       = var.EMAIL
  password    = var.PASSWORD
}

data "bizflycloud_kubernetes_version" "test_k8s_version" {
  version = "v1.29.13"
}

data "bizflycloud_kubernetes_package" "standard_1" {
  provision_type = "standard"
  name = "STANDARD-1"
}

resource "bizflycloud_kubernetes" "ducnv" {
  name           = "ducnv-cluster"
  version        = data.bizflycloud_kubernetes_version.test_k8s_version.id
  vpc_network_id = "aa6f8cd0-98de-42ab-aa3d-5617d3fa66d2"
  tags           = ["tags", "123"]
  package_id     = data.bizflycloud_kubernetes_package.standard_1.id

  worker_pools {
    availability_zone  = "HN1"
    billing_plan       = "on_demand"
    desired_size       = 2
    enable_autoscaling = true
    flavor             = "nix.2c_2g"
    labels             = {
        "UpdateLabel" = "UpdateLabelVal1"
    }
    max_size           = 4
    min_size           = 2
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
        key    = "UpdateTaint1"
        value  = "UpdateTaintVal1"
    }
  }
}


resource "bizflycloud_kubernetes_worker_pool" "tf_pool" {
    availability_zone  = "HN2"
    billing_plan       = "on_demand"
    cluster_id         = resource.bizflycloud_kubernetes.ducnv.id
    desired_size       = 1
    enable_autoscaling = false
    flavor             = "nix.2c_2g"
    labels             = {
      "label-key" = "label-val"
    }
    max_size           = 3
    min_size           = 1
    name               = "pool-d1blq7yd"
    network_plan       = "free_datatransfer"
    profile_type       = "premium"
    tags               = ["pool-tag"]
    volume_size        = 50
    volume_type        = "PREMIUM-HDD1"
    taints {
        effect = "NoSchedule"
        key    = "UpdateTaint1"
        value  = "UpdateTaintVal1"
    }
}

resource "bizflycloud_kubernetes_worker_pool" "tf_pool2" {
    availability_zone  = "HN2"
    billing_plan       = "on_demand"
    cluster_id         = resource.bizflycloud_kubernetes.ducnv.id
    desired_size       = 1
    enable_autoscaling = true
    flavor             = "nix.2c_2g"
    labels             = {
      "label-key" = "label-val"
    }
    max_size           = 3
    min_size           = 1
    name               = "pool-2"
    network_plan       = "free_datatransfer"
    profile_type       = "premium"
    tags               = ["pool-tag"]
    volume_size        = 50
    volume_type        = "PREMIUM-HDD1"
    taints {
        effect = "NoSchedule"
        key    = "UpdateTaint1"
        value  = "UpdateTaintVal1"
    }
}