ALTER TABLE `users`.`producers` ADD COLUMN `coordinator_id` VARCHAR(22) NOT NULL AFTER `admin_id`;
ALTER TABLE `users`.`producers` ADD CONSTRAINT `fk_producers_coordinator_id`
  FOREIGN KEY (`coordinator_id`) REFERENCES `users`.`coordinators` (`admin_id`)
  ON DELETE CASCADE ON UPDATE CASCADE;
