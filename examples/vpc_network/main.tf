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
    cidr = "10.108.16.0/20"
}

resource "bizflycloud_network_interface" "network_interface" {
  name = "test-name"
  network_id = "${bizflycloud_vpc_network.vpc_network.id}"
  attached_server = "21da0a9e-a59f-456f-a4c3-a0248a29eb9c"
  fixed_ip = "10.108.16.5"
  action = "attach_server"
  security_groups = ["4b41c931-bf3d-443f-b311-df3817a3fbc0"]
}

data "bizflycloud_network_interface" "network_interface" {
  network_id = bizflycloud_network_interface.network_interface.network_id
}

