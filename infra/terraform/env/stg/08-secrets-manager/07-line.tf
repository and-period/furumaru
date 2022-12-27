module "line" {
  source = "./../../../modules/secrets-manager"

  #####################################################################
  # Common
  #####################################################################
  tags           = var.tags
  kms_alias_name = ""

  #####################################################################
  # Secrets Manager
  #####################################################################
  name                    = format("%s/%s", local.prefix, "line")
  description             = "staging環境用のLINE API接続用のシークレット"
  recovery_window_in_days = 0
}
