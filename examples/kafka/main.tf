terraform {
  required_providers {
    bizflycloud = {
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

# variable "PROJECT_ID" {
#   type=string
# }

provider "bizflycloud" {
  auth_method = "password"
  email       = var.EMAIL
  password    = var.PASSWORD
#  project_id  = var.PROJECT_ID
  region_name = "HaNoi"
}

## Get version of the kafka
data "bizflycloud_kafka_version" "tf_kafka_version" {
}

## List available kafka version
output "available_BizflyCloud_kafka_versions" {
  description = "List of all available Kafka versions"
  value = [
    for v in data.bizflycloud_kafka_version.tf_kafka_version.versions : {
      id         = v.id
      name       = v.name
      code       = v.code
      is_default = v.is_default
    }
  ]
}

## Get available flavors
data "bizflycloud_kafka_flavor" "tf_kafka_flavor" {
}

## List available kafka flavors
output "available_BizflyCloud_kafka_flavors" {
  description = "List of all available Kafka flavors"
  value = [
    for f in data.bizflycloud_kafka_flavor.tf_kafka_flavor.flavors : {
      id          = f.id
      name        = f.name
      vcpus       = f.vcpus
      ram_mb      = f.ram
      disk_gb     = f.disk
      is_default  = f.is_default
      description = f.description
    }
  ]
}

## Create cluster
resource "bizflycloud_kafka" "tf-kafka1" {
  name              = "tf-kafka1"
  version_id        = "10d9c0ee-5b3b-4a9e-9f21-4416fbb9e8ef"
  nodes             = 1
  flavor            = "2c_4g"
  volume_size       = 20
  availability_zone = "HN2"
  vpc_network_id    = "6046aaac-1bf4-4c1b-bbf3-1eddedaf6ee3"
  public_access     = false
}

data "bizflycloud_kafka" "tf-kafka1" {
  name = bizflycloud_kafka.tf-kafka1.name
}

output "kafka_create_cluster_info" {
  value = {
    name = data.bizflycloud_kafka.tf-kafka1.name
    version_id = data.bizflycloud_kafka.tf-kafka1.version_id
    nodes  = data.bizflycloud_kafka.tf-kafka1.nodes
    volume_size = data.bizflycloud_kafka.tf-kafka1.volume_size
    flavor  = data.bizflycloud_kafka.tf-kafka1.flavor
    availability_zone = data.bizflycloud_kafka.tf-kafka1.availability_zone
    public_access = data.bizflycloud_kafka.tf-kafka1.public_access
    status = data.bizflycloud_kafka.tf-kafka1.status
  }
}

## Get cluster info by ID
data "bizflycloud_kafka" "by_id" {
  id = "cfcec8cb-76c7-4e8c-8792-15ab0103a5bb"
}

output "kafka_cluster_info_by_id" {
  description = "Get cluster info using cluster ID"
  value = {
    id                = data.bizflycloud_kafka.by_id.id
    name              = data.bizflycloud_kafka.by_id.name
    version_id        = data.bizflycloud_kafka.by_id.version_id
    nodes             = data.bizflycloud_kafka.by_id.nodes
    volume_size       = data.bizflycloud_kafka.by_id.volume_size
    flavor            = data.bizflycloud_kafka.by_id.flavor
    availability_zone = data.bizflycloud_kafka.by_id.availability_zone
    vpc_network_id    = data.bizflycloud_kafka.by_id.vpc_network_id
    public_access     = data.bizflycloud_kafka.by_id.public_access
    status            = data.bizflycloud_kafka.by_id.status
  }
}

## Add node to existing cluster
resource "bizflycloud_kafka" "tf-kafka1" {
  name              = "tf-kafka1"
  version_id        = "10d9c0ee-5b3b-4a9e-9f21-4416fbb9e8ef"
  
  # Increase nodes count to scale out (e.g., from 1 to 2)
  # NOTE: Nodes can only be increased, not decreased
  nodes             = 3  # <- Increase this value to add more nodes
  
  flavor            = "2c_4g"
  volume_size       = 20
  availability_zone = "HN2"
  vpc_network_id    = "6046aaac-1bf4-4c1b-bbf3-1eddedaf6ee3"
  public_access     = false
}

# Important notes:
# - Node count can only be INCREASED, not decreased
# - Minimum node count is 1
# - Adding nodes improves cluster capacity and availability
# - The cluster will remain available during scale out

# After changing the nodes count, run:
# terraform plan    # to see what will change
# terraform apply   # to apply the change

# Check the status
data "bizflycloud_kafka" "tf-kafka1" {
  name = bizflycloud_kafka.tf-kafka1.name
}

output "cluster_scale_node_info" {
  value = {
    name   = data.bizflycloud_kafka.tf-kafka1.name
    nodes  = data.bizflycloud_kafka.tf-kafka1.nodes
    status = data.bizflycloud_kafka.tf-kafka1.status
  }
}


## Resize volume cluster from existing cluster
resource "bizflycloud_kafka" "tf-kafka1" {
  name              = "tf-kafka1"
  version_id        = "10d9c0ee-5b3b-4a9e-9f21-4416fbb9e8ef"
  nodes             = 2
  flavor            = "2c_4g"
  
  # Increase volume_size (e.g., from 10GB to 30GB)
  # NOTE: Volume can only be increased, not decreased
  volume_size       = 30  # <- Increase this value (minimum 10GB)
  
  availability_zone = "HN2"
  vpc_network_id    = "6046aaac-1bf4-4c1b-bbf3-1eddedaf6ee3"
  public_access     = false
}

# Important notes:
# - Volume size can only be INCREASED, not decreased
# - Minimum volume size is 10GB
# - Changes will trigger a resize operation
# - The cluster will remain available during resize

# After changing the volume_size, run:
# terraform plan    # to see what will change
# terraform apply   # to apply the change

# Check the status
data "bizflycloud_kafka" "tf-kafka1" {
  name = bizflycloud_kafka.tf-kafka1.name
}

output "cluster_resize_volume_info" {
  value = {
    name        = data.bizflycloud_kafka.tf-kafka1.name
    volume_size = data.bizflycloud_kafka.tf-kafka1.volume_size
    status      = data.bizflycloud_kafka.tf-kafka1.status
  }
}


## Resize flavor cluster from existing cluster
resource "bizflycloud_kafka" "tf-kafka1" {
  name              = "tf-kafka1"
  version_id        = "10d9c0ee-5b3b-4a9e-9f21-4416fbb9e8ef"
  nodes             = 2
  
  # Change flavor to a new one (e.g., from "2c_4g" to "2c_8g")
  flavor            = "2c_8g"  # <- Change this to upgrade/downgrade
  
  volume_size       = 30
  availability_zone = "HN2"
  vpc_network_id    = "6046aaac-1bf4-4c1b-bbf3-1eddedaf6ee3"
  public_access     = false
}

# After changing the flavor, run:
# terraform plan    # to see what will change
# terraform apply   # to apply the change

# Check the status
data "bizflycloud_kafka" "tf-kafka1" {
  name = bizflycloud_kafka.tf-kafka1.name
}

output "cluster_resize_flavor_info" {
  value = {
    name   = data.bizflycloud_kafka.tf-kafka1.name
    flavor  = data.bizflycloud_kafka.tf-kafka1.flavor
    status = data.bizflycloud_kafka.tf-kafka1.status
  }
}


# If you want to delete the cluster, comment out or remove the resource block.
resource "bizflycloud_kafka" "tf-kafka1" {
  name              = "tf-kafka1"
  version_id        = "10d9c0ee-5b3b-4a9e-9f21-4416fbb9e8ef"
  nodes             = 2
  flavor            = "2c_8g"
  volume_size       = 30
  availability_zone = "HN2"
  vpc_network_id    = "6046aaac-1bf4-4c1b-bbf3-1eddedaf6ee3"
  public_access     = true
}

# Or use terraform destroy to delete all resources:
# terraform destroy

# To delete only the Kafka cluster (using cluster ID from state):
# terraform destroy -target=bizflycloud_kafka.tf-kafka1