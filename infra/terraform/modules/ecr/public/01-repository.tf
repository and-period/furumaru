##################################################
# Elastic Container Registry Repository (Public)
# - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/ecrpublic_repository
##################################################
resource "aws_ecrpublic_repository" "this" {
  count = 1

  repository_name = var.repository_name

  catalog_data {
    description       = var.description
    operating_systems = concat(var.operating_systems, [])
    architectures     = concat(var.architectures, [])
    about_text        = var.about_text
    usage_text        = var.usage_text
  }
}
