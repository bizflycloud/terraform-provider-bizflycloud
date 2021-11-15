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

resource "bizflycloud_server" "sapd-server" {
  name              = "sapd-tf-server-2"
  flavor_name       = "4c_2g"
  ssh_key           = data.bizflycloud_ssh_key.ssh_key.name
  os_type           = "image"
  os_id             = "5f218529-ce32-4cb6-8557-920b16307d35"
  category          = "premium"
  availability_zone = "HN1"
  root_disk_type    = "HDD"
  root_disk_size    = 20
  network_plan      = "free_bandwidth"
}

resource "bizflycloud_volume" "volume1" {
  name              = "sapd-volume-tf5"
  size              = 20
  type              = "HDD"
  category          = "premium"
  availability_zone = "HN2"
}
