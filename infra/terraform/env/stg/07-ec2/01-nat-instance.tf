module "nat-instance" {
  source = "./../../../modules/ec2"

  #####################################################################
  # Common
  #####################################################################
  tags                 = var.tags
  vpc_name             = var.vpc_name
  subnet_name          = "furumaru-stg-sub-pub-1a"
  security_group_names = ["furumaru-stg-sg-nat"]
  kms_alias_name       = ""

  #####################################################################
  # EC2 Instance - 基本設定
  #####################################################################
  name = "furumaru-stg-nat"

  ami_name        = "amzn2-ami-hvm-*-x86_64-ebs"
  ami_owners      = ["amazon"]
  ami_most_recent = true

  instance_type = "t3.micro"
  key_name      = var.key_name
  user_data     = file("${path.module}/user_data/nat-instance.sh")

  ##################################################
  # EC2 Instance - ストレージ設定
  ##################################################
  ebs_optimized = false

  root_volume_type                  = "gp3"
  root_volume_size                  = 8 # 8GiB
  root_volume_encrypted             = false
  root_volume_delete_on_termination = true

  ebs_volumes             = []
  ebs_volume_skip_destroy = true

  ##################################################
  # EC2 Instance - ネットワーク設定
  ##################################################
  associate_public_ip_address = true
  source_dest_check           = false # NAT利用の場合、無効にする必要があり
  enable_eip                  = true
  eip_name                    = "furumaru-stg-eip-nat"

  eni_description  = "for nat instance"
  eni_private_ips  = []
  eni_device_index = 0

  ##################################################
  # EC2 Instance - モニタリング設定
  ##################################################
  monitoring = false

  ##################################################
  # EC2 Instance - 削除保護
  ##################################################
  disable_api_termination              = false
  disable_api_stop                     = false
  instance_initiated_shutdown_behavior = "stop"
}
