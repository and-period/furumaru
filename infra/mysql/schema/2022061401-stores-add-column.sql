ALTER TABLE `stores`.`products` ADD COLUMN `created_by` VARCHAR(22) NOT NULL AFTER `created_at`;
ALTER TABLE `stores`.`products` ADD COLUMN `updated_by` VARCHAR(22) NOT NULL AFTER `updated_at`;
