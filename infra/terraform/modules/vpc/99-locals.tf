locals {
  ## Tags
  # 共通のタグ
  tags = merge(
    var.tags,
    { Name = var.name },
  )

  # Route Tables
  route_tables = {
    for rt in var.route_tables : rt => rt
  }

  # Subnets
  subnets = {
    for s in var.subnets : s.name => s
  }

  # NAT Gateways
  nat_gateways = {
    for ngw in var.nat_gateways : ngw.name => ngw
  }
}
