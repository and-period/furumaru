ALTER TABLE `users`.`administrators` ADD COLUMN `deleted_at` DATETIME(3) NULL DEFAULT NULL;
ALTER TABLE `users`.`coordinators` ADD COLUMN `deleted_at` DATETIME(3) NULL DEFAULT NULL;
ALTER TABLE `users`.`producers` ADD COLUMN `deleted_at` DATETIME(3) NULL DEFAULT NULL;
