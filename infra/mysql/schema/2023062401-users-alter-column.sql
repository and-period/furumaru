ALTER TABLE `users`.`coordinators` MODIFY COLUMN `prefecture` BIGINT NOT NULL;
ALTER TABLE `users`.`customers` MODIFY COLUMN `prefecture` BIGINT NOT NULL;
ALTER TABLE `users`.`producers` MODIFY COLUMN `prefecture` BIGINT NOT NULL;
