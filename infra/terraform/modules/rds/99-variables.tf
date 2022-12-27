##################################################
# Common
##################################################
variable "tags" {
  default = {}
}

variable "vpc_name" {
  default = ""
}

variable "security_group_names" {
  type    = list(string)
  default = []
}

variable "kms_alias_name" {
  default = ""
}

##################################################
# Subnet Group
##################################################
variable "subnet_group_name" {
  description = "サブネットグループ名"
  default     = ""
}

variable "subnet_group_description" {
  description = "サブネットグループ説明"
  default     = ""
}

variable "subnet_names" {
  description = "サブネットグループに追加するサブネット名リスト"
  type        = list(string)
  default     = []
}

##################################################
# Parameter Group
##################################################
variable "create_db_parameter_group" {
  description = "パラメータグループの作成"
  type        = bool
  default     = false
}

variable "db_parameter_group_name" {
  description = "パラメータグループ名"
  default     = ""
}

variable "db_parameter_group_description" {
  description = "パラメータグループ説明"
  default     = ""
}

variable "db_parameter_group_family" {
  description = "パラメータグループファミリー"
  default     = "" # e.g.) mysql5.7
}

variable "db_parameters" {
  description = "パラメータグループ追加パラメータ"
  type = list(object({
    name  = string # パラメータ名
    value = string # パラメータ値
  }))
  default = []
}

##################################################
# Option Group
##################################################
variable "create_db_option_group" {
  description = "オプショングループの作成"
  type        = bool
  default     = false
}

variable "db_option_group_name" {
  description = "オプショングループ名"
  default     = ""
}

variable "db_option_group_description" {
  description = "オプショングループ説明"
  default     = ""
}

variable "db_option_group_major_engine_version" {
  description = "メジャーエンジンのバージョン"
  default     = "" # e.g.) mysql -> 5.7
}

variable "db_options" {
  description = "オプショングループ追加オプション"
  type = list(object({
    name = string # オプション名
    port = number # ポート (必要ないときは0をいれる)
    setting = list(object({
      name  = string # オプション設定
      value = string # オプション値
    }))
  }))
  default = []
}

##################################################
# RDS Instance (Master)
##################################################
# 基本設定
variable "db_instance_count" {
  description = "マスターDBの作成数"
  type        = number
  default     = 1 # 基本変えない、0にしたいときぐらい
}

variable "db_identifier" {
  description = "マスターDBインスタンス識別子"
  default     = ""
}

# アカウント設定
variable "db_instance_name" {
  description = "インスタンス作成時に作成するDB名"
  default     = ""
}

variable "db_character_set_name" {
  description = "Oracle系のDBエンコーディング用文字セット"
  default     = "" # 参照: https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Appendix.OracleCharacterSets.html
}

variable "db_master_username" {
  description = "マスターDBインスタンス マスターユーザー名"
  default     = "" # 1字目は英字
}

variable "db_master_password" {
  description = "マスターDBインスタンス マスターパスワード"
  default     = "" # ASCII文字 8文字以上
}

# インスタンス設定
variable "db_engine" {
  description = "DBエンジンのタイプ(Aurora以外)"
  default     = "" # e.g.) mysql
}

variable "db_engine_version" {
  description = "DBエンジンのバージョン"
  default     = "" # e.g.) 5.7
}

variable "db_instance_class" {
  description = "マスターDBインスタンス インスタンスサイズ"
  default     = "db.t2.micro"
}

# ストレージ設定
variable "db_storage_type" {
  description = "マスターDBインスタンス ストレージタイプ"
  default     = "gp2" # gp2 or io1
}

variable "db_storage_iops" {
  description = "マスターDBインスタンス ストレージIOPS(io1のみ)"
  type        = number
  default     = 1000
}

variable "db_storage_encrypted" {
  description = "マスターDBインスタンス ストレージの暗号化"
  type        = bool
  default     = false
}

variable "db_allocated_storage" {
  description = "マスターDBインスタンス ストレージ割り当て(GiB)"
  type        = number
  default     = 20 # min: 100GiB, MAX: 65536GiB
}

variable "db_max_allocated_storage" {
  description = "マスターDBインスタンス 最大ストレージしきい値(GiB)"
  type        = number
  default     = 1000 # min: 101GiB, MAX: 65536GiB
}

# ネットワーク設定
variable "db_port" {
  description = "DBへの接続を受け付けるポート番号"
  type        = number
  default     = 3306
}

