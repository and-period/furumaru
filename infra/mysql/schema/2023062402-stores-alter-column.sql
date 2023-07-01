ALTER TABLE `stores`.`addresses` MODIFY COLUMN `prefecture` BIGINT NOT NULL;

ALTER TABLE `stores`.`addresses` DROP COLUMN `prefecture_code`;
