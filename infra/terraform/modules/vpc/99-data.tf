##################################################
# Certificate Manager
# - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/acm_certificate
##################################################

# Client VPNで使用するACMの証明書を取得
# - 同じドメイン名が複数あった場合、最新のものを取得
data "aws_acm_certificate" "server" {
  count = var.create_client_vpn ? 1 : 0

  domain = var.acm_certificate_server_domain

  most_recent = true
}

data "aws_acm_certificate" "client" {
  count = var.create_client_vpn ? 1 : 0

  domain = var.acm_certificate_client_domain

  most_recent = true
}
