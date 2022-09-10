-- ライブ配信管理テーブル
CREATE TABLE `stores`.`lives` (
  `id`            VARCHAR(22) NOT NULL,          -- テンプレートID
  `schedule_id`   VARCHAR(22) NOT NULL,          -- 開催スケジュールID
  `title`         VARCHAR(64) NOT NULL,          -- タイトル
  `description`   TEXT        NOT NULL,          -- 説明
  `producer_id`   VARCHAR(22) NOT NULL,          -- 生産者ID
  `start_at`      DATETIME    NOT NULL,          -- 配信開始日時
  `end_at`        DATETIME    NOT NULL,          -- 配信終了日時
  `cancel`        TINYINT     NOT NULL,          -- 配信中止フラグ
  `recommends`    JSON        NULL DEFAULT NULL, -- おすすめ商品一覧
  `created_at`    DATETIME    NOT NULL,          -- 登録日時
  `updated_at`    DATETIME    NOT NULL,          -- 更新日時
  `deleted_at`    DATETIME    NOT NULL,          -- 削除日時
  PRIMARY KEY(`id`),
  CONSTRAINT `fk_lives_schedule_id`
  FOREIGN KEY (`schedule_id`) REFERENCES `stores`.`schedules` (`id`)
  ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB;
