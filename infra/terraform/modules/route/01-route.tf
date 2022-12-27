##################################################
# Route
# - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/route
##################################################
resource "aws_route" "internet_gateway" {
  for_each = local.internet_gateways

  route_table_id = data.aws_route_table.this.id

  destination_cidr_block = each.value.destination_cidr_block
  gateway_id             = data.aws_internet_gateway.this[each.key].id
}

resource "aws_route" "nat_gateway" {
  for_each = local.nat_gateways

  route_table_id = data.aws_route_table.this.id

  destination_cidr_block = each.value.destination_cidr_block
  nat_gateway_id         = data.aws_nat_gateway.this[each.key].id
}

resource "aws_route" "vpn_gateway" {
  for_each = local.vpn_gateways

  route_table_id = data.aws_route_table.this.id

  destination_cidr_block = each.value.destination_cidr_block
  gateway_id             = data.aws_vpn_gateway.this[each.key].id
}

resource "aws_route" "transit_gateway" {
  for_each = local.transit_gateways

  route_table_id = data.aws_route_table.this.id

  destination_cidr_block = each.value.destination_cidr_block
  transit_gateway_id     = data.aws_ec2_transit_gateway.this[each.key].id
}

resource "aws_route" "vpc_peering" {
  for_each = local.vpc_peerings

  route_table_id = data.aws_route_table.this.id

  destination_cidr_block    = each.value.destination_cidr_block
  vpc_peering_connection_id = data.aws_vpc_peering_connection.this[each.key].id
}

resource "aws_route" "instance" {
  for_each = local.instnaces

  route_table_id = data.aws_route_table.this.id

  destination_cidr_block = each.value.destination_cidr_block
  network_interface_id   = data.aws_instance.this[each.key].network_interface_id
}

resource "aws_route" "network_interface" {
  for_each = local.network_interfaces

  route_table_id = data.aws_route_table.this.id

  destination_cidr_block = each.value.destination_cidr_block
  network_interface_id   = data.aws_network_interface.this[each.key].id
}
