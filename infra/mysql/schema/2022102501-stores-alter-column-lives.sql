ALTER TABLE `stores`.`lives` ADD COLUMN `shipping_id` VARCHAR(22) NOT NULL AFTER `id`;
ALTER TABLE `stores`.`lives` ADD CONSTRAINT `fk_lives_shipping_id`
    FOREIGN KEY (`shipping_id`) REFERENCES `stores`.`shippings` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE;
