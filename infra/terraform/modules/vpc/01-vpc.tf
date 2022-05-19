##################################################
# VPC
# - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/vpc
##################################################
resource "aws_vpc" "this" {
  count = 1

  cidr_block                       = var.cidr_block
  assign_generated_ipv6_cidr_block = var.assign_generated_ipv6_cidr_block
  instance_tenancy                 = var.instance_tenancy
  enable_dns_support               = var.enable_dns_support
  enable_dns_hostnames             = var.enable_dns_hostnames

  tags = local.tags

  lifecycle {
    ignore_changes = [
      cidr_block,
      assign_generated_ipv6_cidr_block,
      instance_tenancy,
    ]
  }
}
