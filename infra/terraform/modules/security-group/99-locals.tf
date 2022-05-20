locals {
  // ID指定が必要なセキュリティグループ名一覧
  security_group_names = distinct(flatten(
    concat(
      [for ir in var.ingress_rules : [ir.security_group_names]],
      [for er in var.egress_rules : [er.security_group_names]],
    )
  ))

  security_groups = {
    for sg in local.security_group_names : sg => sg
  }

  // インバウンドルール
  ingress_rules = flatten([
    for ir in var.ingress_rules : {
      description     = ir.description
      protocol        = ir.protocol
      from_port       = ir.from_port
      to_port         = ir.to_port
      cidr_blocks     = ir.cidr_blocks
      prefix_list_ids = ir.prefix_list_ids
      security_groups = flatten([
        for n in ir.security_group_names : [
          data.aws_security_group.this[n].id,
        ]
      ])
    }
  ])

  // アウトバウンドルール
  egress_rules = flatten([
    for er in var.egress_rules : {
      description     = er.description
      protocol        = er.protocol
      from_port       = er.from_port
      to_port         = er.to_port
      cidr_blocks     = er.cidr_blocks
      prefix_list_ids = er.prefix_list_ids
      security_groups = flatten([
        for n in er.security_group_names : [
          data.aws_security_group.this[n].id,
        ]
      ])
    }
  ])
}
