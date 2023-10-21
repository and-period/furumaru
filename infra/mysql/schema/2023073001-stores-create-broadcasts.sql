-- ライブ配信情報
CREATE TABLE IF NOT EXISTS `stores`.`broadcasts` (
  `id`          VARCHAR(22) NOT NULL, -- ライブ配信ID
  `schedule_id` VARCHAR(22) NOT NULL, -- 開催スケジュールID
  `status`      INT         NOT NULL, -- ライブ配信状況
  `input_url`   TEXT        NOT NULL, -- ライブ配信URL(入力)
  `output_url`  TEXT        NOT NULL, -- ライブ配信URL(出力)
  `created_at`  DATETIME(3) NOT NULL, -- 登録日時
  `updated_at`  DATETIME(3) NOT NULL, -- 更新日時
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_broadcasts_schedule_id`
    FOREIGN KEY (`schedule_id`) REFERENCES `stores`.`schedules` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB;
