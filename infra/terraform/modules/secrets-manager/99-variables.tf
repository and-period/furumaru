##################################################
# Common
##################################################
variable "tags" {
  default = {}
}

variable "kms_alias_name" {
  default = ""
}

##################################################
# Secrets Manager
##################################################
variable "name" {
  description = "シークレット名"
  default     = ""
}

variable "description" {
  description = "シークレットの説明"
  default     = ""
}

variable "recovery_window_in_days" {
  description = "削除後のリカバリ日数"
  type        = number
  default     = 30
}
