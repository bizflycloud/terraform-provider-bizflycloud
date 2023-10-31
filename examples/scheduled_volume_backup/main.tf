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

resource "bizflycloud_scheduled_volume_backup" "backup_test" {
  volume_id      = "11a2e71b-8701-47a0-b247-41843db17e54"
  frequency      = "2880"
  size           = "4"
  scheduled_hour = 4
}
