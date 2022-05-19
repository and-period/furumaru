##############################
# Subnet
# - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/subnet
##############################
resource "aws_subnet" "this" {
  for_each = local.subnets

  vpc_id = aws_vpc.this[0].id

  cidr_block                       = each.value.cidr_block
  availability_zone                = each.value.availability_zone
  map_public_ip_on_launch          = each.value.map_public_ip_on_launch
  assign_ipv6_address_on_creation  = var.assign_generated_ipv6_cidr_block

  enable_dns64                                   = var.enable_dns64
  enable_resource_name_dns_a_record_on_launch    = var.enable_resource_name_dns_a_record_on_launch
  enable_resource_name_dns_aaaa_record_on_launch = var.enable_resource_name_dns_aaaa_record_on_launch

  tags = merge(
    local.tags,
    { Name = each.key },
  )

  lifecycle {
    ignore_changes = [
      cidr_block,
      availability_zone,
    ]
  }
}
