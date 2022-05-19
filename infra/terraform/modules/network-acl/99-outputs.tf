##################################################
# Network ACL
##################################################
output "network_acls" {
  description = "Network ACL"
  value       = aws_network_acl.this
}
