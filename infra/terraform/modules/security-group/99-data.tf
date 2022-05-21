##################################################
# ネットワークリソースの取得 (タグ指定)
# - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/vpc
# - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/security_group
##################################################
data "aws_vpc" "this" {
  filter {
    name = "tag:Name"
    values = [
      var.vpc_name
    ]
  }
}

data "aws_security_group" "this" {
  for_each = local.security_groups

  filter {
    name = "tag:Name"
    values = [
      each.key
    ]
  }

  vpc_id = data.aws_vpc.this.id
}
