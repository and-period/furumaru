##################################################
# Common
##################################################
variable "tags" {
  default = {}
}

variable "vpc_name" {
  default = ""
}

##################################################
# Security Group
##################################################
variable "name" {
  description = "Security Group名"
  default     = ""
}

variable "description" {
  description = "Security Groupの説明"
  default     = ""
}

variable "ingress_rules" {
  description = "インバウンドルール リスト"
  type = list(object({
    description          = string       # 説明
    protocol             = string       # プロトコル(icmp, tcp, udp or -1)
    from_port            = number       # ポート範囲(全て開けるときは 0)
    to_port              = number       # ポート範囲(全て開けるときは 0)
    cidr_blocks          = list(string) # ルールを適用するIPv4 CIDRブロックリスト
    security_group_names = list(string) # ルールを適用するセキュリティグループIDリスト
    # ipv6_cidr_blocks   = list(string) # ルールを適用するIPv6 CIDRブロックリスト
    prefix_list_ids = list(string) # ルールを適用するエンドポイントリスト
  }))
  default = []
}

variable "egress_rules" {
  description = "アウトバウンドルール リスト"
  type = list(object({
    description          = string       # 説明
    protocol             = string       # プロトコル(icmp, tcp, udp or -1)
    from_port            = number       # ポート範囲(全て開けるときは 0)
    to_port              = number       # ポート範囲(全て開けるときは 0)
    cidr_blocks          = list(string) # ルールを適用するIPv4 CIDRブロックリスト
    security_group_names = list(string) # ルールを適用するセキュリティグループIDリスト
    # ipv6_cidr_blocks   = list(string) # ルールを適用するIPv6 CIDRブロックリスト
    prefix_list_ids = list(string) # ルールを適用するエンドポイントリスト
  }))
  default = []
}
