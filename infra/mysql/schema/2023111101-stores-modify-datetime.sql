ALTER TABLE `stores`.`live_products` MODIFY COLUMN `deleted_at` DATETIME(3) NULL DEFAULT NULL;

ALTER TABLE `stores`.`order_fulfillments` MODIFY COLUMN `shipped_at` DATETIME(3) NOT NULL;
ALTER TABLE `stores`.`order_fulfillments` MODIFY COLUMN `created_at` DATETIME(3) NOT NULL;
ALTER TABLE `stores`.`order_fulfillments` MODIFY COLUMN `updated_at` DATETIME(3) NOT NULL;

ALTER TABLE `stores`.`order_items` MODIFY COLUMN `created_at` DATETIME(3) NOT NULL;
ALTER TABLE `stores`.`order_items` MODIFY COLUMN `updated_at` DATETIME(3) NOT NULL;

ALTER TABLE `stores`.`order_payments` MODIFY COLUMN `ordered_at` DATETIME(3) NULL DEFAULT NULL;
ALTER TABLE `stores`.`order_payments` MODIFY COLUMN `paid_at` DATETIME(3) NULL DEFAULT NULL;
ALTER TABLE `stores`.`order_payments` MODIFY COLUMN `captured_at` DATETIME(3) NULL DEFAULT NULL;
ALTER TABLE `stores`.`order_payments` MODIFY COLUMN `failed_at` DATETIME(3) NULL DEFAULT NULL;
ALTER TABLE `stores`.`order_payments` MODIFY COLUMN `refunded_at` DATETIME(3) NULL DEFAULT NULL;
ALTER TABLE `stores`.`order_payments` MODIFY COLUMN `created_at` DATETIME(3) NOT NULL;
ALTER TABLE `stores`.`order_payments` MODIFY COLUMN `updated_at` DATETIME(3) NOT NULL;

ALTER TABLE `stores`.`orders` MODIFY COLUMN `created_at` DATETIME(3) NOT NULL;
ALTER TABLE `stores`.`orders` MODIFY COLUMN `updated_at` DATETIME(3) NOT NULL;
ALTER TABLE `stores`.`orders` MODIFY COLUMN `deleted_at` DATETIME(3) NULL DEFAULT NULL;

ALTER TABLE `stores`.`product_revisions` MODIFY COLUMN `created_at` DATETIME(3) NOT NULL;
ALTER TABLE `stores`.`product_revisions` MODIFY COLUMN `updated_at` DATETIME(3) NOT NULL;

ALTER TABLE `stores`.`shipping_revisions` MODIFY COLUMN `created_at` DATETIME(3) NOT NULL;
ALTER TABLE `stores`.`shipping_revisions` MODIFY COLUMN `updated_at` DATETIME(3) NOT NULL;

DROP TABLE IF EXISTS `stores`.`stripe_users`;
