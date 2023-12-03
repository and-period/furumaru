ALTER TABLE `stores`.`orders` ADD COLUMN `shipping_message` TEXT NULL DEFAULT NULL;
ALTER TABLE `stores`.`orders` ADD COLUMN `completed_at` DATETIME(3) NULL DEFAULT NULL;
