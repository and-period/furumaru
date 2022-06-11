locals {
  create_lifecycle_policy = var.lifecycle_policy_priority != 0 || var.lifecycle_policy_untagged_priority != 0

  untagged_rules = flatten(var.lifecycle_policy_untagged_priority == 0 ? [] : [{
    rulePriority: var.lifecycle_policy_untagged_priority,
    description: var.lifecycle_policy_untagged_description,
    selection: {
      tagStatus: "untagged",
      countType: "sinceImagePushed",
      countUnit: "days",
      countNumber: var.lifecycle_policy_untagged_retention_days,
    },
    action: {
      type: "expire",
    },
  }])

  all_rules = flatten(var.lifecycle_policy_priority == 0 ? [] : [{
    rulePriority: var.lifecycle_policy_priority,
    description: var.lifecycle_policy_description,
    selection: {
      tagStatus: "any",
      countType: "imageCountMoreThan",
      countNumber: var.lifecycle_policy_image_count,
    },
    action: {
      type: "expire",
    },
  }])

  lifecycle_rules = jsonencode({
    "rules": flatten([
      local.untagged_rules,
      local.all_rules,
    ])
  })
}
