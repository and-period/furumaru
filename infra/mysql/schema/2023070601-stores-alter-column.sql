ALTER TABLE `stores`.`schedules` DROP COLUMN `intermission_video_url`;

ALTER TABLE `stores`.`schedules` ADD COLUMN `thumbnails` JSON NULL DEFAULT NULL;
ALTER TABLE `stores`.`schedules` ADD COLUMN `image_url` TEXT NOT NULL;
ALTER TABLE `stores`.`schedules` ADD COLUMN `images` JSON NULL DEFAULT NULL;
