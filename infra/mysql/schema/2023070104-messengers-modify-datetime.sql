ALTER TABLE `messengers`.`contact_categories` MODIFY COLUMN `created_at` DATETIME(3) NOT NULL;
ALTER TABLE `messengers`.`contact_categories` MODIFY COLUMN `updated_at` DATETIME(3) NOT NULL;

ALTER TABLE `messengers`.`contact_reads` MODIFY COLUMN `created_at` DATETIME(3) NOT NULL;
ALTER TABLE `messengers`.`contact_reads` MODIFY COLUMN `updated_at` DATETIME(3) NOT NULL;

ALTER TABLE `messengers`.`contacts` MODIFY COLUMN `created_at` DATETIME(3) NOT NULL;
ALTER TABLE `messengers`.`contacts` MODIFY COLUMN `updated_at` DATETIME(3) NOT NULL;
ALTER TABLE `messengers`.`contacts` MODIFY COLUMN `deleted_at` DATETIME(3) NULL DEFAULT NULL;

ALTER TABLE `messengers`.`message_templates` MODIFY COLUMN `created_at` DATETIME(3) NOT NULL;
ALTER TABLE `messengers`.`message_templates` MODIFY COLUMN `updated_at` DATETIME(3) NOT NULL;

ALTER TABLE `messengers`.`messages` MODIFY COLUMN `created_at` DATETIME(3) NOT NULL;
ALTER TABLE `messengers`.`messages` MODIFY COLUMN `updated_at` DATETIME(3) NOT NULL;
ALTER TABLE `messengers`.`messages` MODIFY COLUMN `received_at` DATETIME(3) NOT NULL;

ALTER TABLE `messengers`.`notifications` MODIFY COLUMN `created_at` DATETIME(3) NOT NULL;
ALTER TABLE `messengers`.`notifications` MODIFY COLUMN `updated_at` DATETIME(3) NOT NULL;
ALTER TABLE `messengers`.`notifications` MODIFY COLUMN `deleted_at` DATETIME(3) NULL DEFAULT NULL;
ALTER TABLE `messengers`.`notifications` MODIFY COLUMN `published_at` DATETIME(3) NOT NULL;

ALTER TABLE `messengers`.`push_templates` MODIFY COLUMN `created_at` DATETIME(3) NOT NULL;
ALTER TABLE `messengers`.`push_templates` MODIFY COLUMN `updated_at` DATETIME(3) NOT NULL;

ALTER TABLE `messengers`.`received_queues` MODIFY COLUMN `created_at` DATETIME(3) NOT NULL;
ALTER TABLE `messengers`.`received_queues` MODIFY COLUMN `updated_at` DATETIME(3) NOT NULL;

ALTER TABLE `messengers`.`report_templates` MODIFY COLUMN `created_at` DATETIME(3) NOT NULL;
ALTER TABLE `messengers`.`report_templates` MODIFY COLUMN `updated_at` DATETIME(3) NOT NULL;

ALTER TABLE `messengers`.`schedules` MODIFY COLUMN `created_at` DATETIME(3) NOT NULL;
ALTER TABLE `messengers`.`schedules` MODIFY COLUMN `updated_at` DATETIME(3) NOT NULL;

ALTER TABLE `messengers`.`threads` MODIFY COLUMN `created_at` DATETIME(3) NOT NULL;
ALTER TABLE `messengers`.`threads` MODIFY COLUMN `updated_at` DATETIME(3) NOT NULL;
ALTER TABLE `messengers`.`threads` MODIFY COLUMN `deleted_at` DATETIME(3) NULL DEFAULT NULL;
