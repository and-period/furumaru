CREATE SCHEMA IF NOT EXISTS `media` DEFAULT CHARACTER SET utf8mb4;

CREATE TABLE IF NOT EXISTS `media`.`broadcasts` (
  `id` varchar(22) NOT NULL,
  `schedule_id` varchar(22) DEFAULT NULL,
  `type` int NOT NULL,
  `status` int NOT NULL,
  `input_url` text NOT NULL,
  `output_url` text NOT NULL,
  `archive_url` text NOT NULL,
  `cloud_front_distribution_arn` text,
  `media_live_channel_arn` text,
  `media_live_channel_id` varchar(256) DEFAULT NULL,
  `media_live_rtmp_input_arn` text,
  `media_live_rtmp_input_name` varchar(256) DEFAULT NULL,
  `media_live_mp4_input_arn` text,
  `media_live_mp4_input_name` varchar(256) DEFAULT NULL,
  `media_store_container_arn` text,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  `coordinator_id` varchar(22) NOT NULL,
  `archive_fixed` tinyint NOT NULL,
  `youtube_stream_url` text,
  `youtube_stream_key` varchar(256) DEFAULT NULL,
  `youtube_backup_url` text,
  `youtube_account` varchar(256) DEFAULT NULL,
  `youtube_broadcast_id` varchar(256) DEFAULT NULL,
  `youtube_stream_id` varchar(256) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `ui_broadcast_schedule_id` (`schedule_id`)
);

CREATE TABLE IF NOT EXISTS `media`.`broadcast_comments` (
  `id` varchar(22) NOT NULL,
  `broadcast_id` varchar(22) NOT NULL,
  `user_id` varchar(22) NOT NULL,
  `content` varchar(256) NOT NULL,
  `disabled` tinyint NOT NULL,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_broadcast_comments_broadcast_id` (`broadcast_id`),
  CONSTRAINT `fk_broadcast_comments_broadcast_id` FOREIGN KEY (`broadcast_id`) REFERENCES `broadcasts` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS `media`.`broadcast_viewer_logs` (
  `broadcast_id` varchar(22) NOT NULL,
  `session_id` varchar(22) NOT NULL,
  `created_at` datetime(3) NOT NULL,
  `user_id` varchar(22) DEFAULT NULL,
  `user_agent` text NOT NULL,
  `client_ip` varchar(15) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  PRIMARY KEY (`broadcast_id`,`session_id`,`created_at` DESC),
  CONSTRAINT `fk_broadcast_viewer_logs_broadcast_id` FOREIGN KEY (`broadcast_id`) REFERENCES `broadcasts` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS `media`.`videos` (
  `id` varchar(22) NOT NULL,
  `coordinator_id` varchar(22) NOT NULL,
  `title` varchar(255) NOT NULL,
  `description` text NOT NULL,
  `thumbnail_url` text NOT NULL,
  `video_url` text NOT NULL,
  `public` tinyint NOT NULL,
  `limited` tinyint NOT NULL,
  `published_at` datetime(3) NOT NULL,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  `display_product` tinyint NOT NULL,
  `display_experience` tinyint NOT NULL,
  PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `media`.`video_products` (
  `video_id` varchar(22) NOT NULL,
  `product_id` varchar(22) NOT NULL,
  `priority` bigint NOT NULL,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  PRIMARY KEY (`video_id`,`product_id`),
  UNIQUE KEY `ui_video_products_video_id_priority` (`video_id`,`priority`),
  CONSTRAINT `fk_video_products_video_id` FOREIGN KEY (`video_id`) REFERENCES `videos` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS `media`.`video_experiences` (
  `video_id` varchar(22) NOT NULL,
  `experience_id` varchar(22) NOT NULL,
  `priority` bigint NOT NULL,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  PRIMARY KEY (`video_id`,`experience_id`),
  UNIQUE KEY `ui_video_experiences_video_id_priority` (`video_id`,`priority`),
  CONSTRAINT `fk_video_experiences_video_id` FOREIGN KEY (`video_id`) REFERENCES `videos` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS `media`.`video_viewer_logs` (
  `video_id` varchar(22) NOT NULL,
  `session_id` varchar(22) NOT NULL,
  `created_at` datetime(3) NOT NULL,
  `user_id` varchar(22) DEFAULT NULL,
  `user_agent` text NOT NULL,
  `client_ip` varchar(15) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  PRIMARY KEY (`video_id`,`session_id`,`created_at` DESC),
  CONSTRAINT `fk_video_viewer_logs_video_id` FOREIGN KEY (`video_id`) REFERENCES `videos` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
);
