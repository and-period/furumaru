##################################################
# Elastic Container Registry Repository (Private)
# - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/ecr_lifecycle_policy
# - https://docs.aws.amazon.com/AmazonECR/latest/userguide/LifecyclePolicies.html#lifecycle_policy_parameters
##################################################
resource "aws_ecr_lifecycle_policy" "this" {
  count = local.create_lifecycle_policy ? 1 : 0

  repository = aws_ecr_repository.this[0].name
  policy     = local.lifecycle_rules
}
