CREATE TABLE IF NOT EXISTS `ai_chat_sessions` (
  `id`         VARCHAR(22)  NOT NULL,
  `admin_id`   VARCHAR(22)  NOT NULL,
  `product_id` VARCHAR(22)  NOT NULL DEFAULT '',
  `title`      VARCHAR(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `created_at` DATETIME(3)  NOT NULL,
  `updated_at` DATETIME(3)  NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_admin_id` (`admin_id`)
);

CREATE TABLE IF NOT EXISTS `ai_chat_messages` (
  `id`         VARCHAR(22) NOT NULL,
  `session_id` VARCHAR(22) NOT NULL,
  `role`       VARCHAR(16) NOT NULL,
  `content`    LONGTEXT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `created_at` DATETIME(3) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_session_id_created` (`session_id`, `created_at`)
);
