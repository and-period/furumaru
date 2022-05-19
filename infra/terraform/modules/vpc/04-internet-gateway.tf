##################################################
# Internet Gateway
# - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/internet_gateway
##################################################
resource "aws_internet_gateway" "this" {
  count = var.create_internet_gateway ? 1 : 0

  vpc_id = aws_vpc.this[0].id

  tags = merge(
    local.tags,
    { Name = var.internet_gateway_name },
  )
}
