module "dns_and_period_co_jp" {
  source = "./../../../modules/route53"

  #####################################################################
  # Common
  #####################################################################
  tags     = var.tags

  #####################################################################
  # Route53 (zone)
  #####################################################################
  domain  = "and-period.co.jp"
  comment = "&. コーポレートサイト用"

  subdomains = [
    {
      domain  = "hitofusa.and-period.co.jp"
      comment = "ECプラットフォーム Shopify用"
      ttl     = 60
    },
    {
      domain  = "olibio.and-period.co.jp"
      comment = "ECプラットフォーム Olibio用"
      ttl     = 60
    },
  ]
}
