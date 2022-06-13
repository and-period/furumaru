ALTER TABLE `stores`.`products` ADD COLUMN `created_by` VARCHAR(22) NOT NULL AFTER `origin_city`;
ALTER TABLE `stores`.`products` ADD COLUMN `updated_by` VARCHAR(22) NOT NULL AFTER `created_at`;
