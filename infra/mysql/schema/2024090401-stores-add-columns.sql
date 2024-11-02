ALTER TABLE `stores`.`experiences` ADD COLUMN `duration` INT NOT NULL;
ALTER TABLE `stores`.`experiences` ADD COLUMN `direction` TEXT NOT NULL;
ALTER TABLE `stores`.`experiences` ADD COLUMN `business_open_time` VARCHAR(4) NOT NULL;
ALTER TABLE `stores`.`experiences` ADD COLUMN `business_close_time` VARCHAR(4) NOT NULL;
