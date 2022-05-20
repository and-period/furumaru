##################################################
# Route Table
##################################################
output "route_tables" {
  description = "ルーティング"
  value = [{
    id = data.aws_route_table.this.id,
    routes = flatten(
      concat(
        [for igw in aws_route.internet_gateway : igw],
        [for ngw in aws_route.nat_gateway : ngw],
        [for vgw in aws_route.vpn_gateway : vgw],
        [for tgw in aws_route.transit_gateway : tgw],
        [for vp in aws_route.vpc_peering : vp],
        [for i in aws_route.instance : i],
        [for ni in aws_route.network_interface : ni],
      )
    )
  }]
}
