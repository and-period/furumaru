module "route_pri_db" {
  source = "./../../../modules/route"

  #####################################################################
  # Common
  #####################################################################
  vpc_name         = var.vpc_name
  route_table_name = "furumaru-stg-rt-pri-db"

  #####################################################################
  # ルーティング
  #####################################################################
  # | ターゲットリソース           | target_resourceの値 | target_resource_nameの値               |
  # | ---------------------------- | ------------------- | -------------------------------------- |
  # | EC2インスタンス              | instance            | EC2インスタンスのNameタグ              |
  # | インターネットゲートウェイ   | internet_gateway    | インターネットゲートウェイのNameタグ   |
  # | NATゲートウェイ              | nat_gateway         | NATゲートウェイのNameタグ              |
  # | ネットワークインターフェース | network_interface   | ネットワークインターフェースのNameタグ |
  # | ピアリング接続               | vpc_peering         | ピアリング接続リソースのID               |
  # | Transit Gateway              | transit_gateway     | Transit GatewayのリソースID            |
  # | 仮想プライベートゲートウェイ | vpn_gateway         | 仮想プライベートゲートウェイのNameタグ |
  routes = []
}
