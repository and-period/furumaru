##################################################
# Common
##################################################
variable "name" {
  default = ""
}

variable "tags" {
  default = {}
}

##################################################
# VPC
##################################################
variable "cidr_block" {
  description = "VPCのIPv4 CIDRブロック"
  default     = ""
}

variable "assign_generated_ipv6_cidr_block" {
  description = "VPCのIPv6 CIDRブロックの有効化(true: Amazonが提供したIPv6 CIDRブロックが適用)"
  type        = bool
  default     = false # true: Amazonが提供したIPv6 CIDRブロックが適用
}

variable "instance_tenancy" {
  description = "VPCのテナンシー"
  default     = "default" # default or dedicated
}

variable "enable_dns_support" {
  description = "VPCのEC2 DNSホスト名の有効化"
  type        = bool
  default     = true
}

variable "enable_dns_hostnames" {
  description = "VPCのEC2 DNS名前解決の有効化"
  type        = bool
  default     = false
}

##################################################
# DHCP Option Set
##################################################
variable "create_dhcp_options" {
  description = "DHCPオプションセットの作成"
  type        = bool
  default     = false
}

variable "dhcp_options_name" {
  description = "DHCPオプションセット名"
  default     = ""
}

variable "domain_name" {
  description = "/etc/resolv.confに設定する名前解決をするときのデフォルトドメイン名"
  default     = ""
}

variable "domain_name_servers" {
  description = "/etc/resolv.confに設定するDNSサーバのリスト"
  type        = list(string)
  default     = [] # AWSデフォルトのDNSサーバを利用する場合、AmazonProvidedDNSを指定
}

variable "ntp_servers" {
  description = "NTPサーバのリスト"
  type        = list(string)
  default     = []
}

variable "netbios_name_servers" {
  description = "NETBIOSネームサーバのリスト"
  type        = list(string)
  default     = []
}

variable "netbios_node_type" {
  description = "NETBIOSのノードタイプ(1,2,4,8から選択)"
  type        = number
  default     = 2 # AWSではブロードキャストとマルチキャストがサポートされていないため2を推奨とのこと
}

##################################################
# Route Table
##################################################
variable "route_tables" {
  description = "ルートテーブル名一覧"
  type        = list(string)
  default     = []
}

##################################################
# Subnet
##################################################
variable "subnets" {
  description = "サブネット一覧"
  type = list(object({
    name                    = string # サブネット名
    availability_zone       = string # アベイラビリティーゾーン e.g.) us-east-1a
    cidr_block              = string # サブネットのIPv4 CIDRブロック e.g.) 10.0.0.0/24
    map_public_ip_on_launch = bool   # パブリックIPv4アドレス自動割り当ての有効化
    route_table_name        = string # 関連付けるルートテーブル名
  }))
  default = []
}

variable "enable_dns64" {
  description = "DNS64の有効化"
  type        = bool
  default     = false
}

variable "enable_resource_name_dns_a_record_on_launch" {
  description = "リソース名のDNS Aレコードに対するクエリ応答の有効化"
  type        = bool
  default     = false
}

variable "enable_resource_name_dns_aaaa_record_on_launch" {
  description = "リソース名のDNS AAAAレコードに対するクエリ応答の有効化"
  type        = bool
  default     = false
}

##################################################
# Internet Gateway
##################################################
variable "create_internet_gateway" {
  description = "インターネットゲートウェイの作成"
  type        = bool
  default     = false
}

variable "internet_gateway_name" {
  description = "インターネットゲートウェイ名"
  default     = ""
}

##################################################
# NAT Gateway
##################################################
variable "nat_gateways" {
  description = "NATゲートウェイ一覧"
  type = list(object({
    name        = string # NATゲートウェイ名
    subnet_name = string # NATゲートウェイを配置するサブネット名
    eip_name    = string # NATゲートウェイに関連付けるElastic IP名
  }))
  default = []
}

##############################
# Virtual Private Gateway
##############################
variable "create_vpn_gateway" {
  description = "VPNゲートウェイの作成"
  type        = bool
  default     = false
}

variable "vpn_gateway_name" {
  description = "VPNゲートウェイ名"
  default     = ""
}

variable "vpn_gateway_asn" {
  description = "BGPに使用するVPNゲートウェイのASN"
  type        = number
  default     = 0
}

##############################
# Client VPN
##############################
variable "create_client_vpn" {
  description = "Client VPNの作成"
  type        = bool
  default     = false
}

variable "client_vpn_name" {
  description = "Client VPN名"
  default     = ""
}

variable "client_vpn_description" {
  description = "Client VPNの説明"
  default     = ""
}

variable "client_vpn_cidr_block" {
  description = "クライアント IPv4 CIDR"
  default     = ""
}

variable "acm_certificate_server_domain" {
  description = "ACMにあるサーバー証明書のドメイン名"
  default     = ""
}

variable "acm_certificate_client_domain" {
  description = "ACMにあるクライアント証明書のドメイン名"
  default     = ""
}

variable "enable_client_vpn_connection_log" {
  description = "クライアント接続の詳細の記録有効化"
  type        = bool
  default     = false
}

variable "client_vpn_cloudwatch_log_group_name" {
  description = "クライアント接続のログ保存先CloudWatch Logsロググループ名"
  default     = ""
}

variable "client_vpn_cloudwatch_log_stream_name" {
  description = "クライアント接続のログ保存先CloudWatch Logsログストリーム名"
  default     = ""
}

variable "client_vpn_dns_servers" {
  description = "DNSサーバー IPアドレス リスト"
  type        = list(string)
  default     = []
}

variable "client_vpn_transport_protocol" {
  description = "TLSセッションで使用されるトランスポートプロトコル"
  default     = "" # udp or tcp
}

variable "client_vpn_split_tunnel" {
  description = "VPN経由で送信されるトラフィックの制限を有効化"
  type        = bool
  default     = false
}

variable "client_vpn_associate_subnet_names" {
  description = "Client VPNに関連付けるサブネットリスト"
  type        = list(string)
  default     = []
}
