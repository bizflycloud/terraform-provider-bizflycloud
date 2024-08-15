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

# variable "PROJECT_ID" {
#   type=string
# }

provider "bizflycloud" {
  auth_method = "password"
  region_name = "HaNoi"
  email       = var.EMAIL
  password    = var.PASSWORD
}

resource "bizflycloud_firewall" "sample_firewall_1" {
  name = "sample-firewall-tf-1"
}

data "bizflycloud_ssh_key" "ssh_key" {
  name = "test1"
}

data "bizflycloud_volume_type" "example_volume_type" {
  name     = "SSD"
  category = "premium"
}

data "bizflycloud_vpc_network" "vpc_network1" {
  cidr = "10.20.3.0/24"
}

data "bizflycloud_vpc_network" "vpc_network2" {
  cidr = "10.20.2.0/24"
}

resource "bizflycloud_network_interface" "test_lan1" {
  name         = "test_lan2_${count.index}"
  network_id   = data.bizflycloud_vpc_network.vpc_network1.id
  firewall_ids = [
    "9ff7fbdc-1461-4713-b4f8-ba8a22699bb4",
    "0b1d6a90-24a0-4f86-b75f-fb07290e44dd",
    "348b3be8-0610-4397-a781-10be10998b90"
  ]

  count = 2
}

resource "bizflycloud_network_interface" "test_lan2" {
  name       = "test_lan4_${count.index}"
  network_id = data.bizflycloud_vpc_network.vpc_network2.id
  count      = 2
}

resource "bizflycloud_volume" "volume1" {
  name              = "volume-tf_${count.index}"
  size              = 20
  type              = data.bizflycloud_volume_type.example_volume_type.type
  category          = "premium"
  availability_zone = "HN1"
  count             = 2
}

resource "bizflycloud_volume_attachment" "volume_attachment" {
  volume_id = bizflycloud_volume.volume1.*.id[count.index]
  server_id = bizflycloud_server.tf_server1.*.id[count.index]
  count     = 2
}

resource "bizflycloud_wan_ip" "test_wan_1" {
  name              = "sapd-wan-ip-tf5-${count.index}"
  availability_zone = "HN1"
  firewall_ids      = [
    "9ff7fbdc-1461-4713-b4f8-ba8a22699bb4",
    "0b1d6a90-24a0-4f86-b75f-fb07290e44dd"

  ]
  count = 2
}

resource "bizflycloud_server" "tf_server1" {
  count                 = 2
  name                  = "tf_server_5_${count.index}"
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
  default_public_ipv4 {
    firewall_ids = [
      "9ff7fbdc-1461-4713-b4f8-ba8a22699bb4",
    ]
  }
  default_public_ipv6 {
    firewall_ids = [
      "0b1d6a90-24a0-4f86-b75f-fb07290e44dd"
    ]
  }
  user_data = "!/bin/bash"
  network_interfaces {
    id = bizflycloud_wan_ip.test_wan_1.*.id[count.index]
    enabled = true
  }
  network_interfaces {
    id = bizflycloud_network_interface.test_lan1.*.id[count.index]
    enabled = true
  }
  network_interfaces {
    id = bizflycloud_network_interface.test_lan2.*.id[count.index]
    enabled = true
  }
}

data "bizflycloud_volume_snapshot" "volume_snapshot" {
  id = "15d628d5-41cf-4508-a7a6-f571eae235ce"
}

data "bizflycloud_custom_image" "custom_image" {
  id = "d646476d-850c-423e-b02c-6b86aeda3717"
}

resource "bizflycloud_custom_image" "new_custom_image1" {
  name        = "new_custom_image_10"
  disk_format = "qcow2"
  image_url   = "https://releases.ubuntu.com/22.04.3/ubuntu-22.04.3-desktop-amd64.iso"
}

 data "bizflycloud_wan_ip" "wan_ip" {
   ip_address = "103.107.183.114"
 }

 data "bizflycloud_network_interface" "lan_ip_1" {
   ip_address = "10.27.214.140"
 }

 data "bizflycloud_vpc_network" "vpc_network" {
   cidr = "10.20.9.0/24"
 }
