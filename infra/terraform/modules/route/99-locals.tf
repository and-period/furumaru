locals {
  # Internet Gateways
  internet_gateway_routes = flatten([
    for r in var.routes : r.target_resource == "internet_gateway" ? [r] : []
  ])

  internet_gateways = {
    for ir in local.internet_gateway_routes : join("_", [
      ir.target_resource_name,
      replace(ir.destination_cidr_block, "/[.|/]/", "_"),
    ]) => ir
  }

  # NAT Gateways
  nat_gateway_routes = flatten([
    for r in var.routes : r.target_resource == "nat_gateway" ? [r] : []
  ])

  nat_gateways = {
    for nr in local.nat_gateway_routes : join("_", [
      nr.target_resource_name,
      replace(nr.destination_cidr_block, "/[.|/]/", "_"),
    ]) => nr
  }

  # VPN Gateways
  vpn_gateway_routes = flatten([
    for r in var.routes : r.target_resource == "vpn_gateway" ? [r] : []
  ])

  vpn_gateways = {
    for vr in local.vpn_gateway_routes : join("_", [
      vr.target_resource_name,
      replace(vr.destination_cidr_block, "/[.|/]/", "_"),
    ]) => vr
  }

  # Transit Gateways
  transit_gateway_routes = flatten([
    for r in var.routes : r.target_resource == "transit_gateway" ? [r] : []
  ])

  transit_gateways = {
    for tr in local.transit_gateway_routes : join("_", [
      tr.target_resource_name,
      replace(tr.destination_cidr_block, "/[.|/]/", "_"),
    ]) => tr
  }

  # VPC Peerings
  vpc_peering_routes = flatten([
    for r in var.routes : r.target_resource == "vpc_peering" ? [r] : []
  ])

  vpc_peerings = {
    for vr in local.vpc_peering_routes : join("_", [
      vr.target_resource_name,
      replace(vr.destination_cidr_block, "/[.|/]/", "_"),
    ]) => vr
  }

  # EC2 Instances
  instance_routes = flatten([
    for r in var.routes : r.target_resource == "instance" ? [r] : []
  ])

  instnaces = {
    for ir in local.instance_routes : join("_", [
      ir.target_resource_name,
      replace(ir.destination_cidr_block, "/[.|/]/", "_"),
    ]) => ir
  }

  # Network Interface
  network_interface_routes = flatten([
    for r in var.routes : r.target_resource == "network_interface" ? [r] : []
  ])

  network_interfaces = {
    for nr in local.network_interface_routes : join("_", [
      nr.target_resource_name,
      replace(nr.destination_cidr_block, "/[.|/]/", "_"),
    ]) => nr
  }
}
