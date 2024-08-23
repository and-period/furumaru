CREATE TABLE IF NOT EXISTS `stores`.`experience_types` (
  `id`         VARCHAR(22)  NOT NULL,          -- 体験種別ID
  `name`       VARCHAR(32)  NOT NULL,           -- 体験種別名
  `created_at` DATETIME(3)  NOT NULL,          -- 登録日時
  `updated_at` DATETIME(3)  NOT NULL,          -- 更新日時
  PRIMARY KEY (`id`)
);

CREATE UNIQUE INDEX `ui_experience_types_name` ON `stores`.`experience_types` (`name`) VISIBLE;
CREATE FULLTEXT INDEX `ftx_experience_types` ON `stores`.`experience_types` (`name`) WITH PARSER ngram;

CREATE TABLE IF NOT EXISTS `stores`.`experiences` (
  `id`                  VARCHAR(22)  NOT NULL,          -- 体験ID
  `coordinator_id`      VARCHAR(22)  NOT NULL,          -- コーディネータID
  `producer_id`         VARCHAR(22)  NOT NULL,          -- 生産者ID
  `experience_type_id`  VARCHAR(22)  NOT NULL,          -- 体験種別ID
  `title`               VARCHAR(128) NOT NULL,          -- タイトル
  `description`         TEXT         NOT NULL,          -- 説明
  `public`              TINYINT      NOT NULL,          -- 公開設定
  `sold_out`            TINYINT      NOT NULL,          -- 完売設定
  `media`               JSON         NULL DEFAULT NULL, -- メディア一覧
  `recommended_points`  JSON         NULL DEFAULT NULL, -- おすすめポイント一覧
  `promotion_video_url` TEXT         NOT NULL,          -- 紹介動画URL
  `host_prefecture`     BIGINT       NOT NULL,          -- 開催場所（都道府県）
  `host_city`           VARCHAR(32)  NOT NULL,          -- 開催場所（市区町村）
  `start_at`            DATETIME(3)  NOT NULL,          -- 募集開始日時
  `end_at`              DATETIME(3)  NOT NULL,          -- 募集終了日時
  `created_at`          DATETIME(3)  NOT NULL,          -- 登録日時
  `updated_at`          DATETIME(3)  NOT NULL,          -- 更新日時
  `deleted_at`          DATETIME(3)  NULL DEFAULT NULL, -- 削除日時
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_experiences_experience_type_id`
    FOREIGN KEY (`experience_type_id`) REFERENCES `stores`.`experience_type_id` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE,
);

CREATE FULLTEXT INDEX `ftx_experiences_title_description` ON `stores`.`experiences` (`title`, `description`) WITH PARSER ngram;

CREATE TABLE IF NOT EXISTS `stores`.`experience_revisions` (
  `id`                       BIGINT      NOT NULL AUTO_INCREMENT, -- 体験変更履歴ID
  `experience_id`            VARCHAR(22) NOT NULL,                -- 体験ID
  `price_adult`              BIGINT      NOT NULL,                -- 大人料金
  `price_junior_high_school` BIGINT      NOT NULL,                -- 中学生料金
  `price_elementary_school`  BIGINT      NOT NULL,                -- 小学生料金
  `price_preschool`          BIGINT      NOT NULL,                -- 幼児料金
  `created_at`               DATETIME(3) NOT NULL,                -- 登録日時
  `updated_at`               DATETIME(3) NOT NULL,                -- 更新日時
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_experience_revisions_experience_id`
    FOREIGN KEY (`experience_id`) REFERENCES `stores`.`experiences` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
);
