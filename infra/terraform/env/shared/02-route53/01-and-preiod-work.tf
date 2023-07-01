module "dns_and_period_work" {
  source = "./../../../modules/route53"

  #####################################################################
  # Common
  #####################################################################
  tags     = var.tags

  #####################################################################
  # Route53 (zone)
  #####################################################################
  domain  = "and-period.work"
  comment = "&. 検証用"

  subdomains = [
    {
      domain  = "furumaru.and-period.work"
      comment = "ふるさとマルシェ 本番用"
      ttl     = 60
    },
    {
      domain  = "furumaru-stg.and-period.work"
      comment = "ふるさとマルシェ 検証用"
      ttl     = 60
    },
  ]
}
