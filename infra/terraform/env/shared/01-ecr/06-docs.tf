module "ecr_docs" {
  source = "./../../../modules/ecr/private"

  #####################################################################
  # ECR Repository (Private)
  #####################################################################
  repository_name = "furumaru/docs"

  image_tag_mutability = "MUTABLE"
  image_scan_on_push   = false

  lifecycle_policy_untagged_priority       = 1
  lifecycle_policy_untagged_description    = "タグ付けなしイメージの定期削除用"
  lifecycle_policy_untagged_retention_days = 7

  lifecycle_policy_priority    = 11
  lifecycle_policy_description = "イメージの世代管理 (5世代まで)"
  lifecycle_policy_image_count = 5
}
