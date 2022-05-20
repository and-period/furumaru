##################################################
# Security Group
# - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/security_group
##################################################
resource "aws_security_group" "this" {
  count = 1

  name        = var.name
  description = var.description
  vpc_id      = data.aws_vpc.this.id

  dynamic "ingress" {
    for_each = local.ingress_rules

    content {
      description     = ingress.value.description
      protocol        = ingress.value.protocol
      from_port       = ingress.value.from_port
      to_port         = ingress.value.to_port
      cidr_blocks     = concat(ingress.value.cidr_blocks, [])
      prefix_list_ids = concat(ingress.value.prefix_list_ids, [])
      security_groups = concat(ingress.value.security_groups, [])
    }
  }

  dynamic "egress" {
    for_each = local.egress_rules

    content {
      description     = egress.value.description
      protocol        = egress.value.protocol
      from_port       = egress.value.from_port
      to_port         = egress.value.to_port
      cidr_blocks     = concat(egress.value.cidr_blocks, [])
      prefix_list_ids = concat(egress.value.prefix_list_ids, [])
      security_groups = concat(egress.value.security_groups, [])
    }
  }

  tags = merge(
    var.tags,
    { Name = var.name },
  )

  lifecycle {
    ignore_changes = [
      name,
      description,
      vpc_id,
    ]
  }
}
