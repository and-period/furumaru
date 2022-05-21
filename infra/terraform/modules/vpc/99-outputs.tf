##################################################
# VPC
##################################################
output "vpcs" {
  description = "VPC"
  value       = aws_vpc.this
}

##################################################
# DHCP Option Set
##################################################
output "dhcp_option_sets" {
  description = "DHCPオプションセット"
  value       = aws_vpc_dhcp_options.this
}

##################################################
# Route Table
##################################################
output "route_tables" {
  description = "ルートテーブル"
  value       = [for rt in aws_route_table.this : rt]
}

##################################################
# Subnet
##################################################
output "subnets" {
  description = "サブネット"
  value       = [for s in aws_subnet.this : s]
}

##################################################
# Internet Gateway
##################################################
output "internet_gateways" {
  description = "インターネットゲートウェイ"
  value       = aws_internet_gateway.this
}

##################################################
# NAT Gateway
##################################################
output "nat_gateways" {
  description = "NATゲートウェイ"
  value       = [for ngw in aws_nat_gateway.this : ngw]
}

##################################################
# VPN Gateway
##################################################
output "vpn_gateways" {
  description = "VPNゲートウェイ"
  value       = aws_vpn_gateway.this
}

##################################################
# Client VPN
##################################################
output "client_vpns" {
  description = "Client VPN"
  value       = aws_ec2_client_vpn_endpoint.this
}
