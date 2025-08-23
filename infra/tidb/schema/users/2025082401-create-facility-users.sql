CREATE TABLE IF NOT EXISTS `users`.`facility_users` (
  `user_id`        VARCHAR(22)  NOT NULL,
  `external_id`    VARCHAR(36)  NOT NULL,
  `producer_id`    VARCHAR(22)  NOT NULL,
  `lastname`       VARCHAR(16)  NOT NULL,
  `firstname`      VARCHAR(16)  NOT NULL,
  `lastname_kana`  VARCHAR(32)  NOT NULL,
  `firstname_kana` VARCHAR(32)  NOT NULL,
  `provider_type`  INT          NOT NULL,
  `email`          VARCHAR(256) NOT NULL,
  `phone_number`   VARCHAR(18)  NULL DEFAULT NULL,
  `exists`         TINYINT      DEFAULT '1',
  `created_at`     DATETIME(3)  NOT NULL,
  `updated_at`     DATETIME(3)  NOT NULL,
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `ui_guests_email_producer_id` (`exists` DESC,`email`,`producer_id`),
  CONSTRAINT `fk_guests_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_guests_producer_id` FOREIGN KEY (`producer_id`) REFERENCES `producers` (`admin_id`) ON DELETE CASCADE ON UPDATE CASCADE
);
