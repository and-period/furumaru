CREATE TABLE IF NOT EXISTS `stores`.`order_metadata` (
  `order_id`        VARCHAR(22) NOT NULL,
  `pickup_at`       DATETIME(3) NULL DEFAULT NULL,
  `pickup_location` TEXT        NULL DEFAULT NULL,
  `created_at`      DATETIME(3) NOT NULL,
  `updated_at`      DATETIME(3) NOT NULL,
  PRIMARY KEY (`order_id`),
  CONSTRAINT `fk_order_metadata_order_id` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
);
