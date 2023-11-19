terraform {
  required_providers {
    bizflycloud = {
      version = ">= 0.0.5"
      source  = "bizflycloud/bizflycloud"
    }
  }
}

 variable "EMAIL" {
   type = string
 }

 variable "PASSWORD" {
   type = string
 }

 variable "PROJECT_ID" {
   type=string
 }

provider "bizflycloud" {
  auth_method = "password"
  region_name = "HaNoi"
  email       = var.EMAIL
  password    = var.PASSWORD
}

resource "bizflycloud_firewall" "sample_firewall_1" {
  name = "sample-firewall-1"
}

data "bizflycloud_ssh_key" "ssh_key" {
  name = "test1"
}

data "bizflycloud_volume_type" "example_volume_type" {
  name     = "SSD"
  category = "premium"
}

resource "bizflycloud_wan_ip" "test_wan" {
  availability_zone = "HN1"
  name              = "test_wan_5"
}

resource "bizflycloud_network_interface" "test_lan" {
  name       = "test_lan_5"
  network_id = "75286136-bfc1-448b-8626-23cc9b90aa47"
}

resource "bizflycloud_server" "tf_server1" {
  name                  = "tf_server_10"
  flavor_name           = "1c_1g"
  ssh_key               = data.bizflycloud_ssh_key.ssh_key.name
  os_type               = "image"
  os_id                 = "d646476d-850c-423e-b02c-6b86aeda3717"
  category              = "premium"
  availability_zone     = "HN1"
  root_disk_volume_type = data.bizflycloud_volume_type.example_volume_type.type
  root_disk_size        = 20
  network_plan          = "free_bandwidth"
  billing_plan          = "on_demand"
  free_wan {
    ip_version         = 4
    firewall_ids = ["348b3be8-0610-4397-a781-10be10998b90", "9ff7fbdc-1461-4713-b4f8-ba8a22699bb4"]
  }
  free_wan {
    ip_version         = 6
    firewall_ids = ["0b1d6a90-24a0-4f86-b75f-fb07290e44dd", "348b3be8-0610-4397-a781-10be10998b90"]
  }
  vpc_network_ids = ["8b4f2c75-e4dd-440c-b36c-50c9494faa8a"]
  user_data = "!/bin/bash"
}


#
#resource "bizflycloud_volume" "volume1" {
#  name              = "sapd-volume-tf5"
#  size              = 20
#  type              = data.bizflycloud_volume_type.example_volume_type.type
#  category          = "premium"
#  availability_zone = "HN2"
#}

# resource "bizflycloud_wan_ip" "test_wan_1" {
#   name              = "sapd-wan-ip-tf4"
#   availability_zone = "HN1"
#   attached_server   = "61fe3c90-7db0-47ba-b034-06de66a0869b"
# }

# data "bizflycloud_wan_ip" "wan_ip" {
#   ip_address = "103.107.183.114"
# }
#
# data "bizflycloud_wan_ip" "wan_ip_2" {
#   ip_address = "45.124.94.87"
# }
#
# data "bizflycloud_network_interface" "lan_ip_1" {
#   ip_address = "10.27.214.140"
# }
#
# data "bizflycloud_network_interface" "lan_ip_2" {
#   ip_address = "10.27.214.158"
# }
#
# data "bizflycloud_vpc_network" "vpc_network" {
#   cidr = "10.20.9.0/24"
# }
#
# data "bizflycloud_vpc_network" "vpc_network_1" {
#   cidr = "10.20.8.0/24"
# }