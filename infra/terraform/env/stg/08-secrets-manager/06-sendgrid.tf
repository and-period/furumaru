module "sendgrid" {
  source = "./../../../modules/secrets-manager"

  #####################################################################
  # Common
  #####################################################################
  tags           = var.tags
  kms_alias_name = ""

  #####################################################################
  # Secrets Manager
  #####################################################################
  name                    = format("%s/%s", local.prefix, "sendgrid")
  description             = "staging環境用のSendGridアクセス用API Key"
  recovery_window_in_days = 0
}
