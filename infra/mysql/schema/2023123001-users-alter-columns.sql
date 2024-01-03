ALTER TABLE `users`.`admins` MODIFY COLUMN `cognito_id` VARCHAR(36) NULL DEFAULT NULL;
ALTER TABLE `users`.`admins` MODIFY COLUMN `email` VARCHAR(256) NULL DEFAULT NULL;
ALTER TABLE `users`.`admins` MODIFY COLUMN `lastname` VARCHAR(16) NULL DEFAULT NULL;
ALTER TABLE `users`.`admins` MODIFY COLUMN `firstname` VARCHAR(16) NULL DEFAULT NULL;
ALTER TABLE `users`.`admins` MODIFY COLUMN `lastname_kana` VARCHAR(32) NULL DEFAULT NULL;
ALTER TABLE `users`.`admins` MODIFY COLUMN `firstname_kana` VARCHAR(32) NULL DEFAULT NULL;

ALTER TABLE `users`.`producers` MODIFY COLUMN `phone_number` VARCHAR(18) NULL DEFAULT NULL;
ALTER TABLE `users`.`producers` MODIFY COLUMN `postal_code` VARCHAR(16) NULL DEFAULT NULL;
ALTER TABLE `users`.`producers` MODIFY COLUMN `prefecture` BIGINT NULL DEFAULT NULL;
ALTER TABLE `users`.`producers` MODIFY COLUMN `city` VARCHAR(32) NULL DEFAULT NULL;
ALTER TABLE `users`.`producers` MODIFY COLUMN `address_line1` VARCHAR(64) NULL DEFAULT NULL;
ALTER TABLE `users`.`producers` MODIFY COLUMN `address_line2` VARCHAR(64) NULL DEFAULT NULL;