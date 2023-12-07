ALTER TABLE `users`.`address_revisions` MODIFY COLUMN `lastname` VARCHAR(16) NOT NULL;
ALTER TABLE `users`.`address_revisions` MODIFY COLUMN `firstname` VARCHAR(16) NOT NULL;

ALTER TABLE `users`.`members` ADD COLUMN `lastname` VARCHAR(16) NOT NULL;
ALTER TABLE `users`.`members` ADD COLUMN `firstname` VARCHAR(16) NOT NULL;
ALTER TABLE `users`.`members` ADD COLUMN `lastname_kana` VARCHAR(32) NOT NULL;
ALTER TABLE `users`.`members` ADD COLUMN `firstname_kana` VARCHAR(32) NOT NULL;
