##################################################
# DHCP Options
# - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/vpc_dhcp_options
# - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/vpc_dhcp_options_association
##################################################
resource "aws_vpc_dhcp_options" "this" {
  count = var.create_dhcp_options ? 1 : 0

  domain_name          = var.domain_name
  domain_name_servers  = concat(var.domain_name_servers, [])
  ntp_servers          = concat(var.ntp_servers, [])
  netbios_name_servers = concat(var.netbios_name_servers, [])
  netbios_node_type    = var.netbios_node_type

  tags = merge(
    local.tags,
    { Name = var.dhcp_options_name },
  )
}

resource "aws_vpc_dhcp_options_association" "this" {
  count = var.create_dhcp_options ? 1 : 0

  vpc_id          = aws_vpc.this[0].id
  dhcp_options_id = aws_vpc_dhcp_options.this[0].id
}
