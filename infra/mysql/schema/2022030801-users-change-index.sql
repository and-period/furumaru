ALTER TABLE `users`.`users` DROP INDEX `ui_users_email`;
ALTER TABLE `users`.`users` DROP INDEX `ui_users_phone_number`;

CREATE UNIQUE INDEX `ui_users_provider_type_email` ON `users`.`users` (`provider_type` ASC, `email` ASC) VISIBLE;
CREATE UNIQUE INDEX `ui_users_provider_type_phone_number` ON `users`.`users` (`provider_type` ASC, `phone_number` ASC) VISIBLE;
