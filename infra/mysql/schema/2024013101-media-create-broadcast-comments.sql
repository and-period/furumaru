CREATE TABLE IF NOT EXISTS `media`.`broadcast_comments` (
  `id`           VARCHAR(22)  NOT NULL, -- コメントID
  `broadcast_id` VARCHAR(22)  NOT NULL, -- ライブ配信ID
  `user_id`      VARCHAR(22)  NOT NULL, -- ユーザーID
  `content`      VARCHAR(256) NOT NULL, -- コメント内容
  `disabled`     TINYINT      NOT NULL, -- コメント無効フラグ
  `created_at`   DATETIME(3)  NOT NULL, -- 登録日時
  `updated_at`   DATETIME(3)  NOT NULL, -- 更新日時
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_broadcast_comments_broadcast_id`
    FOREIGN KEY (`broadcast_id`) REFERENCES `media`.`broadcasts` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
);
