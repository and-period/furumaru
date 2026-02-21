ALTER TABLE `users`.`producers` DROP FOREIGN KEY `fk_producers_coordinator_id`;
ALTER TABLE `users`.`producers` DROP COLUMN `coordinator_id`;
