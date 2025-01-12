CREATE TABLE IF NOT EXISTS `users`.`admin_policies` (
  `id`          VARCHAR(64)  NOT NULL,
  `name`        VARCHAR(64)  NOT NULL,
  `description` TEXT         NOT NULL,
  `priority`    BIGINT       NOT NULL,
  `path`        TEXT         NOT NULL,
  `method`      TEXT         NOT NULL,
  `action`      VARCHAR(16)  NOT NULL,
  `created_at`  DATETIME(3)  NOT NULL,
  `updated_at`  DATETIME(3)  NOT NULL,
  PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `users`.`admin_roles` (
  `id`          VARCHAR(64)  NOT NULL,
  `name`        VARCHAR(64)  NOT NULL,
  `description` TEXT         NOT NULL,
  `created_at`  DATETIME(3)  NOT NULL,
  `updated_at`  DATETIME(3)  NOT NULL,
  PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `users`.`admin_role_policies` (
  `role_id`    VARCHAR(64) NOT NULL,
  `policy_id`  VARCHAR(64) NOT NULL,
  `created_at` DATETIME(3) NOT NULL,
  `updated_at` DATETIME(3) NOT NULL,
  PRIMARY KEY (`role_id`, `policy_id`),
  CONSTRAINT `fk_admin_role_policies_role_id` FOREIGN KEY (`role_id`) REFERENCES `admin_roles` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_admin_role_policies_policy_id` FOREIGN KEY (`policy_id`) REFERENCES `admin_policies` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS `users`.`admin_groups` (
  `id`               VARCHAR(64)  NOT NULL,
  `type`             INT          NOT NULL,
  `name`             VARCHAR(64)  NOT NULL,
  `description`      TEXT         NOT NULL,
  `created_admin_id` VARCHAR(22)  NULL DEFAULT NULL,
  `updated_admin_id` VARCHAR(22)  NULL DEFAULT NULL,
  `created_at`       DATETIME(3)  NOT NULL,
  `updated_at`       DATETIME(3)  NOT NULL,
  `deleted_at`       DATETIME(3)  NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_admin_groups_created_admin_id` FOREIGN KEY (`created_admin_id`) REFERENCES `admins` (`id`) ON DELETE CASCADE ON UPDATE SET NULL,
  CONSTRAINT `fk_admin_groups_updated_admin_id` FOREIGN KEY (`updated_admin_id`) REFERENCES `admins` (`id`) ON DELETE CASCADE ON UPDATE SET NULL
);

CREATE TABLE IF NOT EXISTS `users`.`admin_group_roles` (
  `group_id`         VARCHAR(64) NOT NULL,
  `role_id`          VARCHAR(64) NOT NULL,
  `created_admin_id` VARCHAR(22) NULL DEFAULT NULL,
  `updated_admin_id` VARCHAR(22) NULL DEFAULT NULL,
  `created_at`       DATETIME(3) NOT NULL,
  `updated_at`       DATETIME(3) NOT NULL,
  `deleted_at`       DATETIME(3) NULL DEFAULT NULL,
  PRIMARY KEY (`group_id`, `role_id`),
  CONSTRAINT `fk_admin_group_roles_group_id` FOREIGN KEY (`group_id`) REFERENCES `admin_groups` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_admin_group_roles_role_id` FOREIGN KEY (`role_id`) REFERENCES `admin_roles` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_admin_group_roles_created_admin_id` FOREIGN KEY (`created_admin_id`) REFERENCES `admins` (`id`) ON DELETE CASCADE ON UPDATE SET NULL,
  CONSTRAINT `fk_admin_group_roles_updated_admin_id` FOREIGN KEY (`updated_admin_id`) REFERENCES `admins` (`id`) ON DELETE CASCADE ON UPDATE SET NULL
);

CREATE TABLE IF NOT EXISTS `users`.`admin_group_users` (
  `group_id`         VARCHAR(64) NOT NULL,
  `admin_id`         VARCHAR(22) NOT NULL,
  `created_admin_id` VARCHAR(22) NULL DEFAULT NULL,
  `updated_admin_id` VARCHAR(22) NULL DEFAULT NULL,
  `expired_at`       DATETIME(3) NULL DEFAULT NULL,
  `created_at`       DATETIME(3) NOT NULL,
  `updated_at`       DATETIME(3) NOT NULL,
  PRIMARY KEY (`group_id`, `user_id`),
  CONSTRAINT `fk_admin_group_users_group_id` FOREIGN KEY (`group_id`) REFERENCES `admin_groups` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_admin_group_users_admin_id` FOREIGN KEY (`admin_id`) REFERENCES `admins` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_admin_group_users_created_admin_id` FOREIGN KEY (`created_admin_id`) REFERENCES `admins` (`id`) ON DELETE CASCADE ON UPDATE SET NULL,
  CONSTRAINT `fk_admin_group_users_updated_admin_id` FOREIGN KEY (`updated_admin_id`) REFERENCES `admins` (`id`) ON DELETE CASCADE ON UPDATE SET NULL
);
