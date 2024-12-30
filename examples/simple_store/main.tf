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


resource "bizflycloud_simple_storage_bucket" "example" {
    name = "newtest"
    location = "hn"
    acl = "private"
    default_storage_class = "COLD"
}

resource "bizflycloud_simple_storage_access_key" "example" {
    subuser_id = "cuong01"
    access_key = "77777"
    secret_key = "2222"
}

resource "bizflycloud_simple_storage_bucket_acl" "example" {
    bucket_name = "newtest"
    acl = "private"
}
// public-read
// private

resource "bizflycloud_simple_storage_bucket_versioning" "example" {
    bucket_name = "newtest"
    versioning = false
}
resource "bizflycloud_simple_storage_bucket_cors" "example" {
  bucket_name = "newtest"

  rules {
    allowed_origin  = "http://ahoho.com"
    allowed_methods = ["PUT"]
    allowed_headers = ["Content-Type"]
    max_age_seconds = 6400
  }

  rules {
    allowed_origin  = "http://another-origin.com"
    allowed_methods = ["POST"]
    allowed_headers = ["Authorization"]
    max_age_seconds = 7200
  }
}

resource "bizflycloud_simple_storage_bucket_website_config" "example" {
    bucket_name = "newtest"
    index = "tttt.html"
    error = "okokokoefe"
}

