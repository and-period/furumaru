ALTER TABLE `users`.`producers` MODIFY COLUMN `coordinator_id` VARCHAR(22) NULL DEFAULT NULL;

ALTER TABLE `users`.`producers` DROP FOREIGN KEY `fk_producers_coordinator_id`;
ALTER TABLE `users`.`producers` ADD CONSTRAINT `fk_producers_coordinator_id`
  FOREIGN KEY (`coordinator_id`) REFERENCES `users`.`coordinators` (`admin_id`)
  ON DELETE SET NULL ON UPDATE CASCADE;
