ALTER TABLE `stores`.`promotions` ADD COLUMN `shop_id` VARCHAR(22) NULL DEFAULT NULL;
ALTER TABLE `stores`.`promotions` ADD COLUMN `target_type` INT NOT NULL DEFAULT 0;

ALTER TABLE `stores`.`promotions` ADD CONSTRAINT `fk_promotions_shop_id`
  FOREIGN KEY (`shop_id`) REFERENCES `stores`.`shops` (`id`)
  ON DELETE SET NULL ON UPDATE CASCADE;
