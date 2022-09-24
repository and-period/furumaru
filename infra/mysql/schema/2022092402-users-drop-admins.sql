DROP TABLE IF EXISTS `users`.`admin_auths`;

ALTER TABLE `users`.`administrators` DROP PRIMARY KEY, ADD PRIMARY KEY (`admin_id`);
ALTER TABLE `users`.`administrators` DROP COLUMN `id`;
ALTER TABLE `users`.`administrators` DROP COLUMN `lastname`;
ALTER TABLE `users`.`administrators` DROP COLUMN `firstname`;
ALTER TABLE `users`.`administrators` DROP COLUMN `lastname_kana`;
ALTER TABLE `users`.`administrators` DROP COLUMN `firstname_kana`;
ALTER TABLE `users`.`administrators` DROP COLUMN `email`;
ALTER TABLE `users`.`administrators` ADD CONSTRAINT `fk_administrators_admin_id`
  FOREIGN KEY (`admin_id`) REFERENCES `users`.`admins` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `users`.`coordinators` DROP PRIMARY KEY, ADD PRIMARY KEY (`admin_id`);
ALTER TABLE `users`.`coordinators` DROP COLUMN `id`;
ALTER TABLE `users`.`coordinators` DROP COLUMN `lastname`;
ALTER TABLE `users`.`coordinators` DROP COLUMN `firstname`;
ALTER TABLE `users`.`coordinators` DROP COLUMN `lastname_kana`;
ALTER TABLE `users`.`coordinators` DROP COLUMN `firstname_kana`;
ALTER TABLE `users`.`coordinators` DROP COLUMN `email`;
ALTER TABLE `users`.`coordinators` ADD CONSTRAINT `fk_coordinators_admin_id`
  FOREIGN KEY (`admin_id`) REFERENCES `users`.`admins` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `users`.`producers` DROP PRIMARY KEY, ADD PRIMARY KEY (`admin_id`);
ALTER TABLE `users`.`producers` DROP COLUMN `id`;
ALTER TABLE `users`.`producers` DROP COLUMN `lastname`;
ALTER TABLE `users`.`producers` DROP COLUMN `firstname`;
ALTER TABLE `users`.`producers` DROP COLUMN `lastname_kana`;
ALTER TABLE `users`.`producers` DROP COLUMN `firstname_kana`;
ALTER TABLE `users`.`producers` DROP COLUMN `email`;
ALTER TABLE `users`.`producers` ADD CONSTRAINT `fk_producers_admin_id`
  FOREIGN KEY (`admin_id`) REFERENCES `users`.`admins` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
