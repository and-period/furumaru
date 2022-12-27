##################################################
# AMI
# - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/ami
##################################################
# 使用するAMIを検索 (以下の情報から検索実施)
# - owners:      AMIのオーナーリスト
# - name:        AMIの名前に含まれるキーワード
# - most_recent: AMIを複数取得した場合、取得した中で最新のAMIを使用
data "aws_ami" "this" {
  count = var.ami_name != "" ? 1 : 0

  most_recent = var.ami_most_recent
  owners      = concat(var.ami_owners, [])

  filter {
    name = "name"
    values = [
      var.ami_name
    ]
  }
}

##################################################
# KMS
# - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/kms_alias
##################################################
data "aws_kms_alias" "this" {
  count = var.kms_alias_name != "" ? 1 : 0

  name = format("alias/%s", var.kms_alias_name)
}

##################################################
# ネットワークリソースの取得 (タグ指定)
# - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/vpc
# - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/subnet
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

# タグ:Name の値より、Subnetの情報を取得
data "aws_subnet" "this" {
  vpc_id = data.aws_vpc.this.id

  filter {
    name = "tag:Name"
    values = [
      var.subnet_name
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
