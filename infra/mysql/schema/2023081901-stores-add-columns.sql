ALTER TABLE `stores`.`broadcasts` ADD COLUMN `media_live_channel_id` VARCHAR(256) NULL DEFAULT NULL;
ALTER TABLE `stores`.`broadcasts` ADD COLUMN `media_live_rtmp_input_name` VARCHAR(256) NULL DEFAULT NULL;
ALTER TABLE `stores`.`broadcasts` ADD COLUMN `media_live_mp4_input_name` VARCHAR(256) NULL DEFAULT NULL;
