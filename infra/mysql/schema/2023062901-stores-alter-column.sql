ALTER TABLE `stores`.`products` ADD COLUMN `start_at` DATETIME NOT NULL;
ALTER TABLE `stores`.`products` ADD COLUMN `end_at` DATETIME NOT NULL;
ALTER TABLE `stores`.`products` ADD COLUMN `business_days` JSON NULL DEFAULT NULL;
