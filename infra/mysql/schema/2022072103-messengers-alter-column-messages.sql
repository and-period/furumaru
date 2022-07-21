ALTER TABLE `messengers`.`message_templates` ADD COLUMN `body_template` TEXT NOT NULL AFTER `template`;
ALTER TABLE `messengers`.`message_templates` ADD COLUMN `title_template` TEXT NOT NULL AFTER `template`;

ALTER TABLE `messengers`.`message_templates` DROP COLUMN `template`;
