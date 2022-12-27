module "route_pri_app" {
  source = "./../../../modules/route"

  #####################################################################
  # Common
  #####################################################################
  vpc_name         = var.vpc_name
  route_table_name = "furumaru-stg-rt-pri-app"

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
  routes = [
    {
      destination_cidr_block = "0.0.0.0/0"
      target_resource        = "instance"
      target_resource_name   = "furumaru-stg-nat"
    },
  ]
}
