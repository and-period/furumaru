ALTER TABLE `users`.`producers` ADD COLUMN `username` VARCHAR(64) NOT NULL;
ALTER TABLE `users`.`producers` ADD COLUMN `profile` TEXT NOT NULL;
ALTER TABLE `users`.`producers` ADD COLUMN `promotion_video_url` TEXT NOT NULL;
ALTER TABLE `users`.`producers` ADD COLUMN `bonus_video_url` TEXT NOT NULL;
ALTER TABLE `users`.`producers` ADD COLUMN `instagram_id` VARCHAR(30) NOT NULL;
ALTER TABLE `users`.`producers` ADD COLUMN `facebook_id` VARCHAR(50) NOT NULL;

ALTER TABLE `users`.`producers` DROP COLUMN `store_name`;
