ALTER TABLE `stores`.`experiences` DROP FOREIGN KEY `fk_experiences_shop_id`;
ALTER TABLE `stores`.`orders` DROP FOREIGN KEY `fk_orders_shop_id`;
ALTER TABLE `stores`.`products` DROP FOREIGN KEY `fk_products_shop_id`;
ALTER TABLE `stores`.`schedules` DROP FOREIGN KEY `fk_schedules_shop_id`;
ALTER TABLE `stores`.`shippings` DROP FOREIGN KEY `fk_shippings_shop_id`;
