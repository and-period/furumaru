##################################################
# Common
##################################################
variable "tags" {
  default = {}
}

##################################################
# DynamoDB Table
##################################################
variable "name" {
  description = "テーブル名"
  default     = ""
}

variable "billing_mode" {
  description = "キャパシティーモード"
  default     = "PROVISIONED" // PROVISIONED or PAY_PER_REQUEST
}

variable "table_class" {
  description = "テーブルクラス"
  default     = "STANDARD" // STANDARD or STANDARD_INFREQUENT_ACCESS
}

variable "read_capacity" {
  description = "読み込みキャパシティー"
  type        = number
  default     = 0
}

variable "write_capacity" {
  description = "書き込みキャパシティー"
  type        = number
  default     = 0
}

variable "ttl_enabled" {
  description = "TTLの有効化"
  type        = bool
  default     = false
}

variable "ttl_attribute_name" {
  description = "TTL属性"
  default     = ""
}

variable "stream_enabled" {
  description = "ストリームの有効化"
  type        = bool
  default     = false
}

variable "stream_view_type" {
  description = "ストリームの表示タイプ"
  default     = "" # KEYS_ONLY, NEW_IMAGE, OLD_IMAGE, NEW_AND_OLD_IMAGES
}

variable "hash_key" {
  description = "パーティションキー"
  default     = ""
}

variable "range_key" {
  description = "ソートキー"
  default     = ""
}

variable "attributes" {
  description = "属性"
  type = list(object({
    name = string # 属性名
    type = string # 型 (S:string,N:number,B:binary)
  }))
  default = []
}

variable "local_secondary_indexes" {
  description = "ローカルセカンダリインデックス"
  type = list(object({
    name               = string # インデックス名
    range_key          = string # ソートキー
    projection_type    = string # 属性の射影 (ALL, INCLUDE or KEYS_ONLY)
    non_key_attributes = string # 属性名 (INCLUDEの場合のみ指定)
  }))
  default = []
}

variable "global_secondary_indexes" {
  description = "グローバルセカンダリインデックス"
  type = list(object({
    name               = string # インデックス名
    hash_key           = string # パーティションキー
    range_key          = string # ソートキー
    projection_type    = string # 属性の射影 (ALL, INCLUDE or KEYS_ONLY)
    non_key_attributes = string # 属性名 (INCLUDEの場合のみ指定)
    read_capacity      = number # 読み取りキャパシティー
    write_capacity     = number # 書き込みキャパシティー
  }))
  default = []
}
