-- 開催スケジュール管理テーブル
CREATE TABLE `stores`.`schedules` (
  `id`            VARCHAR(22) NOT NULL, -- テンプレートID
  `title`         VARCHAR(64) NOT NULL, -- タイトル
  `description`   TEXT        NOT NULL, -- 説明
  `thumbnail_url` TEXT        NOT NULL, -- サムネイルURL
  `start_at`      DATETIME    NOT NULL, -- 開催開始日時
  `end_at`        DATETIME    NOT NULL, -- 開催終了日時
  `cancel`        TINYINT     NOT NULL, -- 開催中止フラグ
  `created_at`    DATETIME    NOT NULL, -- 登録日時
  `updated_at`    DATETIME    NOT NULL, -- 更新日時
  `deleted_at`    DATETIME    NOT NULL, -- 削除日時
  PRIMARY KEY(`id`)
) ENGINE = InnoDB;
