##################################################
# Common
##################################################
variable "tags" {
  default = {}
}

##################################################
# SQS Queue
##################################################
variable "queue_name" {
  description = "AWS SQS名"
  default     = ""
}

variable "fifo_queue" {
  description = "キュータイプ(FIFOの有効化)"
  type        = bool
  default     = false # true: FIFO, false: 標準
}

variable "content_based_deduplication" {
  description = "コンテンツに基づく重複排除 (キュータイプがFIFOのとき)"
  type        = bool
  default     = false
}

variable "deduplication_scope" {
  description = "重複排除スコープ (キュータイプがFIFOのとき)"
  default     = "queue" # messageGroup or queue
}

variable "fifo_throughput_limit" {
  description = "FIFOスループット制限 (キュータイプがFIFOのとき)"
  default     = "perQueue" # perQueue or perMessageGroupId
}

variable "visibility_timeout_seconds" {
  description = "可視性タイムアウト (sec)"
  type        = number
  default     = 30 # default: 30sec. from 0 to 43200 (12 hours)
}

variable "message_retention_seconds" {
  description = "メッセージ保持期間 (sec)"
  type        = number
  default     = 345600 # default: 4days. from 60 (1 minute) to 1209600 (14 days)
}

variable "delay_seconds" {
  description = "配信遅延 (sec)"
  type        = number
  default     = 0 # default: 0sec. from 0 to 900 (15 minutes)
}

variable "max_message_size" {
  description = "最大メッセージサイズ (KiB)"
  type        = number
  default     = 262144 # default: 256KiB. from 1024 bytes (1 KiB) up to 262144 bytes (256 KiB)
}

variable "receive_wait_time_seconds" {
  description = "メッセージ受信待機時間 (sec)"
  type        = number
  default     = 0 # default: 0sec. from 0 to 20 (seconds)
}
