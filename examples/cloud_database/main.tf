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
  email       = var.username
  password    = var.password
}


# Get information about backup
data "bizflycloud_cloud_database_datastore" "ds" {
  type = "Postgres"
  name = "patroni-13.11"
}

# Get information about backup
data "bizflycloud_cloud_database_backup" "init" {
  id = "7201f1dd-6185-46a2-b488-87f2a4fce70e"
}

# Restore a backup - Create a cloud database instance from a backup and secondary
resource "bizflycloud_cloud_database_instance" "terraform_patroni_mongo" {
  name = "mongo"
  autoscaling = {
    enable           = 1
    volume_limited   = 100
    volume_threshold = 90
  }
  backup_id = "7201f1dd-6185-46a2-b488-87f2a4fce70e"
  datastore = {
    type       = data.bizflycloud_cloud_database_backup.init.datastore.type
    name       = data.bizflycloud_cloud_database_backup.init.datastore.name
    version_id = data.bizflycloud_cloud_database_backup.init.datastore.version_id
  }
  availability_zone = "HN1"
  flavor_name       = "2c_4g"
  instance_type     = "enterprise"

  secondaries {
    availability_zone = "HN1"
    quantity          = 2
  }

  network_ids = [
    "0706f928-02f6-4c4a-935c-fc22fdccbe81"
  ]

  init_databases = ["games", "genshin", "impact"]
  users {
    username  = "yanfei"
    password  = "password"
    databases = ["games", "genshin", "impact"]
  }

  users {
    username  = "nahida"
    password  = "password"
    databases = ["games", "genshin", "impact"]
  }

  users {
    username  = "nilou"
    password  = "password"
    databases = ["games", "genshin", "impact"]
  }

  public_access = true
  volume_size   = 20
}

