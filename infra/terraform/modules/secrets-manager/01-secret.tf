##################################################
# Secrets Manager Secret
# - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/secretsmanager_secret
##################################################
resource "aws_secretsmanager_secret" "this" {
  name        = var.name
  description = var.description

  kms_key_id              = var.kms_alias_name != "" ? data.aws_kms_alias.this[0].target_key_arn : null
  recovery_window_in_days = var.recovery_window_in_days

  tags = merge(
    var.tags,
    { Name = var.name },
  )
}
