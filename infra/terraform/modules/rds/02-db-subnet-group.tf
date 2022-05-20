##################################################
# DB Subnet Group
# - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/db_subnet_group
##################################################
resource "aws_db_subnet_group" "this" {
  count = 1

  name        = var.subnet_group_name
  description = var.subnet_group_description

  subnet_ids = flatten(data.aws_subnet.this[*].id)

  tags = merge(
    var.tags,
    { Name = var.subnet_group_name },
  )
}
