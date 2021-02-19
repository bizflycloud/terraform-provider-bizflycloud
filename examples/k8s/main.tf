terraform {
    required_providers {
        bizflycloud = {
            versions = ["0.0.1"]
            source = "bizflycloud/bizflycloud"
        }
    }
}

provider "bizflycloud" {
    auth_method = "password"
    region_name = "HN"
    version = "0.0.1"
}

resource "bizflycloud_kubernets" "test_k8s" {
    name = "tung491-test-k8s"
    version = "5f6425f3d0d3befd40e7a31f"
    auto_upgrade = false
    enable_cloud = true
    tags = ["string"]
    worker_pools = [
        {
            name = "pool"
            version = "v1.18.0"
            flavor = "8c_8g"
            profile_type = "premium"
            volume_type = "SSD"
            volume_size = 40,
            availability_zone = "HN1"
            desired_size = 1
            enable_autoscaling = true
            min_size = 1
            max_size = 3
            tags = ["string"]
        }
    ]
    }
