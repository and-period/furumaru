##################################################
# ネットワークリソースの取得 (タグ指定)
# - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/vpc
# - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/route_table
# - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/internet_gateway
# - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/nat_gateway
# - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/vpn_gateway
# - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/ec2_transit_gateway
# - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/vpc_peering_connection
# - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/instance
# - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/network_interface
##################################################
data "aws_vpc" "this" {
  filter {
    name = "tag:Name"
    values = [var.vpc_name]
  }
}

data "aws_route_table" "this" {
  vpc_id = data.aws_vpc.this.id

  filter {
    name = "tag:Name"
    values = [var.route_table_name]
  }
}

data "aws_internet_gateway" "this" {
  for_each = local.internet_gateways

  filter {
    name   = "tag:Name"
    values = [each.value.target_resource_name]
  }
}

data "aws_nat_gateway" "this" {
  for_each = local.nat_gateways

  vpc_id = data.aws_vpc.this.id

  filter {
    name   = "tag:Name"
    values = [each.value.target_resource_name]
  }
}

data "aws_vpn_gateway" "this" {
  for_each = local.vpn_gateways

  attached_vpc_id = data.aws_vpc.this.id

  filter {
    name   = "tag:Name"
    values = [each.value.target_resource_name]
  }
}

data "aws_ec2_transit_gateway" "this" {
  for_each = local.transit_gateways

  id = each.value.target_resource_name
}

data "aws_vpc_peering_connection" "this" {
  for_each = local.vpc_peerings

  id = each.value.target_resource_name
}

data "aws_instance" "this" {
  for_each = local.instnaces

  filter {
    name   = "tag:Name"
    values = [each.value.target_resource_name]
  }
}

data "aws_network_interface" "this" {
  for_each = local.network_interfaces

  filter {
    name   = "tag:Name"
    values = [each.value.target_resource_name]
  }
}
