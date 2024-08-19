CREATE TABLE IF NOT EXISTS `media`.`video_categories` (
  `id`         VARCHAR(22)  NOT NULL, -- カテゴリID
  `name`       VARCHAR(255) NOT NULL, -- カテゴリ名
  `created_at` DATETIME(3)  NOT NULL, -- 登録日時
  `updated_at` DATETIME(3)  NOT NULL, -- 更新日時
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_video_categories_name` (`name`)
);

CREATE TABLE IF NOT EXISTS `media`.`videos` (
  `id`                  VARCHAR(22)  NOT NULL,          -- 動画ID
  `coordinator_id`      VARCHAR(22)  NOT NULL,          -- コーディネータID
  `category_ids`        JSON         NULL DEFAULT NULL, -- カテゴリID一覧
  `product_ids`         JSON         NULL DEFAULT NULL, -- 商品ID一覧
  `title`               VARCHAR(255) NOT NULL,          -- タイトル
  `description`         TEXT         NOT NULL,          -- 説明
  `thumbnail_url`       VARCHAR(255) NOT NULL,          -- サムネイルURL
  `original_video_url`  VARCHAR(255) NOT NULL,          -- オリジナル動画URL
  `processed_video_url` VARCHAR(255) NOT NULL,          -- 加工済み動画URL
  `public`              TINYINT      NOT NULL,          -- 公開設定
  `limited`             TINYINT      NOT NULL,          -- 限定公開設定
  `published_at`        DATETIME(3)  NOT NULL,          -- 公開日時
  `created_at`          DATETIME(3)  NOT NULL,          -- 登録日時
  `updated_at`          DATETIME(3)  NOT NULL,          -- 更新日時
  PRIMARY KEY (`id`)
);
