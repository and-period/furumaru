CREATE TABLE IF NOT EXISTS `media`.`video_comments` (
  `id`         VARCHAR(22)  NOT NULL, -- コメントID
  `video_id`   VARCHAR(22)  NOT NULL, -- オンデマンド配信ID
  `user_id`    VARCHAR(22)  NOT NULL, -- ユーザーID
  `content`    VARCHAR(256) NOT NULL, -- コメント内容
  `disabled`   TINYINT      NOT NULL, -- コメント無効フラグ
  `created_at` DATETIME(3)  NOT NULL, -- 登録日時
  `updated_at` DATETIME(3)  NOT NULL, -- 更新日時
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_video_comments_video_id`
    FOREIGN KEY (`video_id`) REFERENCES `media`.`videos` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
);
