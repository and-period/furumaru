CREATE TABLE IF NOT EXISTS `users`.`admin_auth_providers` (
  `admin_id`      VARCHAR(22)  NOT NULL,
  `provider_type` INT          NOT NULL,
  `account_id`    VARCHAR(64)  NOT NULL,
  `email`         VARCHAR(255) NOT NULL,
  `created_at`    DATETIME(3)  NOT NULL,
  `updated_at`    DATETIME(3)  NOT NULL,
  PRIMARY KEY (`admin_id`, `provider_type`),
  CONSTRAINT `fk_admin_auth_providers_admin_id` FOREIGN KEY (`admin_id`) REFERENCES `admins` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
);
