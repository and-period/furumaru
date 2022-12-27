##################################################
# KMS
# - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/kms_alias
##################################################
data "aws_kms_alias" "this" {
  count = var.kms_alias_name != "" ? 1 : 0

  name = format("alias/%s", var.kms_alias_name)
}
