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
  email       = "example@domain.com"
  password    = "thisispassword"
}

resource "bizflycloud_autoscaling_launch_configuration" "basic-centos-terrafrom" {
  name              = "basic-centos-terrafrom"
  ssh_key           = "ministry"
  availability_zone = "HN1"
  flavor            = "1c_1g"
  instance_type     = "basic"

  networks {
    network_id = "network_id"
  }

  os {
    uuid        = "4cdbe57f-6ba1-4f40-a6fb-beb1ed974168"
    create_from = "image"

  }
  rootdisk {
    volume_type = "BASIC_SSD1"
    volume_size = 20
  }
  data_disks {
    volume_type = "BASIC_SSD1"
    volume_size = 20
  }
  data_disks {
    volume_type           = "BASIC_SSD1"
    volume_size           = 40
    delete_on_termination = true
  }
  user_data = "#!/bin/sh \n echo \"Hello World\" > /tmp/greeting.txt"
}

resource "bizflycloud_autoscaling_group" "hutao" {
  name                    = "hutao"
  launch_configuration_id = bizflycloud_autoscaling_launch_configuration.basic-centos-terrafrom.id
  max_size                = 2
  min_size                = 1
  desired_capacity        = 1
}

resource "bizflycloud_autoscaling_scalein_policy" "name" {
  cluster_id  = bizflycloud_autoscaling_group.hutao.id
  metric_type = "ram_used"
  threshold   = 10
  range_time  = 600
  cooldown    = 600
}

resource "bizflycloud_autoscaling_scaleout_policy" "name" {
  cluster_id  = bizflycloud_autoscaling_group.hutao.id
  metric_type = "ram_used"
  threshold   = 90
  range_time  = 600
  cooldown    = 600
}

resource "bizflycloud_autoscaling_scaleout_policy" "namee" {
  cluster_id  = bizflycloud_autoscaling_group.hutao.id
  metric_type = "cpu_used"
  threshold   = 80
  range_time  = 600
  cooldown    = 600
}

resource "bizflycloud_autoscaling_deletion_policy" "deletion_policy" {
  cluster_id = bizflycloud_autoscaling_group.hutao.id
  criteria   = "YOUNGEST_FIRST"
}
