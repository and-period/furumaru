##################################################
# Elastic Container Registry Repository (Private)
# - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/ecr_repository
##################################################
resource "aws_ecr_repository" "this" {
  count = 1

  name                 = var.repository_name
  image_tag_mutability = var.image_tag_mutability

  image_scanning_configuration {
    scan_on_push = var.image_scan_on_push
  }

  tags = merge(
    var.tags,
    { Name = var.repository_name },
  )
}
