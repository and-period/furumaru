-- 通知スケジュール管理テーブル
CREATE TABLE `messengers`.`schedules` (
  `message_type` INT         NOT NULL, -- メッセージ種別
  `message_id`   VARCHAR(22) NOT NULL, -- メッセージID
  `status`       INT         NOT NULL, -- 実行ステータス
  `count`        BIGINT      NOT NULL, -- 実行回数
  `sent_at`      DATETIME    NOT NULL, -- 送信日時
  `created_at`   DATETIME    NOT NULL, -- 登録日時
  `updated_at`   DATETIME    NOT NULL, -- 更新日時
  PRIMARY KEY(`message_type`, `message_id`)
) ENGINE = InnoDB;
