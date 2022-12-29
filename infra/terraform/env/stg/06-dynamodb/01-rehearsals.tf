module "rehearsals" {
  source = "./../../../modules/dynamodb"

  #####################################################################
  # Common
  #####################################################################
  tags = var.tags

  #####################################################################
  # DynamoDB
  #####################################################################
  name           = format("%s_%s", local.table_prefix, "rehearsals")
  billing_mode   = "PROVISIONED"
  table_class    = "STANDARD"
  read_capacity  = 1
  write_capacity = 1

  hash_key           = "live_id"
  range_key          = ""

  ttl_enabled        = true
  ttl_attribute_name = "expires_at"

  stream_enabled = true
  stream_view_type = "NEW_AND_OLD_IMAGES"

  attributes = [
    { name = "live_id", type = "S" },    # ライブ配信ID
  ]
}
