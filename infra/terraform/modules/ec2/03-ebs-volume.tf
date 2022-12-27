##################################################
# EBS Volume
# - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/ebs_volume
##################################################
resource "aws_ebs_volume" "this" {
  for_each = local.ebs_volumes

  availability_zone = data.aws_subnet.this.availability_zone

  type = each.value.type
  size = each.value.size

  encrypted  = each.value.encrypted
  kms_key_id = var.kms_alias_name != "" ? data.aws_kms_alias.this[0].target_key_arn : null

  tags = merge(
    var.tags,
    { Name = each.key },
  )
}

resource "aws_volume_attachment" "this" {
  for_each = local.ebs_volumes

  instance_id = aws_instance.this.id
  volume_id   = aws_ebs_volume.this[each.key].id

  device_name = each.value.device_name

  skip_destroy = var.ebs_volume_skip_destroy
}
