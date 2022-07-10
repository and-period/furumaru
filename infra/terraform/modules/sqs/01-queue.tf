##################################################
# VPC
# - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/sqs_queue
##################################################
resource "aws_sqs_queue" "this" {
  name = var.queue_name

  fifo_queue                  = var.fifo_queue
  content_based_deduplication = var.fifo_queue ? var.content_based_deduplication : null
  deduplication_scope         = var.fifo_queue ? var.deduplication_scope : null
  fifo_throughput_limit       = var.fifo_queue ? var.fifo_throughput_limit : null

  visibility_timeout_seconds = var.visibility_timeout_seconds
  message_retention_seconds  = var.message_retention_seconds
  delay_seconds              = var.delay_seconds
  max_message_size           = var.max_message_size
  receive_wait_time_seconds  = var.receive_wait_time_seconds

  tags = merge(
    var.tags,
    { Name = var.queue_name },
  )
}
