ALTER TABLE `users`.`address_revisions` MODIFY COLUMN `created_at` DATETIME(3) NOT NULL;
ALTER TABLE `users`.`address_revisions` MODIFY COLUMN `updated_at` DATETIME(3) NOT NULL;

ALTER TABLE `users`.`addresses` MODIFY COLUMN `created_at` DATETIME(3) NOT NULL;
ALTER TABLE `users`.`addresses` MODIFY COLUMN `updated_at` DATETIME(3) NOT NULL;
ALTER TABLE `users`.`addresses` MODIFY COLUMN `deleted_at` DATETIME(3) NULL DEFAULT NULL;
