##################################################
# Common
##################################################
variable "tags" {
  default = {}
}

variable "vpc_name" {
  default = ""
}

variable "subnet_names" {
  type    = list(string)
  default = []
}

##################################################
# Network ACL
##################################################
variable "name" {
  description = "Network ACL名"
  default     = ""
}

variable "ingress_rules" {
  description = "インバウンドルール一覧"
  type = list(object({
    rule_no    = number # ルール番号
    action     = string # 許可/拒否 (許可: allow, 拒否: deny)
    protocol   = string # プロトコル(icmp, tcp, udp or -1)
    from_port  = number # ポート範囲(全て開けるときは 0)
    to_port    = number # ポート範囲(全て開けるときは 0)
    cidr_block = string # ルールを適用するIPv4 CIDRブロック
  }))
  default = []
}

variable "egress_rules" {
  description = "アウトバウンドルール一覧"
  type = list(object({
    rule_no    = number # ルール番号
    action     = string # 許可/拒否 (許可: allow, 拒否: deny)
    protocol   = string # プロトコル(icmp, tcp, udp or -1)
    from_port  = number # ポート範囲(全て開けるときは 0)
    to_port    = number # ポート範囲(全て開けるときは 0)
    cidr_block = string # ルールを適用するIPv4 CIDRブロック
  }))
  default = []
}
