-- 注文履歴テーブルの修正
ALTER TABLE `stores`.`orders` ADD COLUMN `coordinator_id` VARCHAR(22) NOT NULL AFTER `id`;
ALTER TABLE `stores`.`orders` ADD COLUMN `schedule_id` VARCHAR(22) NOT NULL AFTER `id`;
ALTER TABLE `stores`.`orders` ADD CONSTRAINT `fk_orders_schedule_id`
  FOREIGN KEY (`schedule_id`) REFERENCES `stores`.`schedules` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- 商品テーブルの修正
ALTER TABLE `stores`.`products` DROP COLUMN `created_by`;
ALTER TABLE `stores`.`products` DROP COLUMN `updated_by`;
