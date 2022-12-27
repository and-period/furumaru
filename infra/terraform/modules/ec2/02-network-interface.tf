##################################################
# Elastic Network Interface
# - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/network_interface
##################################################
resource "aws_network_interface" "this" {
  subnet_id       = data.aws_subnet.this.id
  security_groups = concat(data.aws_security_group.this[*].id, [])

  description       = var.eni_description
  private_ips       = length(var.eni_private_ips) > 0 ? concat(var.eni_private_ips, []) : null
  source_dest_check = var.source_dest_check

  tags = merge(
    var.tags,
    { Name = var.name },
  )
}
