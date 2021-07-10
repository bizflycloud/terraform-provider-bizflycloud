output "status" {
  value = bizflycloud_vpc_network.vpc_network.status
}

output "created_at" {
  value = bizflycloud_vpc_network.vpc_network.created_at
}

output "updated_at" {
  value = bizflycloud_vpc_network.vpc_network.updated_at
}

output "availability_zones" {
  value = bizflycloud_vpc_network.vpc_network.availability_zones
}

output "mtu" {
  value = bizflycloud_vpc_network.vpc_network.mtu
}

output "tags" {
  value = bizflycloud_vpc_network.vpc_network.tags
}


output "subnets" {
  value = bizflycloud_vpc_network.vpc_network.subnets
}
