terraform {
  required_providers {
    bizflycloud = {
      version = ">= 0.0.5"
      source  = "bizflycloud/bizflycloud"
    }
  }
}

provider "bizflycloud" {
    auth_method = "password"
    region_name = "HN"
    email = "username"
    password = ""
}

data "bizflycloud_vpc_network" "vpc_network" {
  name = bizflycloud_vpc_network.vpc_network.name
}

resource "bizflycloud_vpc_network" "vpc_network" {
    name = var.vpc_network_name
    description = "test vpc network"
    is_default = false
}