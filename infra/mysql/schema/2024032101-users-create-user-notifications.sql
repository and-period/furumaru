CREATE TABLE IF NOT EXISTS `users`.`user_notifications` (
  `user_id`        VARCHAR(22) NOT NULL, -- ユーザーID
  `email_disabled` TINYINT     NOT NULL, -- メール通知の停止
  `created_at`     DATETIME(3) NOT NULL, -- 登録日時
  `updated_at`     DATETIME(3) NOT NULL, -- 更新日時
  PRIMARY KEY (`user_id`),
  CONSTRAINT `fk_user_notifications_user_id`
    FOREIGN KEY (`user_id`) REFERENCES `users`.`users` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
);
