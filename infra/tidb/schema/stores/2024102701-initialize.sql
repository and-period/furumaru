CREATE SCHEMA IF NOT EXISTS `stores` DEFAULT CHARACTER SET utf8mb4;

CREATE TABLE IF NOT EXISTS `stores`.`payment_systems` (
  `method_type` int NOT NULL,
  `status` int NOT NULL,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  PRIMARY KEY (`method_type`)
);

CREATE TABLE IF NOT EXISTS `stores`.`categories` (
  `id` varchar(22) NOT NULL,
  `name` varchar(32) NOT NULL,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `ui_categories_name` (`name`)
);

CREATE TABLE IF NOT EXISTS `stores`.`product_types` (
  `id` varchar(22) NOT NULL,
  `category_id` varchar(22) NOT NULL,
  `name` varchar(32) NOT NULL,
  `icon_url` text NOT NULL,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `ui_product_types_category_id_name` (`category_id`,`name`),
  CONSTRAINT `fk_product_types_category_id` FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS `stores`.`product_tags` (
  `id` varchar(22) NOT NULL,
  `name` varchar(32) NOT NULL,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `ui_product_tags_name` (`name`)
);

CREATE TABLE IF NOT EXISTS `stores`.`products` (
  `id` varchar(22) NOT NULL,
  `producer_id` varchar(22) NOT NULL,
  `product_type_id` varchar(22) DEFAULT NULL,
  `name` varchar(128) NOT NULL,
  `description` text NOT NULL,
  `public` tinyint NOT NULL,
  `inventory` bigint NOT NULL,
  `weight` bigint NOT NULL,
  `weight_unit` int NOT NULL,
  `item` bigint NOT NULL,
  `item_unit` varchar(16) NOT NULL,
  `item_description` varchar(64) NOT NULL,
  `media` json DEFAULT NULL,
  `delivery_type` int NOT NULL,
  `box60_rate` bigint NOT NULL,
  `box80_rate` bigint NOT NULL,
  `box100_rate` bigint NOT NULL,
  `origin_prefecture` bigint NOT NULL,
  `origin_city` varchar(32) NOT NULL,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `product_tag_ids` json DEFAULT NULL,
  `recommended_points` json DEFAULT NULL,
  `expiration_date` bigint NOT NULL,
  `storage_method_type` int NOT NULL,
  `start_at` datetime(3) NOT NULL,
  `end_at` datetime(3) NOT NULL,
  `coordinator_id` varchar(22) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_products_product_type_id` (`product_type_id`),
  CONSTRAINT `fk_products_product_type_id` FOREIGN KEY (`product_type_id`) REFERENCES `product_types` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS `stores`.`product_revisions` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `product_id` varchar(22) NOT NULL,
  `price` bigint NOT NULL,
  `cost` bigint NOT NULL,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_product_revisions_product_id` (`product_id`),
  CONSTRAINT `fk_product_revisions_product_id` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS `stores`.`experience_types` (
  `id` varchar(22) NOT NULL,
  `name` varchar(32) NOT NULL,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `ui_experience_types_name` (`name`)
);

CREATE TABLE IF NOT EXISTS `stores`.`experiences` (
  `id` varchar(22) NOT NULL,
  `coordinator_id` varchar(22) NOT NULL,
  `producer_id` varchar(22) NOT NULL,
  `experience_type_id` varchar(22) DEFAULT NULL,
  `title` varchar(128) NOT NULL,
  `description` text NOT NULL,
  `public` tinyint NOT NULL,
  `sold_out` tinyint NOT NULL,
  `media` json DEFAULT NULL,
  `recommended_points` json DEFAULT NULL,
  `promotion_video_url` text NOT NULL,
  `host_prefecture` bigint NOT NULL,
  `host_city` varchar(32) NOT NULL,
  `start_at` datetime(3) NOT NULL,
  `end_at` datetime(3) NOT NULL,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `host_postal_code` varchar(16) NOT NULL,
  `host_address_line1` varchar(64) NOT NULL,
  `host_address_line2` varchar(64) NOT NULL,
  `host_latitude` decimal(10, 6) NOT NULL,
  `host_longitude` decimal(10, 6) NOT NULL,
  `duration` int NOT NULL,
  `direction` text NOT NULL,
  `business_open_time` varchar(4) NOT NULL,
  `business_close_time` varchar(4) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_experiences_experience_type_id` (`experience_type_id`),
  CONSTRAINT `fk_experiences_experience_type_id` FOREIGN KEY (`experience_type_id`) REFERENCES `experience_types` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS `stores`.`experience_revisions` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `experience_id` varchar(22) NOT NULL,
  `price_adult` bigint NOT NULL,
  `price_junior_high_school` bigint NOT NULL,
  `price_elementary_school` bigint NOT NULL,
  `price_preschool` bigint NOT NULL,
  `price_senior` bigint NOT NULL,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_experience_revisions_experience_id` (`experience_id`),
  CONSTRAINT `fk_experience_revisions_experience_id` FOREIGN KEY (`experience_id`) REFERENCES `experiences` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS `stores`.`shippings` (
  `id` varchar(22) NOT NULL,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  `coordinator_id` varchar(22) NOT NULL,
  PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `stores`.`shipping_revisions` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `shipping_id` varchar(22) NOT NULL,
  `box60_rates` json DEFAULT NULL,
  `box60_frozen` bigint NOT NULL,
  `box80_rates` json DEFAULT NULL,
  `box80_frozen` bigint NOT NULL,
  `box100_rates` json DEFAULT NULL,
  `box100_frozen` bigint NOT NULL,
  `has_free_shipping` tinyint NOT NULL,
  `free_shipping_rates` bigint NOT NULL,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_shipping_revisions_shipping_id` (`shipping_id`),
  CONSTRAINT `fk_shipping_revisions_shipping_id` FOREIGN KEY (`shipping_id`) REFERENCES `shippings` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS `stores`.`schedules` (
  `id` varchar(22) NOT NULL,
  `coordinator_id` varchar(22) NOT NULL,
  `title` varchar(64) NOT NULL,
  `description` text NOT NULL,
  `thumbnail_url` text NOT NULL,
  `start_at` datetime(3) NOT NULL,
  `end_at` datetime(3) NOT NULL,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `public` tinyint NOT NULL,
  `approved` tinyint NOT NULL,
  `approved_admin_id` varchar(22) NOT NULL,
  `opening_video_url` text NOT NULL,
  `image_url` text NOT NULL,
  PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `stores`.`lives` (
  `id` varchar(22) NOT NULL,
  `schedule_id` varchar(22) NOT NULL,
  `producer_id` varchar(22) NOT NULL,
  `start_at` datetime(3) NOT NULL,
  `end_at` datetime(3) NOT NULL,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  `comment` text NOT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_lives_schedule_id` (`schedule_id`),
  CONSTRAINT `fk_lives_schedule_id` FOREIGN KEY (`schedule_id`) REFERENCES `schedules` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS `stores`.`live_products` (
  `live_id` varchar(22) NOT NULL,
  `product_id` varchar(22) NOT NULL,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `priority` bigint NOT NULL,
  PRIMARY KEY (`live_id`,`product_id`),
  KEY `fk_live_products_product_id` (`product_id`),
  CONSTRAINT `fk_live_products_live_id` FOREIGN KEY (`live_id`) REFERENCES `lives` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_live_products_product_id` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS `stores`.`promotions` (
  `id` varchar(22) NOT NULL,
  `title` varchar(64) NOT NULL,
  `description` text NOT NULL,
  `public` tinyint NOT NULL,
  `published_at` datetime(3) DEFAULT NULL,
  `discount_type` int NOT NULL,
  `discount_rate` bigint NOT NULL,
  `code` varchar(8) NOT NULL,
  `code_type` int NOT NULL,
  `start_at` datetime(3) NOT NULL,
  `end_at` datetime(3) NOT NULL,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `ui_promotions_code` (`code`)
);

CREATE TABLE IF NOT EXISTS `stores`.`orders` (
  `id` varchar(22) NOT NULL,
  `user_id` varchar(22) NOT NULL,
  `coordinator_id` varchar(22) NOT NULL,
  `promotion_id` varchar(22) DEFAULT NULL,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `shipping_message` text,
  `completed_at` datetime(3) DEFAULT NULL,
  `management_id` bigint NOT NULL,
  `status` int NOT NULL,
  `session_id` varchar(22) NOT NULL,
  `type` int NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `ui_orders_coordinator_id_management_id` (`coordinator_id`,`management_id` DESC),
  KEY `fk_orders_promotion_id` (`promotion_id`),
  CONSTRAINT `fk_orders_promotion_id` FOREIGN KEY (`promotion_id`) REFERENCES `promotions` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS `stores`.`order_payments` (
  `order_id` varchar(22) NOT NULL,
  `address_revision_id` bigint NOT NULL,
  `status` int NOT NULL,
  `transaction_id` varchar(256) NOT NULL,
  `method_type` int NOT NULL,
  `subtotal` bigint NOT NULL,
  `discount` bigint NOT NULL,
  `shipping_fee` bigint NOT NULL,
  `tax` bigint NOT NULL,
  `total` bigint NOT NULL,
  `refund_total` bigint NOT NULL,
  `refund_type` int NOT NULL,
  `refund_reason` text NOT NULL,
  `ordered_at` datetime(3) DEFAULT NULL,
  `paid_at` datetime(3) DEFAULT NULL,
  `captured_at` datetime(3) DEFAULT NULL,
  `failed_at` datetime(3) DEFAULT NULL,
  `refunded_at` datetime(3) DEFAULT NULL,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  `payment_id` varchar(256) DEFAULT NULL,
  `canceled_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`order_id`),
  CONSTRAINT `fk_order_payments_order_id` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS `stores`.`order_fulfillments` (
  `id` varchar(22) NOT NULL,
  `order_id` varchar(22) NOT NULL,
  `address_revision_id` bigint NOT NULL,
  `status` int NOT NULL,
  `tracking_number` varchar(32) DEFAULT NULL,
  `shipping_carrier` int NOT NULL,
  `shipping_type` int NOT NULL,
  `box_number` bigint NOT NULL,
  `box_size` int NOT NULL,
  `shipped_at` datetime(3) DEFAULT NULL,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  `box_rate` bigint NOT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_order_fulfillments_order_id` (`order_id`),
  CONSTRAINT `fk_order_fulfillments_order_id` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS `stores`.`order_items` (
  `fulfillment_id` varchar(22) NOT NULL,
  `product_revision_id` bigint NOT NULL,
  `order_id` varchar(22) NOT NULL,
  `quantity` bigint NOT NULL,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  PRIMARY KEY (`fulfillment_id`,`product_revision_id`),
  KEY `fk_order_items_order_id` (`order_id`),
  KEY `fk_order_items_product_revision_id` (`product_revision_id`),
  CONSTRAINT `fk_order_items_fulfillment_id` FOREIGN KEY (`fulfillment_id`) REFERENCES `order_fulfillments` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_order_items_order_id` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_order_items_product_revision_id` FOREIGN KEY (`product_revision_id`) REFERENCES `product_revisions` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS `stores`.`order_experiences` (
  `order_id` varchar(22) NOT NULL,
  `experience_revision_id` bigint NOT NULL,
  `adult_count` bigint NOT NULL,
  `junior_high_school_count` bigint NOT NULL,
  `elementary_school_count` bigint NOT NULL,
  `preschool_count` bigint NOT NULL,
  `senior_count` bigint NOT NULL,
  `remarks` json DEFAULT NULL,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  PRIMARY KEY (`order_id`,`experience_revision_id`),
  KEY `fk_order_experiences_experience_revision_id` (`experience_revision_id`),
  CONSTRAINT `fk_order_experiences_experience_revision_id` FOREIGN KEY (`experience_revision_id`) REFERENCES `experience_revisions` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_order_experiences_order_id` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
);
