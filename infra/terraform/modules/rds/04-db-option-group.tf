##################################################
# DB Option Group
# - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/db_option_group
##################################################
resource "aws_db_option_group" "this" {
  count = var.create_db_option_group ? 1 : 0

  name                     = var.db_option_group_name
  option_group_description = var.db_option_group_description

  engine_name          = var.db_engine
  major_engine_version = var.db_option_group_major_engine_version

  dynamic "option" {
    for_each = var.db_options

    content {
      option_name = option.value.name

      dynamic "option_settings" {
        for_each = option.value.settings

        content {
          name  = option_settings.value.key
          value = option_settings.value.value
        }
      }
    }
  }

  tags = merge(
    var.tags,
    { Name = var.db_option_group_name },
  )
}
