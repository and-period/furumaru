##################################################
# Subnet Group
##################################################
output "subnet_groups" {
  description = "DB サブネットグループ"
  value       = aws_db_subnet_group.this
}

##################################################
# RDS Instance
##################################################
output "rds_instances" {
  description = "RDS インスタンス"
  value = concat(
    aws_db_instance.master,
    aws_db_instance.replica,
  )
}

##################################################
# RDS Group
##################################################
# Parameter Group
output "rds_parameter_groups" {
  description = "RDS パラメータグループ"
  value       = aws_db_parameter_group.this
}

# Option Group
output "rds_option_groups" {
  description = "RDS オプショングループ"
  value       = aws_db_option_group.this
}
