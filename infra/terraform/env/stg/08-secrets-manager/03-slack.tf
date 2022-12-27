module "slack" {
  source = "./../../../modules/secrets-manager"

  #####################################################################
  # Common
  #####################################################################
  tags           = var.tags
  kms_alias_name = ""

  #####################################################################
  # Secrets Manager
  #####################################################################
  name                    = format("%s/%s", local.prefix, "slack")
  description             = "staging環境用のSlackアクセス用API Key"
  recovery_window_in_days = 0
}
