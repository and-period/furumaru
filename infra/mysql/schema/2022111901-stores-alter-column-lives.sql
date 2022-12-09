ALTER TABLE `stores`.`lives` ADD COLUMN `stream_key_arn` VARCHAR(256) NULL DEFAULT NULL AFTER `channel_arn`;
