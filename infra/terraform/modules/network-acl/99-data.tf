##################################################
# ネットワークリソースの取得 (タグ指定)
# - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/vpc
# - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/subnet
##################################################
data "aws_vpc" "this" {
  filter {
    name = "tag:Name"
    values = [
      var.vpc_name
    ]
  }
}

data "aws_subnet" "this" {
  count = length(var.subnet_names)

  vpc_id = data.aws_vpc.this.id

  filter {
    name = "tag:Name"
    values = [
      element(var.subnet_names, count.index),
    ]
  }
}
