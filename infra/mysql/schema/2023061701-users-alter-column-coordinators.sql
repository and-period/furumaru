ALTER TABLE `users`.`coordinators` ADD COLUMN `marche_name` VARCHAR(64) NOT NULL;
ALTER TABLE `users`.`coordinators` ADD COLUMN `username` VARCHAR(64) NOT NULL;
ALTER TABLE `users`.`coordinators` ADD COLUMN `profile` TEXT NOT NULL;
ALTER TABLE `users`.`coordinators` ADD COLUMN `product_type_ids` JSON NULL DEFAULT NULL;
ALTER TABLE `users`.`coordinators` ADD COLUMN `promotion_video_url` TEXT NOT NULL;
ALTER TABLE `users`.`coordinators` ADD COLUMN `bonus_video_url` TEXT NOT NULL;
ALTER TABLE `users`.`coordinators` ADD COLUMN `instagram_id` VARCHAR(30) NOT NULL;
ALTER TABLE `users`.`coordinators` ADD COLUMN `facebook_id` VARCHAR(50) NOT NULL;

ALTER TABLE `users`.`coordinators` DROP COLUMN `company_name`;
ALTER TABLE `users`.`coordinators` DROP COLUMN `store_name`;
ALTER TABLE `users`.`coordinators` DROP COLUMN `twitter_account`;
ALTER TABLE `users`.`coordinators` DROP COLUMN `instagram_account`;
ALTER TABLE `users`.`coordinators` DROP COLUMN `facebook_account`;
