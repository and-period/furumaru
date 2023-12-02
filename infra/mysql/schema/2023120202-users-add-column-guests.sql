ALTER TABLE `users`.`guests` ADD COLUMN `deleted_at` DATETIME(3) NULL DEFAULT NULL;
ALTER TABLE `users`.`guests` ADD COLUMN `exists` TINYINT NULL DEFAULT 1;

DROP INDEX `ui_guests_email` ON `users`.`guests`;
CREATE UNIQUE INDEX `ui_guests_email` ON `users`.`guests` (`exists` DESC, `email` ASC) VISIBLE;
