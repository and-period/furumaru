##################################################
# Elastic IP
# - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/eip
##################################################
resource "aws_eip" "this" {
  count = var.enable_eip ? 1 : 0

  vpc = true

  tags = merge(
    var.tags,
    { Name = var.eip_name != "" ? var.eip_name : var.name },
  )
}

resource "aws_eip_association" "this" {
  count = var.enable_eip ? 1 : 0

  network_interface_id = aws_network_interface.this.id
  allocation_id        = aws_eip.this[0].id
}
