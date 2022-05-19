##############################
# Client VPN Endpoint
# - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/ec2_client_vpn_endpoint
# - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/ec2_client_vpn_network_association
##############################
resource "aws_ec2_client_vpn_endpoint" "this" {
  count = var.create_client_vpn ? 1 : 0

  description       = var.client_vpn_description
  client_cidr_block = var.client_vpn_cidr_block

  # とりあえず Active Directory 認証は考慮しない
  server_certificate_arn = data.aws_acm_certificate.server[0].arn
  authentication_options {
    type                       = "certificate-authentication"
    root_certificate_chain_arn = data.aws_acm_certificate.client[0].arn
  }

  connection_log_options {
    enabled               = var.enable_client_vpn_connection_log
    cloudwatch_log_group  = var.enable_client_vpn_connection_log ? var.client_vpn_cloudwatch_log_group_name : null
    cloudwatch_log_stream = var.enable_client_vpn_connection_log ? var.client_vpn_cloudwatch_log_stream_name : null
  }

  dns_servers        = concat(var.client_vpn_dns_servers, [])
  transport_protocol = var.client_vpn_transport_protocol
  split_tunnel       = var.client_vpn_split_tunnel

  tags = merge(
    local.tags,
    { Name = var.client_vpn_name },
  )

  lifecycle {
    ignore_changes = [
      client_cidr_block,
      authentication_options,
      transport_protocol,
    ]
  }
}

# Client VPNとサブネットを関連付け
resource "aws_ec2_client_vpn_network_association" "this" {
  count = var.create_client_vpn ? length(var.client_vpn_associate_subnet_names) : 0

  client_vpn_endpoint_id = aws_ec2_client_vpn_endpoint.this[0].id
  subnet_id = element(
    aws_subnet.this[*].id,
    index(
      aws_subnet.this[*].tags.Name,
      element(var.client_vpn_associate_subnet_names, count.index)
    ),
  )
}
