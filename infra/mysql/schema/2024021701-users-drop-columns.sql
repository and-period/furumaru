
ALTER TABLE `users`.`guests` DROP COLUMN `phone_number`;

ALTER TABLE `users`.`guests` ADD COLUMN `lastname` VARCHAR(16) NOT NULL;
ALTER TABLE `users`.`guests` ADD COLUMN `firstname` VARCHAR(16) NOT NULL;
ALTER TABLE `users`.`guests` ADD COLUMN `lastname_kana` VARCHAR(32) NOT NULL;
ALTER TABLE `users`.`guests` ADD COLUMN `firstname_kana` VARCHAR(32) NOT NULL;
