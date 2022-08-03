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
  email       = ""
  password    = ""
}

data "bizflycloud_ssh_key" "ssh_key" {
  name = "test1"
}

data "bizflycloud_volume_type" "example_volume_type" {
  name = "HDD"
  category = "premium"
}

resource "bizflycloud_server" "tf_server1" {
  name                   = "tf_server_7"
  flavor_name            = "2c_2g"
  ssh_key                = data.bizflycloud_ssh_key.ssh_key.name
  os_type                = "image"
  os_id                  = "d646476d-850c-423e-b02c-6b86aeda3717"
  category               = "premium"
  availability_zone      = "HN1"
  root_disk_volume_type         = data.bizflycloud_volume_type.example_volume_type.type
  root_disk_size         = 20
  network_plan           = "free_bandwidth"
  billing_plan          = "on_demand"
  wan_network_interfaces = [data.bizflycloud_wan_ip.wan_ip.id, data.bizflycloud_wan_ip.wan_ip_2.id]
  network_interfaces = [data.bizflycloud_network_interface.lan_ip_1.id, data.bizflycloud_network_interface.lan_ip_2.id]
}

resource "bizflycloud_volume" "volume1" {
  name              = "sapd-volume-tf5"
  size              = 20
  type              = data.bizflycloud_volume_type.example_volume_type.type
  category          = "premium"
  availability_zone = "HN2"
}

resource "bizflycloud_wan_ip" "test_wan_1" {
  name              = "sapd-wan-ip-tf4"
  availability_zone = "HN1"
  attached_server   = "61fe3c90-7db0-47ba-b034-06de66a0869b"
}

data "bizflycloud_wan_ip" "wan_ip" {
  ip_address = "103.148.57.46"
}

data "bizflycloud_wan_ip" "wan_ip_2" {
  ip_address = "45.124.94.87"
}

data "bizflycloud_network_interface" "lan_ip_1" {
  ip_address = "10.27.214.140"
}

data "bizflycloud_network_interface" "lan_ip_2" {
  ip_address = "10.27.214.158"
}

data "bizflycloud_vpc_network" "vpc_network" {
  cidr = "10.20.9.0/24"
}

data "bizflycloud_vpc_network" "vpc_network_1" {
  cidr = "10.20.8.0/24"
}