ALTER TABLE `stores`.`activities` MODIFY COLUMN `created_at` DATETIME(3) NOT NULL;
ALTER TABLE `stores`.`activities` MODIFY COLUMN `updated_at` DATETIME(3) NOT NULL;

ALTER TABLE `stores`.`addresses` MODIFY COLUMN `created_at` DATETIME(3) NOT NULL;
ALTER TABLE `stores`.`addresses` MODIFY COLUMN `updated_at` DATETIME(3) NOT NULL;
ALTER TABLE `stores`.`addresses` MODIFY COLUMN `deleted_at` DATETIME(3) NULL DEFAULT NULL;

ALTER TABLE `stores`.`cart_items` MODIFY COLUMN `created_at` DATETIME(3) NOT NULL;
ALTER TABLE `stores`.`cart_items` MODIFY COLUMN `updated_at` DATETIME(3) NOT NULL;

ALTER TABLE `stores`.`carts` MODIFY COLUMN `created_at` DATETIME(3) NOT NULL;
ALTER TABLE `stores`.`carts` MODIFY COLUMN `updated_at` DATETIME(3) NOT NULL;

ALTER TABLE `stores`.`categories` MODIFY COLUMN `created_at` DATETIME(3) NOT NULL;
ALTER TABLE `stores`.`categories` MODIFY COLUMN `updated_at` DATETIME(3) NOT NULL;

ALTER TABLE `stores`.`fulfillments` MODIFY COLUMN `created_at` DATETIME(3) NOT NULL;
ALTER TABLE `stores`.`fulfillments` MODIFY COLUMN `updated_at` DATETIME(3) NOT NULL;
ALTER TABLE `stores`.`fulfillments` MODIFY COLUMN `deleted_at` DATETIME(3) NULL DEFAULT NULL;

ALTER TABLE `stores`.`live_products` MODIFY COLUMN `created_at` DATETIME(3) NOT NULL;
ALTER TABLE `stores`.`live_products` MODIFY COLUMN `updated_at` DATETIME(3) NOT NULL;

ALTER TABLE `stores`.`lives` MODIFY COLUMN `created_at` DATETIME(3) NOT NULL;
ALTER TABLE `stores`.`lives` MODIFY COLUMN `updated_at` DATETIME(3) NOT NULL;
ALTER TABLE `stores`.`lives` MODIFY COLUMN `deleted_at` DATETIME(3) NULL DEFAULT NULL;
ALTER TABLE `stores`.`lives` MODIFY COLUMN `start_at` DATETIME(3) NOT NULL;
ALTER TABLE `stores`.`lives` MODIFY COLUMN `end_at` DATETIME(3) NOT NULL;

ALTER TABLE `stores`.`order_items` MODIFY COLUMN `created_at` DATETIME(3) NOT NULL;
ALTER TABLE `stores`.`order_items` MODIFY COLUMN `updated_at` DATETIME(3) NOT NULL;

ALTER TABLE `stores`.`orders` MODIFY COLUMN `created_at` DATETIME(3) NOT NULL;
ALTER TABLE `stores`.`orders` MODIFY COLUMN `updated_at` DATETIME(3) NOT NULL;
ALTER TABLE `stores`.`orders` MODIFY COLUMN `deleted_at` DATETIME(3) NULL DEFAULT NULL;
ALTER TABLE `stores`.`orders` MODIFY COLUMN `ordered_at` DATETIME(3) NULL DEFAULT NULL;
ALTER TABLE `stores`.`orders` MODIFY COLUMN `paid_at` DATETIME(3) NULL DEFAULT NULL;
ALTER TABLE `stores`.`orders` MODIFY COLUMN `captured_at` DATETIME(3) NULL DEFAULT NULL;
ALTER TABLE `stores`.`orders` MODIFY COLUMN `failed_at` DATETIME(3) NULL DEFAULT NULL;
ALTER TABLE `stores`.`orders` MODIFY COLUMN `refunded_at` DATETIME(3) NULL DEFAULT NULL;
ALTER TABLE `stores`.`orders` MODIFY COLUMN `shipped_at` DATETIME(3) NULL DEFAULT NULL;

ALTER TABLE `stores`.`payment_cards` MODIFY COLUMN `created_at` DATETIME(3) NOT NULL;
ALTER TABLE `stores`.`payment_cards` MODIFY COLUMN `updated_at` DATETIME(3) NOT NULL;

ALTER TABLE `stores`.`payments` MODIFY COLUMN `created_at` DATETIME(3) NOT NULL;
ALTER TABLE `stores`.`payments` MODIFY COLUMN `updated_at` DATETIME(3) NOT NULL;
ALTER TABLE `stores`.`payments` MODIFY COLUMN `deleted_at` DATETIME(3) NULL DEFAULT NULL;

ALTER TABLE `stores`.`product_tags` MODIFY COLUMN `created_at` DATETIME(3) NOT NULL;
ALTER TABLE `stores`.`product_tags` MODIFY COLUMN `updated_at` DATETIME(3) NOT NULL;

ALTER TABLE `stores`.`product_types` MODIFY COLUMN `created_at` DATETIME(3) NOT NULL;
ALTER TABLE `stores`.`product_types` MODIFY COLUMN `updated_at` DATETIME(3) NOT NULL;

ALTER TABLE `stores`.`products` MODIFY COLUMN `created_at` DATETIME(3) NOT NULL;
ALTER TABLE `stores`.`products` MODIFY COLUMN `updated_at` DATETIME(3) NOT NULL;
ALTER TABLE `stores`.`products` MODIFY COLUMN `deleted_at` DATETIME(3) NULL DEFAULT NULL;
ALTER TABLE `stores`.`products` MODIFY COLUMN `start_at` DATETIME(3) NOT NULL;
ALTER TABLE `stores`.`products` MODIFY COLUMN `end_at` DATETIME(3) NOT NULL;

ALTER TABLE `stores`.`promotions` MODIFY COLUMN `created_at` DATETIME(3) NOT NULL;
ALTER TABLE `stores`.`promotions` MODIFY COLUMN `updated_at` DATETIME(3) NOT NULL;
ALTER TABLE `stores`.`promotions` MODIFY COLUMN `published_at` DATETIME(3) NULL DEFAULT NULL;
ALTER TABLE `stores`.`promotions` MODIFY COLUMN `start_at` DATETIME(3) NOT NULL;
ALTER TABLE `stores`.`promotions` MODIFY COLUMN `end_at` DATETIME(3) NOT NULL;

ALTER TABLE `stores`.`schedules` MODIFY COLUMN `created_at` DATETIME(3) NOT NULL;
ALTER TABLE `stores`.`schedules` MODIFY COLUMN `updated_at` DATETIME(3) NOT NULL;
ALTER TABLE `stores`.`schedules` MODIFY COLUMN `deleted_at` DATETIME(3) NULL DEFAULT NULL;
ALTER TABLE `stores`.`schedules` MODIFY COLUMN `start_at` DATETIME(3) NOT NULL;
ALTER TABLE `stores`.`schedules` MODIFY COLUMN `end_at` DATETIME(3) NOT NULL;

ALTER TABLE `stores`.`shippings` MODIFY COLUMN `created_at` DATETIME(3) NOT NULL;
ALTER TABLE `stores`.`shippings` MODIFY COLUMN `updated_at` DATETIME(3) NOT NULL;

ALTER TABLE `stores`.`stripe_users` MODIFY COLUMN `created_at` DATETIME(3) NOT NULL;
ALTER TABLE `stores`.`stripe_users` MODIFY COLUMN `updated_at` DATETIME(3) NOT NULL;
