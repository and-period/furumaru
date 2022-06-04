locals {
  subdomains = {
    for s in var.subdomains : s.domain => s
  }
}
