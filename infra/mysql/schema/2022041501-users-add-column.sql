ALTER TABLE `users`.`users` ADD COLUMN `account_id` VARCHAR(32) NOT NULL;
ALTER TABLE `users`.`users` ADD COLUMN `username`   VARCHAR(32) NOT NULL;

CREATE UNIQUE INDEX `ui_users_account_id` ON `users`.`users` (`account_id` ASC) VISIBLE;
