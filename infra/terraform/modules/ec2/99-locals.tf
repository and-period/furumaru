locals {
  # インスタンスタイプがバーストタイプ対応かの検証
  is_t_instance_type = replace(var.instance_type, "/^t(2|3|3a){1}\\..*$/", "1") == "1" ? true : false

  # EBSボリューム
  ebs_volumes = {
    for e in var.ebs_volumes : e.name => e
  }
}
