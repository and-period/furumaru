module "newrelic" {
  source = "./../../../modules/secrets-manager"

  #####################################################################
  # Common
  #####################################################################
  tags           = var.tags
  kms_alias_name = ""

  #####################################################################
  # Secrets Manager
  #####################################################################
  name                    = format("%s/%s", local.prefix, "newrelic")
  description             = "staging環境用のNew Relicへのメトリクス送信用"
  recovery_window_in_days = 0
}
