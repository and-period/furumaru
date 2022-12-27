module "sg_nat" {
  source = "./../../../modules/security-group"

  #####################################################################
  # Common
  #####################################################################
  vpc_name = var.vpc_name
  tags     = var.tags

  #####################################################################
  # Security Group
  #####################################################################
  name        = "furumaru-stg-sg-nat"
  description = "for nat instance"

  ingress_rules = [
    {
      description          = "all icmp"
      protocol             = "icmp"
      from_port            = 0
      to_port              = 0
      cidr_blocks          = ["0.0.0.0/0"]
      prefix_list_ids      = []
      security_group_names = []
    },
    {
      description          = "ssh from external"
      protocol             = "tcp"
      from_port            = 22
      to_port              = 22
      cidr_blocks          = ["0.0.0.0/0"]
      prefix_list_ids      = []
      security_group_names = []
    },
    {
      description          = "all from internal"
      protocol             = "-1"
      from_port            = 0
      to_port              = 0
      cidr_blocks          = ["10.110.0.0/16"]
      prefix_list_ids      = []
      security_group_names = []
    },
  ]

  egress_rules = [
    {
      description          = "all"
      protocol             = "-1"
      from_port            = 0
      to_port              = 0
      cidr_blocks          = ["0.0.0.0/0"]
      prefix_list_ids      = []
      security_group_names = []
    },
  ]
}
