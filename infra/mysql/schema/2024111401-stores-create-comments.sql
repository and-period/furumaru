-- ユーザーが投稿したスポット情報
CREATE TABLE IF NOT EXISTS `stores`.`spots` (
  `id`                VARCHAR(22)    NOT NULL, -- ID
  `user_id`           VARCHAR(22)    NOT NULL, -- ユーザID
  `name`              VARCHAR(64)    NOT NULL, -- スポット名
  `description`       TEXT           NOT NULL, -- 説明
  `thumbnail_url`     TEXT           NOT NULL, -- サムネイルURL
  `latitude`          decimal(10, 6) NOT NULL, -- 緯度
  `longitude`         decimal(10, 6) NOT NULL, -- 経度
  `approved`          TINYINT        NOT NULL, -- 承認フラグ
  `approved_admin_id` VARCHAR(22)    NOT NULL, -- 承認した管理者ID
  `created_at`        DATETIME(3)    NOT NULL, -- 登録日時
  `updated_at`        DATETIME(3)    NOT NULL, -- 更新日時
  PRIMARY KEY (`id`)
);

CREATE FULLTEXT INDEX `ftx_spots_name_description` ON `stores`.`spots` (`name`, `description`) WITH PARSER ngram;
