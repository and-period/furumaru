module "furumaru_db" {
  source = "./../../../modules/rds"

  #####################################################################
  # Common
  #####################################################################
  vpc_name = var.vpc_name
  tags = var.tags

  #####################################################################
  # サブネットグループ
  #####################################################################
  subnet_group_name = "furumaru-stg-sub-db"
  subnet_group_description = "For database"

  subnet_names = [
    "furumaru-stg-sub-pri-db-1a",
    "furumaru-stg-sub-pri-db-1c",
  ]

  #####################################################################
  # パラメータグループ
  #####################################################################
  create_db_parameter_group = true

  db_parameter_group_name = "furumaru-stg-db-params"
  db_parameter_group_description = "For staging database"
  db_parameter_group_family = "mysql8.0"

  db_parameters = [
    {name="character_set_client", value="utf8mb4"},
    {name="character_set_server", value="utf8mb4"},
    {name="character_set_database", value="utf8mb4"},
    {name="lc_time_names", value="ja_JP"},
  ]

  #####################################################################
  # オプショングループ
  #####################################################################
  create_db_option_group = false

  #####################################################################
  # DBインスタンス (マスター)
  #####################################################################
  ### 基本情報 ############################################
  db_identifier = "furumaru-stg-db"

  ### 認証情報の設定 ######################################
  db_master_username = var.db_master_username
  db_master_password = var.db_master_password

  ### DBエンジンの設定 ####################################
  db_engine = "mysql"
  db_engine_version = "8.0.31"
  db_instance_class = "db.t4g.micro"

  ### ストレージ設定 ######################################
  db_storage_type = "gp2"
  db_storage_encrypted = false
  db_allocated_storage = 20
  db_max_allocated_storage = 0

  ### 可用性と耐久性設定 ##################################
  enable_db_multi_az = false
  db_availability_zone = "ap-northeast-1a"

  ### ネットワーク設定 ####################################
  security_group_names = ["furumaru-stg-sg-db"]
  db_port = 3306

  ### データベース設定 ####################################
  db_instance_name = ""
  db_character_set_name = ""

  ### バックアップ設定 ####################################
  db_backup_retention_period = 1
  db_backup_window = "18:00-21:00"
  db_copy_tags_to_snapshot = true

  ### 暗号化設定 ##########################################
  kms_alias_name = ""

  ### パフォーマンスインサイト設定 ########################
  db_performance_insights_enabled = false
  db_performance_insights_retention_period = 7

  ### モニタリング設定 ####################################
  db_monitoring_interval = 0
  db_monitoring_role_arn = ""

  ### (WIP) ログのエクスポート設定 ########################

  ### メンテナンス設定 ####################################
  db_maintenance_window = ""
  db_auto_minor_version_upgrade = true

  ### 削除保護設定 ########################################
  db_deletion_protection = false
  db_skip_final_snapshot = false

  #####################################################################
  # DBインスタンス (レプリカ)
  #####################################################################
  db_replica_count = 0
}
