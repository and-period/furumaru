module "sg_db" {
  source = "./../../../modules/security-group"

  #####################################################################
  # Common
  #####################################################################
  vpc_name = var.vpc_name
  tags     = var.tags

  #####################################################################
  # Security Group
  #####################################################################
  name        = "furumaru-stg-sg-db"
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
      description          = "mysql from app"
      protocol             = "tcp"
      from_port            = 3306
      to_port              = 3306
      cidr_blocks          = ["10.110.10.0/24", "10.110.11.0/24"]
      prefix_list_ids      = []
      security_group_names = []
    },
    {
      description          = "mysql from debug instance"
      protocol             = "tcp"
      from_port            = 3306
      to_port              = 3306
      cidr_blocks          = ["10.110.0.0/24"]
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
