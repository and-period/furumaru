ALTER TABLE `stores`.`orders` DROP COLUMN `shipping_type`;
ALTER TABLE `stores`.`order_fulfillments` MODIFY COLUMN `address_revision_id` INT NULL DEFAULT 0;
ALTER TABLE `stores`.`order_payments` MODIFY COLUMN `address_revision_id` INT NULL DEFAULT 0;
ALTER TABLE `stores`.`order_metadata` ADD COLUMN `shipping_message` TEXT NULL DEFAULT NULL;
