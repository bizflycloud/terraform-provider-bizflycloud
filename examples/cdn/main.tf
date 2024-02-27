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
  region_name = "HaNoi"
  email       = ""
  password    = ""
}

resource "bizflycloud_cdn" "domain_com" {
  domain = "cdn.domain.com"
  origin =  {
    upstream_addrs = "origin.domain.com"
    upstream_host = "origin.domain.com"
    upstream_proto = "https"
    name = "origin-domain"
  }
}