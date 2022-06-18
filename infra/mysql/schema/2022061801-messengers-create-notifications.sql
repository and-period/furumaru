CREATE SCHEMA IF NOT EXISTS `messengerss` DEFAULT CHARACTER SET utf8mb4;

CREATE TABLE `messengers`.`notifications` (
  `id`                VARCHAR(22)  NOT NULL,          -- お知らせID (Primary Key用)
  `created_by`        VARCHAR(22)  NOT NULL,          -- 登録者ID (Foreign Key)
  `creater_name`      VARCHAR(32)  NOT NULL,          -- 登録者名
  `updated_by`        VARCHAR(22)  NOT NULL,          -- 更新者ID (Foreign Key)
  `title`             VARCHAR(128) NOT NULL,          -- タイトル
  `body`              TEXT         NOT NULL,          -- 本文
  `post_period_start` DATETIME     NOT NULL,          -- 投稿開始時間
  `post_period_end`   DATETIME     NOT NULL,          -- 投稿終了時間
  `post_targets`      JSON         NULL DEFAULT NULL, -- 投稿範囲(PostTarget[])
  `images`            JSON         NULL DEFAULT NULL, -- 参考画像(S3URL)
  `created_at`        DATETIME     NOT NULL,          -- 登録日時
  `updated_at`        DATETIME     NOT NULL,          -- 更新日時
  PRIMARY KEY(`id`),
  CONSTRAINT `fk_notifications_asmin_id`
    FOREIGN KEY (`created_id`) REFERENCES `users`.`admins` (`id`)
    ON DELETE SET CASCADE ON UPDATE CASCADE
    FOREIGN KEY (`updated_id`) REFERENCES `users`.`admins` (`id`)
    ON DELETE SET CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB;
