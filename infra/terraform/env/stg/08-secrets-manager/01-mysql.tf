module "mysql" {
  source = "./../../../modules/secrets-manager"

  #####################################################################
  # Common
  #####################################################################
  tags           = var.tags
  kms_alias_name = ""

  #####################################################################
  # Secrets Manager
  #####################################################################
  name                    = format("%s/%s", local.prefix, "mysql")
  description             = "staging環境用のAmazon RDS for MySQL接続用の認証情報"
  recovery_window_in_days = 0
}
