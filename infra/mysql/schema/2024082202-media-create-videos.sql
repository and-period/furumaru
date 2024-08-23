CREATE TABLE IF NOT EXISTS `media`.`videos` (
  `id`             VARCHAR(22)  NOT NULL,          -- 動画ID
  `coordinator_id` VARCHAR(22)  NOT NULL,          -- コーディネータID
  `title`          VARCHAR(255) NOT NULL,          -- タイトル
  `description`    TEXT         NOT NULL,          -- 説明
  `thumbnail_url`  TEXT         NOT NULL,          -- サムネイルURL
  `video_url`      TEXT         NOT NULL,          -- オリジナル動画URL
  `public`         TINYINT      NOT NULL,          -- 公開設定
  `limited`        TINYINT      NOT NULL,          -- 限定公開設定
  `published_at`   DATETIME(3)  NOT NULL,          -- 公開日時
  `created_at`     DATETIME(3)  NOT NULL,          -- 登録日時
  `updated_at`     DATETIME(3)  NOT NULL,          -- 更新日時
  PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `media`.`video_products` (
  `video_id`   VARCHAR(22)  NOT NULL, -- 動画ID
  `product_id` VARCHAR(22)  NOT NULL, -- 商品ID
  `priority`   BIGINT       NOT NULL, -- 優先度
  `created_at` DATETIME(3)  NOT NULL, -- 登録日時
  `updated_at` DATETIME(3)  NOT NULL, -- 更新日時
  PRIMARY KEY (`video_id`, `product_id`),
  CONSTRAINT `fk_video_products_video_id`
    FOREIGN KEY (`video_id`) REFERENCES `media`.`videos` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE UNIQUE INDEX `ui_video_products_video_id_priority` ON `media`.`video_products` (`video_id`, `priority`);

CREATE TABLE IF NOT EXISTS `media`.`video_experiences` (
  `video_id`      VARCHAR(22)  NOT NULL, -- 動画ID
  `experience_id` VARCHAR(22)  NOT NULL, -- 体験ID
  `priority`      BIGINT       NOT NULL, -- 優先度
  `created_at`    DATETIME(3)  NOT NULL, -- 登録日時
  `updated_at`    DATETIME(3)  NOT NULL, -- 更新日時
  PRIMARY KEY (`video_id`, `experience_id`),
  CONSTRAINT `fk_video_experiences_video_id`
    FOREIGN KEY (`video_id`) REFERENCES `media`.`videos` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE UNIQUE INDEX `ui_video_experiences_video_id_priority` ON `media`.`video_experiences` (`video_id`, `priority`);
