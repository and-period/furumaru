ALTER TABLE `stores`.`orders` ADD COLUMN `management_id` BIGINT NULL DEFAULT NULL;

CREATE UNIQUE INDEX `ui_orders_coordinator_id_management_id` ON `users`.`orders` (`coordinator_id`, `management_id` DESC) VISIBLE;
