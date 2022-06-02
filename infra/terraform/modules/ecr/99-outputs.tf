##################################################
# ECR Repository (Public)
##################################################
output "ecrpublic_repository" {
  description = "ECR パブリックリポジトリ"
  value       = aws_ecrpublic_repository.this
}
