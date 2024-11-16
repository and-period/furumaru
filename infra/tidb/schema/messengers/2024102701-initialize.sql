CREATE SCHEMA IF NOT EXISTS `messengers` DEFAULT CHARACTER SET utf8mb4;

CREATE TABLE IF NOT EXISTS `messengers`.`contact_categories` (
  `id` varchar(22) NOT NULL,
  `title` varchar(64) NOT NULL,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `title` (`title`)
);

CREATE TABLE IF NOT EXISTS `messengers`.`contacts` (
  `id` varchar(22) NOT NULL,
  `category_id` varchar(22) DEFAULT NULL,
  `title` varchar(64) NOT NULL,
  `content` text NOT NULL,
  `username` varchar(64) NOT NULL,
  `user_id` varchar(22) DEFAULT NULL,
  `email` varchar(256) NOT NULL,
  `phone_number` varchar(18) NOT NULL,
  `status` int NOT NULL,
  `responder_id` varchar(22) DEFAULT NULL,
  `note` text NOT NULL,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_contacts_caterogy_id` (`category_id`),
  CONSTRAINT `fk_contacts_caterogy_id` FOREIGN KEY (`category_id`) REFERENCES `contact_categories` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS `messengers`.`contact_reads` (
  `id` varchar(22) NOT NULL,
  `contact_id` varchar(22) DEFAULT NULL,
  `user_id` varchar(22) DEFAULT NULL,
  `user_type` int NOT NULL,
  `read` tinyint NOT NULL,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_contact_reads_contact_id` (`contact_id`),
  CONSTRAINT `fk_contact_reads_contact_id` FOREIGN KEY (`contact_id`) REFERENCES `contacts` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS `messengers`.`threads` (
  `id` varchar(22) NOT NULL,
  `contact_id` varchar(22) DEFAULT NULL,
  `user_id` varchar(22) DEFAULT NULL,
  `user_type` int NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_threads_contact_id` (`contact_id`),
  CONSTRAINT `fk_threads_contact_id` FOREIGN KEY (`contact_id`) REFERENCES `contacts` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS `messengers`.`notifications` (
  `id` varchar(22) NOT NULL,
  `created_by` varchar(22) NOT NULL,
  `updated_by` varchar(22) NOT NULL,
  `title` varchar(128) NOT NULL,
  `body` text NOT NULL,
  `published_at` datetime(3) NOT NULL,
  `targets` json DEFAULT NULL,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `type` int NOT NULL,
  `note` text NOT NULL,
  `promotion_id` varchar(22) DEFAULT NULL,
  PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `messengers`.`push_templates` (
  `id` varchar(64) NOT NULL,
  `title_template` text NOT NULL,
  `body_template` text NOT NULL,
  `image_url` text NOT NULL,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `messengers`.`report_templates` (
  `id` varchar(64) NOT NULL,
  `template` text NOT NULL,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `messengers`.`message_templates` (
  `id` varchar(64) NOT NULL,
  `title_template` text NOT NULL,
  `body_template` text NOT NULL,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `messengers`.`messages` (
  `id` varchar(22) NOT NULL,
  `user_type` int NOT NULL,
  `user_id` varchar(22) NOT NULL,
  `type` int NOT NULL,
  `title` varchar(256) NOT NULL,
  `body` text NOT NULL,
  `link` text NOT NULL,
  `read` tinyint NOT NULL,
  `received_at` datetime(3) NOT NULL,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_messages_user_type_user_id` (`user_type`,`user_id`)
);

CREATE TABLE IF NOT EXISTS `messengers`.`schedules` (
  `message_type` int NOT NULL,
  `message_id` varchar(22) NOT NULL,
  `status` int NOT NULL,
  `count` bigint NOT NULL,
  `sent_at` datetime NOT NULL,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  `deadline` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`message_type`,`message_id`)
);

CREATE TABLE IF NOT EXISTS `messengers`.`received_queues` (
  `id` varchar(22) NOT NULL,
  `event_type` int NOT NULL,
  `user_type` int NOT NULL,
  `user_ids` json DEFAULT NULL,
  `done` tinyint NOT NULL,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  `notify_type` int NOT NULL,
  PRIMARY KEY (`id`,`notify_type`)
);
