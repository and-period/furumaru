ALTER TABLE `messengers`.`notifications` ADD COLUMN `type` INT NOT NULL;
ALTER TABLE `messengers`.`notifications` ADD COLUMN `note` TEXT NOT NULL;
ALTER TABLE `messengers`.`notifications` ADD COLUMN `promotion_id` VARCHAR(22) NULL DEFAULT NULL;

ALTER TABLE `messengers`.`notifications` DROP COLUMN `public`;
ALTER TABLE `messengers`.`notifications` DROP COLUMN `creator_name`;
