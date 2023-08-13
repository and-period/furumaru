ALTER TABLE `stores`.`broadcasts` ADD COLUMN `cloud_front_distribution_arn` TEXT NULL DEFAULT NULL;
ALTER TABLE `stores`.`broadcasts` ADD COLUMN `media_live_channel_arn` TEXT NULL DEFAULT NULL;
ALTER TABLE `stores`.`broadcasts` ADD COLUMN `media_live_rtmp_input_arn` TEXT NULL DEFAULT NULL;
ALTER TABLE `stores`.`broadcasts` ADD COLUMN `media_live_mp4_input_arn` TEXT NULL DEFAULT NULL;
ALTER TABLE `stores`.`broadcasts` ADD COLUMN `media_store_container_arn` TEXT NULL DEFAULT NULL;
