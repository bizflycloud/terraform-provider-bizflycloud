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
  region_name = "HN"
  email       = var.EMAIL
  password    = var.PASSWORD
}

data "bizflycloud_kubernetes_version" "test_k8s_version" {
  version = "v1.27.4"
}

resource "bizflycloud_kubernetes" "test_k8s_cluster" {
  name           = "ducnv"
  version        = data.bizflycloud_kubernetes_version.test_k8s_version.id
  vpc_network_id = "0f05d672-6e25-439b-a27f-e7ffb44a07b2"
  tags           = ["tags", "123"]
  worker_pools {
    name               = "pool-69646"
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
    labels = {
      "UpdateLabel" = "UpdateLabelVal"
    }
    taints {
      effect = "NoSchedule"
      key = "UpdateTaint"
      value = "UpdateTaintVal"
    }
  }

  worker_pools {
    name               = "pool-69645"
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
    labels = {
      "UpdateLabel" = "UpdateLabelVal"
      "label2" = "labelVal2"
    }
    taints {
      effect = "NoSchedule"
      key = "UpdateTaint2"
      value = "UpdateTaintVal2"
    }
    taints {
      effect = "NoSchedule"
      key = "UpdateTaint1"
      value = "UpdateTaintVal1"
    }
  }
}