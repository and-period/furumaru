##################################################
# Virtual Private Gateway (Site-to-Site)
# - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/vpn_gateway
# - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/vpn_gateway_attachment
##################################################
resource "aws_vpn_gateway" "this" {
  count = var.create_vpn_gateway ? 1 : 0

  vpc_id = aws_vpc.this[0].id

  amazon_side_asn = var.vpn_gateway_asn != 0 ? var.vpn_gateway_asn : null

  tags = merge(
    local.tags,
    { Name = var.vpn_gateway_name },
  )

  lifecycle {
    ignore_changes = [
      amazon_side_asn,
    ]
  }
}

resource "aws_vpn_gateway_attachment" "this" {
  count = var.create_vpn_gateway ? 1 : 0

  vpc_id         = aws_vpc.this[0].id
  vpn_gateway_id = aws_vpn_gateway.this[0].id
}
