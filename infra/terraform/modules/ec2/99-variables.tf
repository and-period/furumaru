##################################################
# Common
##################################################
variable "tags" {
  default = {}
}

variable "vpc_name" {
  default = ""
}

variable "subnet_name" {
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
# EC2 Instance - 基本設定
##################################################
variable "name" {
  description = "EC2インスタンス名"
  default     = ""
}

variable "ami_name" {
  description = "使用するAMI"
  default     = ""
}

variable "ami_owners" {
  description = "AMI オーナー名"
  default     = []
}

variable "ami_most_recent" {
  description = "AMI が複数取得できた場合最新のを使用"
  type        = bool
  default     = true
}

variable "instance_type" {
  description = "インスタンスタイプ"
  default     = "t3.micro"
}

variable "key_name" {
  description = "EC2へログインするのに使用するキーペア名"
  default     = ""
}

variable "iam_role" {
  description = "インスタンスに割り当てるIAMロール"
  default     = ""
}

variable "user_data" {
  description = "ユーザーデータ（インスタンスの起動時に実行するコマンドスクリプト）"
  default     = ""
}

##################################################
# EC2 Instance - ストレージ設定
##################################################
variable "ebs_optimized" {
  description = "EBS 最適化インスタンス設定"
  type        = bool
  default     = false
}

variable "root_volume_type" {
  description = "ルートディスクのストレージタイプ"
  default     = "gp3" # standard, gp2, gp3, io1, sc1, or st1
}

variable "root_volume_size" {
  description = "ルートディスクの容量 (GiB)"
  type        = number
  default     = 8
}

variable "root_volume_encrypted" {
  description = "ルートディスクの暗号化"
  type        = bool
  default     = false
}

variable "root_volume_delete_on_termination" {
  description = "EC2削除時、ルートディスクも削除"
  type        = bool
  default     = false
}

variable "ebs_volumes" {
  description = "追加ディスク設定"
  type = list(object({
    name        = string # 追加ディスク名
    device_name = string # 追加ディスクの割り当て先デバイス名
    type        = string # 追加ディスクタイプ # standard, gp2, gp3, io1, sc1, or st1
    size        = number # 追加ディスクサイズ
    encrypted   = bool   # EBSの暗号化
  }))
  default = []
}

variable "ebs_volume_skip_destroy" {
  description = "EC2インスタンス削除時、追加ディスクも削除"
  type        = bool
  default     = false
}

##################################################
# EC2 Instance - ネットワーク設定
##################################################
variable "associate_public_ip_address" {
  description = "パブリックIPの自動割り当て"
  type        = bool
  default     = true
}

variable "source_dest_check" {
  description = "送信元/送信先チェックの有効化"
  type        = bool
  default     = true
}

variable "enable_eip" {
  description = "Elastic IPの関連付け"
  type        = bool
  default     = false
}

variable "eip_name" {
  description = "Elastic IP名"
  default     = ""
}

variable "eni_description" {
  description = "ネットワークインターフェースの説明"
  default     = ""
}

variable "eni_private_ips" {
  description = "ネットワークインターフェースに割り当てるプライベートアドレス一覧"
  type        = list(string)
  default     = []
}

variable "eni_device_index" {
  description = "ネットワークインターフェースを割り当てるデバイス名(数字を指定)"
  type        = number
  default     = 0 # e.g.) 0 -> eth0, 2 -> eth2
}

##################################################
# EC2 Instance - モニタリング設定
##################################################
variable "monitoring" {
  description = "EC2のモニタリング連携"
  type        = bool
  default     = false
}

##################################################
# EC2 Instance - 削除保護
##################################################
variable "disable_api_termination" {
  description = "インスタンスの終了保護"
  type        = bool
  default     = false
}

variable "disable_api_stop" {
  description = "インスタンスの停止保護"
  type        = bool
  default     = false
}

variable "instance_initiated_shutdown_behavior" {
  description = "シャットダウン時の動作" # https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/terminating-instances.html#Using_ChangingInstanceInitiatedShutdownBehavior
  default     = "stop"        # stop or terminate
}
