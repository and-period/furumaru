##################################################
# Route Table
# - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/route_table
# - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/route_table_association
##################################################
resource "aws_route_table" "this" {
  for_each = local.route_tables

  vpc_id = aws_vpc.this[0].id

  tags = merge(
    local.tags,
    { Name = each.key },
  )
}

resource "aws_route_table_association" "this" {
  for_each = local.subnets

  route_table_id = aws_route_table.this[each.value.route_table_name].id
  subnet_id      = aws_subnet.this[each.key].id
}
