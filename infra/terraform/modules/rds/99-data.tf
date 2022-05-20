##################################################
# ネットワークリソースの取得 (タグ指定)
# - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/vpc
# - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/subnet
# - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/security_group
# - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/kms_alias
##################################################
data "aws_vpc" "this" {
  filter {
    name = "tag:Name"
    values = [
      var.vpc_name
    ]
  }
}

# タグ:Name の値より、Subnetの情報を取得
data "aws_subnet" "this" {
  count = length(var.subnet_names)

  vpc_id = data.aws_vpc.this.id

  filter {
    name = "tag:Name"
    values = [
      element(var.subnet_names, count.index)
    ]
  }
}

# タグ:Name の値より、Security Groupの情報を取得
data "aws_security_group" "this" {
  count = length(var.security_group_names)

  vpc_id = data.aws_vpc.this.id

  filter {
    name = "tag:Name"
    values = [
      element(var.security_group_names, count.index),
    ]
  }
}

# KMS Alias名より、KMSキー情報を取得
data "aws_kms_alias" "this" {
  count = var.kms_alias_name != "" ? 1 : 0

  name = format("alias/%s", var.kms_alias_name)
}
