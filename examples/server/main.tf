terraform {
    required_providers {
        bizflycloud = {
            versions = ["v0.0.1"]
            source = "bizflycloud"
        }
    }
}
provider "bizflycloud" {
    auth_method = "password"
    region_name = "HN"
    version = "v0.0.1"
}

resource "bizflycloud_server" "sapd-server" {
    name = "sapd-tf-server"
    flavor_name = "4c_2g"
    ssh_key = "sapd1"
    os_type = "image"
    os_id = "5f218529-ce32-4cb6-8557-920b16307d35"
    category = "premium"
    availability_zone = "HN1"
    root_disk_type = "HDD"
    root_disk_size = 20    
}

resource "bizflycloud_volume" "volume1" {
    name = "sapd-volume-tf3"
    size = 20
    type = "HDD"
    category = "premium"
    availability_zone = "HN2"
}