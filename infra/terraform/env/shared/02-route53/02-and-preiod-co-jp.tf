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

  subdomains = []
}
