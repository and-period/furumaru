##################################################
# Common
##################################################
variable "tags" {
  default = {}
}

##################################################
# ECR Repository
##################################################
variable "repository_name" {
  description = "ECR レポジトリ名"
  default     = ""
}

variable "image_tag_mutability" {
  description = "タグのイミュータビリティ"
  default     = "IMUTABLE" # MUTABLE, IMMUTABLE
}

variable "image_scan_on_push" {
  description = "プッシュ時のスキャン設定"
  type        = bool
  default     = false
}

##################################################
# ECR Repository - ライフサイクルポリシー
##################################################
variable "lifecycle_policy_priority" {
  description = "ライフサイクルポリシー(対象:すべて) ルールの優先度"
  type        = number
  default     = 0
}

variable "lifecycle_policy_description" {
  description = "ライフサイクルポリシー(対象:すべて) ルールの説明"
  default     = ""
}

variable "lifecycle_policy_image_count" {
  description = "ライフサイクルポリシー(対象:すべて) 保持するイメージ数"
  type        = number
  default     = 10
}

variable "lifecycle_policy_untagged_priority" {
  description = "ライフサイクルポリシー(対象:タグ付けなし) ルールの優先度"
  type        = number
  default     = 0
}

variable "lifecycle_policy_untagged_description" {
  description = "ライフサイクルポリシー(対象:タグ付けなし) ルールの説明"
  default     = ""
}

variable "lifecycle_policy_untagged_retention_days" {
  description = "ライフサイクルポリシー(対象:タグ付けなし) 保持する期間(days)"
  type        = number
  default     = 14
}
