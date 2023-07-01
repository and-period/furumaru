ALTER TABLE `users`.`administrators` MODIFY COLUMN `created_at` DATETIME(3) NOT NULL;
ALTER TABLE `users`.`administrators` MODIFY COLUMN `updated_at` DATETIME(3) NOT NULL;
ALTER TABLE `users`.`administrators` MODIFY COLUMN `deleted_at` DATETIME(3) NULL DEFAULT NULL;

ALTER TABLE `users`.`admins` MODIFY COLUMN `created_at` DATETIME(3) NOT NULL;
ALTER TABLE `users`.`admins` MODIFY COLUMN `updated_at` DATETIME(3) NOT NULL;
ALTER TABLE `users`.`admins` MODIFY COLUMN `deleted_at` DATETIME(3) NULL DEFAULT NULL;
ALTER TABLE `users`.`admins` MODIFY COLUMN `first_sign_in_at` DATETIME(3) NULL DEFAULT NULL;
ALTER TABLE `users`.`admins` MODIFY COLUMN `last_sign_in_at` DATETIME(3) NULL DEFAULT NULL;

ALTER TABLE `users`.`coordinators` MODIFY COLUMN `created_at` DATETIME(3) NOT NULL;
ALTER TABLE `users`.`coordinators` MODIFY COLUMN `updated_at` DATETIME(3) NOT NULL;
ALTER TABLE `users`.`coordinators` MODIFY COLUMN `deleted_at` DATETIME(3) NULL DEFAULT NULL;

ALTER TABLE `users`.`customers` MODIFY COLUMN `created_at` DATETIME(3) NOT NULL;
ALTER TABLE `users`.`customers` MODIFY COLUMN `updated_at` DATETIME(3) NOT NULL;

ALTER TABLE `users`.`guests` MODIFY COLUMN `created_at` DATETIME(3) NOT NULL;
ALTER TABLE `users`.`guests` MODIFY COLUMN `updated_at` DATETIME(3) NOT NULL;

ALTER TABLE `users`.`members` MODIFY COLUMN `created_at` DATETIME(3) NOT NULL;
ALTER TABLE `users`.`members` MODIFY COLUMN `updated_at` DATETIME(3) NOT NULL;
ALTER TABLE `users`.`members` MODIFY COLUMN `deleted_at` DATETIME(3) NULL DEFAULT NULL;
ALTER TABLE `users`.`members` MODIFY COLUMN `verified_at` DATETIME(3) NULL DEFAULT NULL;

ALTER TABLE `users`.`producers` MODIFY COLUMN `created_at` DATETIME(3) NOT NULL;
ALTER TABLE `users`.`producers` MODIFY COLUMN `updated_at` DATETIME(3) NOT NULL;
ALTER TABLE `users`.`producers` MODIFY COLUMN `deleted_at` DATETIME(3) NULL DEFAULT NULL;

ALTER TABLE `users`.`users` MODIFY COLUMN `created_at` DATETIME(3) NOT NULL;
ALTER TABLE `users`.`users` MODIFY COLUMN `updated_at` DATETIME(3) NOT NULL;
ALTER TABLE `users`.`users` MODIFY COLUMN `deleted_at` DATETIME(3) NULL DEFAULT NULL;
