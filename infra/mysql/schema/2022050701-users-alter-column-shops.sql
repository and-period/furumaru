ALTER TABLE `users`.`shops` ADD COLUMN `lastname` VARCHAR(16) NOT NULL AFTER `name`;
ALTER TABLE `users`.`shops` ADD COLUMN `firstname` VARCHAR(16) NOT NULL AFTER `lastname`;
ALTER TABLE `users`.`shops` ADD COLUMN `lastname_kana` VARCHAR(32) NOT NULL AFTER `firstname`;
ALTER TABLE `users`.`shops` ADD COLUMN `firstname_kana` VARCHAR(32) NOT NULL AFTER `lastname_kana`;

ALTER TABLE `users`.`shops` DROP COLUMN `name`;
