##################################################
# DB Parameter Group
# - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/db_parameter_group
##################################################
resource "aws_db_parameter_group" "this" {
  count = var.create_db_parameter_group ? 1 : 0

  name        = var.db_parameter_group_name
  description = var.db_parameter_group_description
  family      = var.db_parameter_group_family

  dynamic "parameter" {
    for_each = var.db_parameters

    content {
      name  = parameter.value.name
      value = parameter.value.value
    }
  }

  tags = merge(
    var.tags,
    { Name = var.db_parameter_group_name },
  )
}
