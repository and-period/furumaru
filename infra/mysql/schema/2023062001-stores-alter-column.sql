ALTER TABLE `stores`.`products` ADD COLUMN `product_tag_ids` JSON NULL DEFAULT NULL;
ALTER TABLE `stores`.`products` ADD COLUMN `cost` BIGINT NOT NULL;
ALTER TABLE `stores`.`products` ADD COLUMN `recommended_points` JSON NULL DEFAULT NULL;
ALTER TABLE `stores`.`products` ADD COLUMN `expiration_date` BIGINT NOT NULL;
ALTER TABLE `stores`.`products` ADD COLUMN `storage_method_type` INT NOT NULL;

ALTER TABLE `stores`.`products` MODIFY COLUMN `origin_prefecture` BIGINT NOT NULL;
