##################################################
# NAT Gateway
# - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/nat_gateway
##################################################
resource "aws_nat_gateway" "this" {
  # Internet Gatewayがない状態ではNAT Gatewayは作成できないため
  for_each = var.create_internet_gateway ? local.nat_gateways : {}

  allocation_id = aws_eip.nat_gateway[each.key].id
  subnet_id     = aws_subnet.this[each.value.subnet_name].id

  tags = merge(
    local.tags,
    { Name = each.key },
  )

  depends_on = [
    aws_internet_gateway.this,
    aws_eip.nat_gateway,
  ]

  lifecycle {
    ignore_changes = [
      allocation_id,
      subnet_id,
    ]
  }
}
