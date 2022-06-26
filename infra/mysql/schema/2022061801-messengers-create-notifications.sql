CREATE SCHEMA IF NOT EXISTS `messengers` DEFAULT CHARACTER SET utf8mb4;

DROP TABLE IF EXISTS `messengers`.`notifications`;

CREATE TABLE `messengers`.`notifications` (
  `id`           VARCHAR(22)  NOT NULL,          -- お知らせID (Primary Key用)
  `created_by`   VARCHAR(22)  NOT NULL,          -- 登録者ID
  `creator_name` VARCHAR(32)  NOT NULL,          -- 登録者名
  `updated_by`   VARCHAR(22)  NOT NULL,          -- 更新者ID
  `title`        VARCHAR(128) NOT NULL,          -- タイトル
  `body`         TEXT         NOT NULL,          -- 本文
  `published_at` DATETIME     NOT NULL,          -- 掲載開始時間
  `targets`      JSON         NULL DEFAULT NULL, -- 掲載範囲(PostTarget[])
  `public`       TINYINT      NOT NULL,          -- 公開フラグ
  `created_at`   DATETIME     NOT NULL,          -- 登録日時
  `updated_at`   DATETIME     NOT NULL,          -- 更新日時
  PRIMARY KEY(`id`)
) ENGINE = InnoDB;
