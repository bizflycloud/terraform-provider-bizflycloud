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

resource "bizflycloud_vpc_network" "vpc_network" {
    name = "vpc_toannd"
    description = "test"
    cidr = "10.108.16.0/20"
    is_default = false
}

data "bizflycloud_vpc_network" "vpc_network" {
  id = bizflycloud_vpc_network.vpc_network.id
  name = bizflycloud_vpc_network.vpc_network.name
}