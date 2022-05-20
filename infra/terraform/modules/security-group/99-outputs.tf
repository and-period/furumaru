##################################################
# Security Group
##################################################
output "security_groups" {
  description = "セキュリティグループ一覧"
  value       = aws_security_group.this
}
