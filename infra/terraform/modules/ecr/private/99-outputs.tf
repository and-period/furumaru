##################################################
# ECR Repository (Public)
##################################################
output "ecr_repository" {
  description = "ECR プライベートリポジトリ"
  value       = aws_ecr_repository.this
}
