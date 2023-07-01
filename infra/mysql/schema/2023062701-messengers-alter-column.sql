ALTER TABLE `messengers`.`contact_reads` DROP FOREIGN KEY `fk_contact_reads_contact_id`;
ALTER TABLE `messengers`.`contact_reads` DROP COLUMN `thread_id`;

ALTER TABLE `messengers`.`contact_reads` ADD COLUMN `contact_id` VARCHAR(22) NULL DEFAULT NULL AFTER `id`;
ALTER TABLE `messengers`.`contact_reads` ADD CONSTRAINT `fk_contact_reads_contact_id`
  FOREIGN KEY (`contact_id`) REFERENCES `messengers`.`contacts` (`id`)
  ON DELETE CASCADE ON UPDATE CASCADE;
