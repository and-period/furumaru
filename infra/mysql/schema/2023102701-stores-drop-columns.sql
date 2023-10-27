ALTER TABLE `stores`.`schedules` DROP FOREIGN KEY `fk_schedules_shipping_id`;

ALTER TABLE `stores`.`schedules` DROP COLUMN `shipping_id`;
ALTER TABLE `stores`.`products` DROP COLUMN `business_days`;
