##################################################
# EC2 Instance
# - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/instance
##################################################
resource "aws_instance" "this" {
  # --- 基本設定 -----------------
  ami                  = var.ami_name != "" ? data.aws_ami.this[0].id : null
  instance_type        = var.instance_type
  key_name             = var.key_name
  iam_instance_profile = var.iam_role != "" ? var.iam_role : null
  user_data            = var.user_data != "" ? var.user_data : null

  # --- ストレージ設定 --------------
  ebs_optimized = var.ebs_optimized

  # ブートディスク設定
  root_block_device {
    volume_type           = var.root_volume_type
    volume_size           = var.root_volume_size
    encrypted             = var.root_volume_encrypted
    kms_key_id            = var.kms_alias_name != "" ? data.aws_kms_alias.this[0].target_key_arn : null
    delete_on_termination = var.root_volume_delete_on_termination
  }

  # --- ネットワーク設定 -------------
  availability_zone = data.aws_subnet.this.availability_zone

  network_interface {
    network_interface_id = aws_network_interface.this.id
    device_index         = var.eni_device_index # 固定
  }

  # --- モニタリング設定 -------------
  monitoring = var.monitoring

  # --- 削除保護 -----------------
  disable_api_termination              = var.disable_api_termination
  disable_api_stop                     = var.disable_api_stop
  instance_initiated_shutdown_behavior = var.instance_initiated_shutdown_behavior

  tags = merge(
    var.tags,
    { Name = var.name },
  )
  volume_tags = merge(
    var.tags,
    { Name = var.name },
  )

  depends_on = [
    aws_network_interface.this,
  ]
  lifecycle {
    ignore_changes = [ami]
  }
}
