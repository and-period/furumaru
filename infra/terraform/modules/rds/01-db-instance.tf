##################################################
# RDS Instance (Master)
# - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/db_instance
##################################################
resource "aws_db_instance" "master" {
  count = var.db_instance_count

  # 基本情報
  identifier = var.db_identifier

  # インスタンス設定
  engine         = var.db_engine
  engine_version = var.db_engine_version
  instance_class = var.db_instance_class

  # アカウント設定
  name               = var.db_instance_name
  character_set_name = var.db_character_set_name
  username           = var.db_master_username
  password           = var.db_master_password

  # ストレージ設定
  storage_type          = var.db_storage_type
  iops                  = var.db_storage_type == "io1" ? var.db_storage_iops : null
  storage_encrypted     = var.db_storage_encrypted
  allocated_storage     = var.db_allocated_storage
  max_allocated_storage = var.db_max_allocated_storage
  kms_key_id            = var.kms_alias_name != "" ? data.aws_kms_alias.this[0].target_key_arn : null

  # ネットワーク設定
  port                   = var.db_port
  multi_az               = var.enable_db_multi_az
  availability_zone      = var.enable_db_multi_az ? null : var.db_availability_zone
  db_subnet_group_name   = aws_db_subnet_group.this[0].id
  vpc_security_group_ids = length(data.aws_security_group.this) > 0 ? concat(data.aws_security_group.this[*].id, []) : null

  parameter_group_name = length(var.db_parameters) > 0 ? aws_db_parameter_group.this[0].id : null
  option_group_name    = length(var.db_options) > 0 ? aws_db_option_group.this[0].id : null

  # バックアップ設定
  backup_retention_period = var.db_backup_retention_period
  backup_window           = var.db_backup_window
  copy_tags_to_snapshot   = var.db_copy_tags_to_snapshot

  # Performance Insights設定
  performance_insights_enabled          = var.db_performance_insights_enabled
  performance_insights_retention_period = var.db_performance_insights_enabled ? var.db_performance_insights_retention_period : null
  performance_insights_kms_key_id       = var.db_performance_insights_enabled ? data.aws_kms_alias.this[0].target_key_arn : null

  # モニタリング設定
  monitoring_interval             = var.db_monitoring_interval
  monitoring_role_arn             = var.db_monitoring_role_arn
  enabled_cloudwatch_logs_exports = var.db_enabled_cloudwatch_logs_exports

  # メンテナンス設定
  apply_immediately          = var.apply_immediately
  auto_minor_version_upgrade = var.db_auto_minor_version_upgrade
  maintenance_window         = var.db_backup_window != "" ? null : var.db_maintenance_window

  # 削除保護
  deletion_protection = var.db_deletion_protection
  skip_final_snapshot = var.db_skip_final_snapshot

  lifecycle {
    ignore_changes = [
      engine,
      name,
      username,
      performance_insights_kms_key_id,
    ]
  }

  tags = merge(
    var.tags,
    { Name = var.db_identifier },
  )
}

##################################################
# RDS Instance (Replica)
##################################################
resource "aws_db_instance" "replica" {
  count = var.db_replica_count

  replicate_source_db = aws_db_instance.master[0].identifier
  identifier = var.db_replica_count == 1 ? var.db_replica_identifier : format(
    "%s-%s",
    var.db_replica_identifier,
    count.index
  )

  instance_class = var.db_replica_instance_class

  storage_type      = var.db_replica_storage_type
  iops              = var.db_replica_storage_type == "io1" ? var.db_replica_storage_iops : null
  storage_encrypted = var.db_storage_encrypted
  kms_key_id        = var.kms_alias_name != "" ? data.aws_kms_alias.this[0].target_key_arn : null

  multi_az = var.enable_db_replica_multi_az
  port     = var.db_replica_port
  availability_zone = var.enable_db_replica_multi_az ? null : element(
    var.db_replica_availability_zones,
    count.index
  )

  copy_tags_to_snapshot = var.db_replica_copy_tags_to_snapshot

  monitoring_interval = var.db_replica_monitoring_interval
  monitoring_role_arn = var.db_replica_monitoring_role_arn

  apply_immediately          = var.apply_immediately
  auto_minor_version_upgrade = var.db_replica_auto_minor_version_upgrade

  deletion_protection = var.db_replica_deletion_protection

  tags = merge(
    var.tags,
    {
      Name = var.db_replica_count == 1 ? var.db_replica_identifier : format(
        "%s-%s",
        var.db_replica_identifier,
        count.index,
      )
    }
  )
}
