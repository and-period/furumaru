module "sg_app" {
  source = "./../../../modules/security-group"

  #####################################################################
  # Common
  #####################################################################
  vpc_name = var.vpc_name
  tags     = var.tags

  #####################################################################
  # Security Group
  #####################################################################
  name        = "furumaru-stg-sg-app"
  description = "for database"

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
      description          = "gateway from external"
      protocol             = "tcp"
      from_port            = 9000
      to_port              = 9000
      cidr_blocks          = ["0.0.0.0/0"]
      prefix_list_ids      = []
      security_group_names = []
    },
    {
      description          = "metrics from external"
      protocol             = "tcp"
      from_port            = 9001
      to_port              = 9001
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
