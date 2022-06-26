module "messenger" {
  source = "./../../../modules/sqs"

  #####################################################################
  # Common
  #####################################################################
  tags = var.tags

  #####################################################################
  # SQS
  #####################################################################
  queue_name = "furumaru-stg-messenger"
  fifo_queue = false

  visibility_timeout_seconds = 30     # 30 sec
  message_retention_seconds  = 86400  # 1 day
  delay_seconds              = 0      # 0 sec
  max_message_size           = 262144 # 256 KiB
  receive_wait_time_seconds  = 0      # 0 sec
}
