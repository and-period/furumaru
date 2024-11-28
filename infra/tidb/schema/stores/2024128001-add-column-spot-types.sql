ALTER TABLE `stores`.`spots` ADD COLUMN `postal_code` VARCHAR(16) NOT NULL DEFAULT '';
ALTER TABLE `stores`.`spots` ADD COLUMN `prefecture` BIGINT NOT NULL DEFAULT 0;
ALTER TABLE `stores`.`spots` ADD COLUMN `city` VARCHAR(32) NOT NULL DEFAULT '';
ALTER TABLE `stores`.`spots` ADD COLUMN `address_line1` VARCHAR(64) NOT NULL DEFAULT '';
ALTER TABLE `stores`.`spots` ADD COLUMN `address_line2` VARCHAR(64) NOT NULL DEFAULT '';
