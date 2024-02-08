ALTER TABLE `messengers`.`received_queues` ADD COLUMN `notify_type` INT NOT NULL;
ALTER TABLE `messengers`.`received_queues` DROP PRIMARY KEY, ADD PRIMARY KEY (`id`, `notify_type`);
