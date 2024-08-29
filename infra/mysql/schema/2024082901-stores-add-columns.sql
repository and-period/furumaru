ALTER TABLE `stores`.`experiences` ADD COLUMN `host_postal_code` VARCHAR(16) NOT NULL;
ALTER TABLE `stores`.`experiences` ADD COLUMN `host_address_line1` VARCHAR(64) NOT NULL;
ALTER TABLE `stores`.`experiences` ADD COLUMN `host_address_line2` VARCHAR(64) NOT NULL;
ALTER TABLE `stores`.`experiences` ADD COLUMN `host_geolocation` GEOMETRY NOT NULL;
