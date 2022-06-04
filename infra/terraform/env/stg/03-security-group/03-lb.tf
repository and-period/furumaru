module "sg_lb" {
  source = "./../../../modules/security-group"

  #####################################################################
  # Common
  #####################################################################
  vpc_name = var.vpc_name
  tags     = var.tags

  #####################################################################
  # Security Group
  #####################################################################
  name        = "furumaru-stg-sg-lb"
  description = "for elb"

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
      description          = "http from external"
      protocol             = "tcp"
      from_port            = 80
      to_port              = 80
      cidr_blocks          = ["0.0.0.0/0"]
      prefix_list_ids      = []
      security_group_names = []
    },
    {
      description          = "https from external"
      protocol             = "tcp"
      from_port            = 443
      to_port              = 443
      cidr_blocks          = ["0.0.0.0/0"]
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
