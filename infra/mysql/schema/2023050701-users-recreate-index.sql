ALTER TABLE `users`.`guests` DROP INDEX `ui_guests_email_phone_number`;

CREATE UNIQUE INDEX `ui_guests_email` ON `users`.`guests` (`email` ASC) VISIBLE;
