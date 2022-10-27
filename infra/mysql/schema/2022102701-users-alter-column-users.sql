ALTER TABLE `users`.`users` ADD COLUMN `deleted_at` DATETIME NULL DEFAULT NULL;
ALTER TABLE `users`.`members` ADD COLUMN `exists` TINYINT NULL DEFAULT 1 AFTER `thumbnail_url`;

DROP INDEX `ui_users_account_id` ON `users`.`members`;
DROP INDEX `ui_users_provider_type_email` ON `users`.`members`;
DROP INDEX `ui_users_provider_type_phone_number` ON `users`.`members`;

CREATE UNIQUE INDEX `ui_members_account_id` ON `users`.`members` (`exists` DESC, `account_id` ASC) VISIBLE;
CREATE UNIQUE INDEX `ui_members_provider_type_email` ON `users`.`members` (`exists` DESC, `provider_type` ASC, `email` ASC) VISIBLE;
CREATE UNIQUE INDEX `ui_members_provider_type_phone_number` ON `users`.`members` (`exists` DESC, `provider_type` ASC, `phone_number` ASC) VISIBLE;

ALTER TABLE `users`.`customers` DROP FOREIGN KEY `fk_customer_user_id`;
ALTER TABLE `users`.`customers` ADD CONSTRAINT `fk_customers_user_id`
  FOREIGN KEY (`user_id`) REFERENCES `users`.`users` (`id`)
  ON DELETE CASCADE ON UPDATE CASCADE;
