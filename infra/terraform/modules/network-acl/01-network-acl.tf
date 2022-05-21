##################################################
# Network ACL
# - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/network_acl
##################################################
resource "aws_network_acl" "this" {
  count = 1

  vpc_id     = data.aws_vpc.this.id
  subnet_ids = concat(data.aws_subnet.this[*].id, [])

  dynamic "ingress" {
    for_each = var.ingress_rules

    content {
      rule_no    = ingress.value.rule_no
      action     = ingress.value.action
      protocol   = ingress.value.protocol
      from_port  = ingress.value.from_port
      to_port    = ingress.value.to_port
      cidr_block = ingress.value.cidr_block
    }
  }

  dynamic "egress" {
    for_each = var.egress_rules

    content {
      rule_no    = egress.value.rule_no
      action     = egress.value.action
      protocol   = egress.value.protocol
      from_port  = egress.value.from_port
      to_port    = egress.value.to_port
      cidr_block = egress.value.cidr_block
    }
  }

  tags = merge(
    var.tags,
    { Name = var.name },
  )
}
