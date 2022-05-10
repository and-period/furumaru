ALTER TABLE `users`.`admins` ADD COLUMN `lastname` VARCHAR(16) NOT NULL AFTER `cognito_id`;
ALTER TABLE `users`.`admins` ADD COLUMN `firstname` VARCHAR(16) NOT NULL AFTER `lastname`;
ALTER TABLE `users`.`admins` ADD COLUMN `lastname_kana` VARCHAR(32) NOT NULL AFTER `firstname`;
ALTER TABLE `users`.`admins` ADD COLUMN `firstname_kana` VARCHAR(32) NOT NULL AFTER `lastname_kana`;
ALTER TABLE `users`.`admins` ADD COLUMN `thumbnail_url` TEXT NOT NULL AFTER `firstname_kana`;
