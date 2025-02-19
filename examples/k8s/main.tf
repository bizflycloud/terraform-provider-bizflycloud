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

resource "bizflycloud_kubernetes" "ducnv" {
  name           = "ducnv-cluster"
  version        = data.bizflycloud_kubernetes_version.tf_k8s_version.id
  vpc_network_id = data.bizflycloud_vpc_network.tf_vpc.id
  tags           = ["tags", "123"]
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
      "pool_tag", ""
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
  cluster_id         = resource.bizflycloud_kubernetes.ducnv.id
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