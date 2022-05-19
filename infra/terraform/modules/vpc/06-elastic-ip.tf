##################################################
# Elastic IP
# - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/eip
##################################################
resource "aws_eip" "nat_gateway" {
  # Internet Gatewayがない状態ではNAT Gatewayは作成できないため
  for_each = var.create_internet_gateway ? local.nat_gateways : {}

  vpc = true

  tags = merge(
    local.tags,
    { Name = each.value.eip_name },
  )

  depends_on = [
    aws_internet_gateway.this,
  ]
}
