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

provider "bizflycloud" {
  auth_method = "password"
  region_name = "HaNoi"
  email       = var.EMAIL
  password    = var.PASSWORD
}

data "bizflycloud_vpc_network" "ducnv" {
    cidr = "10.20.2.0/24"
}

resource "bizflycloud_internet_gateway" "tf_igw" {
    name = "igw03"
    description = "Internet gateway 03"
    vpc_network_id = data.bizflycloud_vpc_network.ducnv.id
}
