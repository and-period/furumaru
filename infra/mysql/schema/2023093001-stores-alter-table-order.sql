ALTER TABLE `stores`.`orders` DROP COLUMN `cancel_type`;
ALTER TABLE `stores`.`orders` DROP COLUMN `cancel_reason`;
ALTER TABLE `stores`.`payments` DROP COLUMN `method_id`;

ALTER TABLE `stores`.`orders` ADD COLUMN `refund_reason` TEXT NOT NULL;
