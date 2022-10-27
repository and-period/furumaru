ALTER TABLE `users`.`admins` ADD COLUMN `deleted_at` DATETIME NULL DEFAULT NULL;
ALTER TABLE `users`.`admins` ADD COLUMN `exists` TINYINT NULL DEFAULT 1 AFTER `device`;

DROP INDEX `ui_admins_email` ON `users`.`admins`;
CREATE UNIQUE INDEX `ui_admins_email` ON `users`.`admins` (`exists` DESC, `email` ASC) VISIBLE;
