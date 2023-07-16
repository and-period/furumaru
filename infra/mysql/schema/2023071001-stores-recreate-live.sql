ALTER TABLE `stores`.`lives` DROP COLUMN `title`;
ALTER TABLE `stores`.`lives` DROP COLUMN `description`;
ALTER TABLE `stores`.`lives` DROP COLUMN `status`;
ALTER TABLE `stores`.`lives` DROP COLUMN `channel_arn`;
ALTER TABLE `stores`.`lives` DROP COLUMN `stream_key_arn`;
ALTER TABLE `stores`.`lives` DROP COLUMN `deleted_at`;

ALTER TABLE `stores`.`lives` ADD COLUMN `comment` TEXT NOT NULL;