variable "enable_db_multi_az" {
  description = "マルチAZ配置 (DBインスタンスのスタンバイを別のAZに配置するかを指定)"
  type        = bool
  default     = false
}

variable "db_availability_zone" {
  description = "DBインスタンを作成するAZ (※マルチAZ配置がtrueのとき、選択できない)"
  default     = ""
}

# バックアップ設定
variable "db_backup_retention_period" {
  description = "バックアップ保存期間(日) (※0にした場合、自動バックアップが無効になる)"
  type        = number
  default     = 0 # 0 - 35
}

variable "db_backup_window" {
  description = "自動バックアップの作成時間(UTC)"
  default     = "" # Format: 01:23-23:01
}

variable "db_copy_tags_to_snapshot" {
  description = "タグをスナップショットへコピー"
  type        = bool
  default     = false
}

# パフォーマンスインサイト設定
variable "db_performance_insights_enabled" {
  description = "Performance Insightsの有効化"
  type        = bool
  default     = false
}

variable "db_performance_insights_retention_period" {
  description = "Performance Insightsの保持期間"
  type        = number
  default     = 7
}

# モニタリング設定
variable "db_monitoring_interval" {
  description = "モニタリングのインターバル"
  type        = number
  default     = 0 # 0, 1, 5, 10, 15, 30 or 60
}

variable "db_monitoring_role_arn" {
  description = "モニタリングのロールARN" # <- TODO: Data経由にする
  default     = ""
}

variable "db_enabled_cloudwatch_logs_exports" {
  description = "cloudwatchlogに出力するログ"
  type        = list(string)
  default     = []
}

# メンテナンス設定
variable "apply_immediately" {
  description = "変更をすぐに反映"
  type        = bool
  default     = false
}

variable "db_maintenance_window" {
  description = "DBに適用される保留中、メンテナンスの期間 (※db_backup_windowがtrueの際設定不要)"
  default     = "Mon:00:00-Mon:01:00" # Format: Mon:01:23-Mon:23:01
}

variable "db_auto_minor_version_upgrade" {
  description = "マイナーバージョン自動アップグレードの有効化"
  type        = bool
  default     = true
}

# 削除保護設定
variable "db_deletion_protection" {
  description = "DB削除保護の有効化"
  type        = bool
  default     = false
}

variable "db_skip_final_snapshot" {
  description = "DBインスタンスの削除時、スナップショットを作成しない"
  type        = bool
  default     = false
}

##################################################
# RDS Instance (Replica)
##################################################
variable "db_replica_count" {
  description = "レプリカDBの作成数"
  type        = number
  default     = 0
}

variable "db_replica_identifier" {
  description = "レプリカDBインスタンス識別子 (複数作成する場合、末尾に自動で数字割り振る)"
  default     = ""
}

variable "db_replica_instance_class" {
  description = "レプリカDBインスタンス インスタンスサイズ"
  default     = "db.t2.small"
}

variable "db_replica_storage_type" {
  description = "レプリカDBインスタンス ストレージタイプ"
  default     = "gp2" # gp2 or io1
}

variable "db_replica_storage_iops" {
  description = "レプリカDBインスタンス ストレージIOPS(io1のみ)"
  type        = number
  default     = 1000
}

variable "enable_db_replica_multi_az" {
  description = "マルチAZ配置 (スタンバイインスタンを別のアベイラビリティーゾーンに配置するかを指定)"
  type        = bool
  default     = true
}

variable "db_replica_port" {
  description = "DBへの接続を受け付けるポート番号"
  type        = number
  default     = 3306
}

variable "db_replica_availability_zones" {
  description = "DBインスタンスを作成するAZ(db_replica_countで指定した個数指定) (※マルチAZ配置がtrueのとき、選択できない)"
  default     = ""
}

variable "db_replica_copy_tags_to_snapshot" {
  description = "タグをスナップショットへコピー"
  type        = bool
  default     = false
}

variable "db_replica_monitoring_interval" {
  description = "モニタリングのインターバル"
  type        = number
  default     = 0
}

variable "db_replica_monitoring_role_arn" {
  description = "モニタリングのロールARN" # <- TODO: Data経由にする
  default     = ""
}

variable "db_replica_auto_minor_version_upgrade" {
  description = "マイナーバージョン自動アップグレードの有効化"
  type        = bool
  default     = true
}


variable "db_replica_deletion_protection" {
  description = "DB削除保護の有効化"
  type        = bool
  default     = false
}
