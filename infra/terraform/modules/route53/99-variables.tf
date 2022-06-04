##################################################
# Common
##################################################
variable "tags" {
  default = {}
}

##################################################
# Route 53
##################################################
variable "domain" {
  description = "Route53 ドメイン名"
  default = ""
}

variable "comment" {
  description = "Route53 説明"
  default = ""
}

variable "subdomains" {
  description = "Route53 サブドメイン一覧"
  type = list(object({
    domain  = string # サブドメイン名
    comment = string # 説明
    ttl     = number # TTL (sec)
  }))
  default = []
}
