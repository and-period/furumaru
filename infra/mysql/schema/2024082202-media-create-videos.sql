CREATE TABLE IF NOT EXISTS `media`.`videos` (
  `id`                  VARCHAR(22)  NOT NULL,          -- 動画ID
  `coordinator_id`      VARCHAR(22)  NOT NULL,          -- コーディネータID
  `product_ids`         JSON         NULL DEFAULT NULL, -- 商品ID一覧
  `experience_ids`      JSON         NULL DEFAULT NULL, -- 体験ID一覧
  `title`               VARCHAR(255) NOT NULL,          -- タイトル
  `description`         TEXT         NOT NULL,          -- 説明
  `thumbnail_url`       TEXT         NOT NULL,          -- サムネイルURL
  `video_url`           TEXT         NOT NULL,          -- オリジナル動画URL
  `public`              TINYINT      NOT NULL,          -- 公開設定
  `limited`             TINYINT      NOT NULL,          -- 限定公開設定
  `published_at`        DATETIME(3)  NOT NULL,          -- 公開日時
  `created_at`          DATETIME(3)  NOT NULL,          -- 登録日時
  `updated_at`          DATETIME(3)  NOT NULL,          -- 更新日時
  PRIMARY KEY (`id`)
);
