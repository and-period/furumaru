module "stripe" {
  source = "./../../../modules/secrets-manager"

  #####################################################################
  # Common
  #####################################################################
  tags           = var.tags
  kms_alias_name = ""

  #####################################################################
  # Secrets Manager
  #####################################################################
  name                    = format("%s/%s", local.prefix, "stripe")
  description             = "staging環境用のStripe認証用キー"
  recovery_window_in_days = 0
}
