##################################################
# Common
##################################################

##################################################
# ECR Repository
##################################################
variable "repository_name" {
  description = "ECR レポジトリ名"
  default     = ""
}

variable "description" {
  description = "ECR 簡単な説明(マークダウン形式)"
  default     = ""
}

variable "operating_systems" {
  description = "ECR オペレーティングシステム一覧"
  type        = list(string)
  default     = [] # Linux, Windows
}

variable "architectures" {
  description = "ECR アーキテクチャ一覧"
  type        = list(string)
  default     = [] # ARM, ARM 64, x86, x86-64
}

variable "about_text" {
  description = "ECR 詳細な説明(マークダウン形式)"
  default     = ""
}

variable "usage_text" {
  description = "ECR イメージの使用方法に関する詳細情報(マークダウン形式)"
  default     = ""
}
