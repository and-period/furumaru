CREATE SCHEMA IF NOT EXISTS `users` DEFAULT CHARACTER SET utf8mb4;

CREATE TABLE IF NOT EXISTS `users`.`admins` (
  `id` varchar(22) NOT NULL,
  `cognito_id` varchar(36) DEFAULT NULL,
  `lastname` varchar(16) DEFAULT NULL,
  `firstname` varchar(16) DEFAULT NULL,
  `lastname_kana` varchar(32) DEFAULT NULL,
  `firstname_kana` varchar(32) DEFAULT NULL,
  `role` int NOT NULL,
  `email` varchar(256) DEFAULT NULL,
  `device` varchar(256) DEFAULT NULL,
  `exists` tinyint DEFAULT '1',
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `first_sign_in_at` datetime(3) DEFAULT NULL,
  `last_sign_in_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `ui_admins_cognito_id` (`cognito_id`),
  UNIQUE KEY `ui_admins_email` (`exists` DESC,`email`)
);

CREATE TABLE IF NOT EXISTS `users`.`administrators` (
  `admin_id` varchar(22) NOT NULL,
  `phone_number` varchar(18) NOT NULL,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  PRIMARY KEY (`admin_id`),
  CONSTRAINT `fk_administrators_admin_id` FOREIGN KEY (`admin_id`) REFERENCES `admins` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS `users`.`coordinators` (
  `admin_id` varchar(22) NOT NULL,
  `thumbnail_url` text NOT NULL,
  `header_url` text NOT NULL,
  `phone_number` varchar(18) NOT NULL,
  `postal_code` varchar(16) NOT NULL,
  `prefecture` bigint NOT NULL,
  `city` varchar(32) NOT NULL,
  `address_line1` varchar(64) NOT NULL,
  `address_line2` varchar(64) NOT NULL,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  `marche_name` varchar(64) NOT NULL,
  `username` varchar(64) NOT NULL,
  `profile` text NOT NULL,
  `product_type_ids` json DEFAULT NULL,
  `promotion_video_url` text NOT NULL,
  `bonus_video_url` text NOT NULL,
  `instagram_id` varchar(30) NOT NULL,
  `facebook_id` varchar(50) NOT NULL,
  `business_days` json DEFAULT NULL,
  PRIMARY KEY (`admin_id`),
  CONSTRAINT `fk_coordinators_admin_id` FOREIGN KEY (`admin_id`) REFERENCES `admins` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS `users`.`producers` (
  `admin_id` varchar(22) NOT NULL,
  `coordinator_id` varchar(22) DEFAULT NULL,
  `thumbnail_url` text NOT NULL,
  `header_url` text NOT NULL,
  `phone_number` varchar(18) DEFAULT NULL,
  `postal_code` varchar(16) DEFAULT NULL,
  `prefecture` bigint DEFAULT NULL,
  `city` varchar(32) DEFAULT NULL,
  `address_line1` varchar(64) DEFAULT NULL,
  `address_line2` varchar(64) DEFAULT NULL,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  `username` varchar(64) NOT NULL,
  `profile` text NOT NULL,
  `promotion_video_url` text NOT NULL,
  `bonus_video_url` text NOT NULL,
  `instagram_id` varchar(30) NOT NULL,
  `facebook_id` varchar(50) NOT NULL,
  PRIMARY KEY (`admin_id`),
  KEY `fk_producers_coordinator_id` (`coordinator_id`),
  CONSTRAINT `fk_producers_admin_id` FOREIGN KEY (`admin_id`) REFERENCES `admins` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_producers_coordinator_id` FOREIGN KEY (`coordinator_id`) REFERENCES `coordinators` (`admin_id`) ON DELETE SET NULL ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS `users`.`users` (
  `id` varchar(22) NOT NULL,
  `registered` tinyint NOT NULL,
  `device` varchar(256) NOT NULL,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `users`.`members` (
  `user_id` varchar(22) NOT NULL,
  `cognito_id` varchar(36) NOT NULL,
  `account_id` varchar(32) DEFAULT NULL,
  `username` varchar(32) NOT NULL,
  `provider_type` int NOT NULL,
  `email` varchar(256) DEFAULT NULL,
  `phone_number` varchar(18) DEFAULT NULL,
  `thumbnail_url` text NOT NULL,
  `exists` tinyint DEFAULT '1',
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  `verified_at` datetime(3) DEFAULT NULL,
  `lastname` varchar(16) NOT NULL,
  `firstname` varchar(16) NOT NULL,
  `lastname_kana` varchar(32) NOT NULL,
  `firstname_kana` varchar(32) NOT NULL,
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `ui_users_cognito_id` (`cognito_id`),
  UNIQUE KEY `ui_members_account_id` (`exists` DESC,`account_id`),
  UNIQUE KEY `ui_members_provider_type_email` (`exists` DESC,`provider_type`,`email`),
  CONSTRAINT `fk_accounts_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS `users`.`guests` (
  `user_id` varchar(22) NOT NULL,
  `email` varchar(256) NOT NULL,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  `exists` tinyint DEFAULT '1',
  `lastname` varchar(16) NOT NULL,
  `firstname` varchar(16) NOT NULL,
  `lastname_kana` varchar(32) NOT NULL,
  `firstname_kana` varchar(32) NOT NULL,
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `ui_guests_email` (`exists` DESC,`email`),
  CONSTRAINT `fk_guests_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS `users`.`user_notifications` (
  `user_id` varchar(22) NOT NULL,
  `disabled` tinyint NOT NULL,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  PRIMARY KEY (`user_id`),
  CONSTRAINT `fk_user_notifications_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS `users`.`addresses` (
  `id` varchar(22) NOT NULL,
  `user_id` varchar(22) NOT NULL,
  `is_default` tinyint NOT NULL,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_addresses_user_id` (`user_id`),
  CONSTRAINT `fk_addresses_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS `users`.`address_revisions` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `address_id` varchar(22) NOT NULL,
  `lastname` varchar(16) NOT NULL,
  `firstname` varchar(16) NOT NULL,
  `postal_code` varchar(16) NOT NULL,
  `prefecture` int NOT NULL,
  `city` varchar(32) NOT NULL,
  `address_line1` varchar(64) NOT NULL,
  `address_line2` varchar(64) NOT NULL,
  `phone_number` varchar(18) NOT NULL,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  `lastname_kana` varchar(32) NOT NULL,
  `firstname_kana` varchar(32) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_address_revisions_address_id` (`address_id`),
  CONSTRAINT `fk_address_revisions_address_id` FOREIGN KEY (`address_id`) REFERENCES `addresses` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
);
