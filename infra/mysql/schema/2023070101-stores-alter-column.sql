ALTER TABLE `stores`.`schedules` DROP COLUMN `canceled`;

ALTER TABLE `stores`.`schedules` ADD COLUMN `public` TINYINT NOT NULL;
ALTER TABLE `stores`.`schedules` ADD COLUMN `approved` TINYINT NOT NULL;
ALTER TABLE `stores`.`schedules` ADD COLUMN `approved_admin_id` VARCHAR(22) NOT NULL;
ALTER TABLE `stores`.`schedules` ADD COLUMN `opening_video_url` TEXT NOT NULL;
ALTER TABLE `stores`.`schedules` ADD COLUMN `intermission_video_url` TEXT NOT NULL;
