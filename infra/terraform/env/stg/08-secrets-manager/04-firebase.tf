module "firebase" {
  source = "./../../../modules/secrets-manager"

  #####################################################################
  # Common
  #####################################################################
  tags           = var.tags
  kms_alias_name = ""

  #####################################################################
  # Secrets Manager
  #####################################################################
  name                    = format("%s/%s", local.prefix, "google/admin")
  description             = "staging環境用のFirebase(admin)接続用API Key"
  recovery_window_in_days = 0
}
