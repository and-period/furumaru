ALTER TABLE `stores`.`product_types` ADD COLUMN `icon_url` TEXT NOT NULL AFTER `name`;

ALTER TABLE `stores`.`products` DROP COLUMN `icon_url`;
