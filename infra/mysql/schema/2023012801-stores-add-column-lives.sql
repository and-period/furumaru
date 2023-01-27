ALTER TABLE `stores`.`lives` DROP COLUMN `published`;
ALTER TABLE `stores`.`lives` DROP COLUMN `canceled`;

ALTER TABLE `stores`.`lives` ADD COLUMN `status` VARCHAR(16) NOT NULL;
