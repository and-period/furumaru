ALTER TABLE `stores`.`shippings` ADD COLUMN `is_default` TINYINT NOT NULL;
ALTER TABLE `stores`.`shippings` ADD COLUMN `coordinator_id` VARCHAR(22) NOT NULL;
