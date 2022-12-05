ALTER TABLE `stores`.`schedules` ADD COLUMN `shipping_id` VARCHAR(22) NOT NULL AFTER `id`;
ALTER TABLE `stores`.`schedules` ADD CONSTRAINT `fk_schedules_shipping_id`
    FOREIGN KEY (`shipping_id`) REFERENCES `stores`.`shippings` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE;
