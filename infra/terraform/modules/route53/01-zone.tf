##################################################
# Route 53
# - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/route53_zone
##################################################
resource "aws_route53_zone" "domain" {
  name    = var.domain
  comment = var.comment

  force_destroy = true

  tags = merge(
    var.tags,
    { Name = var.domain },
  )
}

resource "aws_route53_zone" "subdomain" {
  for_each = local.subdomains

  name    = each.value.domain
  comment = each.value.comment

  force_destroy = true

  tags = merge(
    var.tags,
    { Name = var.domain },
  )
}

resource "aws_route53_record" "subdomain" {
  for_each = local.subdomains

  zone_id = aws_route53_zone.domain.zone_id
  records = aws_route53_zone.subdomain[each.key].name_servers

  name = each.value.domain
  type = "NS"
  ttl  = each.value.ttl
}
