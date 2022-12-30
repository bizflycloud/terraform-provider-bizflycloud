terraform {
  required_providers {
    bizflycloud = {
      source  = "bizflycloud/bizflycloud"
    }
  }
}

provider "bizflycloud" {
  auth_method  = "password"
  region_name  = "HN"
  email        = var.username
  password     = var.password
}

resource "bizflycloud_cloud_database_instance" "test_instance_1" {
  name                         = "test_tera_mongo_2"
  flavor_name                  = "1c_2g"
  instance_type                = "enterprise"
  volume_size                  = 20
  datastore_type               = "MongoDB"
  datastore_version_id         = "b48a59df-7a71-49a2-838c-50c9369976bc"
  network_ids                  = ["489a1cb8-92f4-4393-9df5-3866c025bcac"]
  public_access                = true
  availability_zone            = "HN1"
  autoscaling_enable           = true
  autoscaling_volume_threshold = 90
  autoscaling_volume_limited   = 90
}
resource "bizflycloud_cloud_database_instance" "test_instance_2" {
  name                         = "test_tera_redis_2"
  flavor_name                  = "1c_2g"
  instance_type                = "enterprise"
  volume_size                  = 20
  datastore_type               = "Redis"
  datastore_version_id         = "8b3b46fa-1141-46ab-8d05-ce2720a0dcb2"
  network_ids                  = ["489a1cb8-92f4-4393-9df5-3866c025bcac"]
  public_access                = true
  availability_zone            = "HN1"
    autoscaling_enable           = true
    autoscaling_volume_threshold = 90
    autoscaling_volume_limited   = 90
}
resource "bizflycloud_cloud_database_instance" "test_instance_3" {
  name                         = "test_tera_maria_3"
  flavor_name                  = "1c_2g"
  instance_type                = "enterprise"
  volume_size                  = 10
  datastore_type               = "MariaDB"
  datastore_version_id         = "550aebf7-df97-49f1-bf24-7cd7b69fa365"
  network_ids                  = ["489a1cb8-92f4-4393-9df5-3866c025bcac"]
  public_access                = true
  availability_zone            = "HN1"
}

resource "bizflycloud_cloud_database_configuration" "test_configuration_1" {
  name = "test_configuration_1"
  parameters = {
    "auditLog.format" = "test12345434"
    "net.ipv6" = false,
    "net.maxIncomingConnections" = 345
  }
  datastore_type = "MongoDB"
  datastore_version_name = "4.4.7"
}

resource "bizflycloud_cloud_database_backup" "test_backup" {
  name = "test_backup_1"
  node_id = bizflycloud_cloud_database_instance.test_instance_1.nodes.0.id
}

resource "bizflycloud_cloud_database_schedule" "test_schedule" {
  name = "test_schedule_1"
  limit_backup = 1
  schedule_type = "monthly"
  minute = [20, 50]
  hour = [7]
  day_of_month = [10]
  node_id = bizflycloud_cloud_database_instance.test_instance_1.nodes.0.id
}

resource "bizflycloud_cloud_database_node" "test_node_1" {
  replica_of = bizflycloud_cloud_database_instance.test_instance_1.id
}
