ALTER TABLE `messengers`.`notifications` ADD COLUMN `deleted_at` DATETIME NULL DEFAULT NULL AFTER `updated_at`;
