##################################################
# Common
##################################################
variable "vpc_name" {
  default = ""
}

variable "route_table_name" {
  default = ""
}

##################################################
# Route
##################################################
# ターゲットリソース名
# | ターゲットリソース           | target_resourceに記述する値 |
# | ---------------------------- | --------------------------- |
# | EC2インスタンス              | instance                    |
# | インターネットゲートウェイ   | internet_gateway            |
# | NATゲートウェイ              | nat_gateway                 |
# | ネットワークインターフェース | network_interface           |
# | ピアリング接続               | vpc_peering                 |
# | Transit Gateway              | transit_gateway             |
# | 仮想プライベートゲートウェイ | vpn_gateway                 |
variable "routes" {
  type = list(object({
    destination_cidr_block = string # 送信先IPv4 CIDRブロック
    target_resource        = string # ターゲットリソース名
    target_resource_name   = string # ターゲット名
  }))
  default = []
}
