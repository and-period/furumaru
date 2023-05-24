module "dns_olibio_jp" {
  source = "./../../../modules/route53"

  #####################################################################
  # Common
  #####################################################################
  tags     = var.tags

  #####################################################################
  # Route53 (zone)
  #####################################################################
  domain  = "olibio.jp"
  comment = "ECプラットフォーム用"

  subdomains = []
}
