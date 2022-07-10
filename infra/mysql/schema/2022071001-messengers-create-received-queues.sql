-- 通知関連のキュー管理テーブル
CREATE TABLE `messengers`.`received_queues` (
  `id`           VARCHAR(22) NOT NULL,          -- 通知キューID
  `event_type`   INT         NOT NULL,          -- 通知種別
  `user_type`    INT         NOT NULL,          -- 通知先ユーザー種別
  `user_ids`     JSON        NULL DEFAULT NULL, -- 通知先ユーザーID一覧
  `done`         TINYINT     NOT NULL,          -- 完了フラグ
  `created_at`   DATETIME    NOT NULL,          -- 登録日時
  `updated_at`   DATETIME    NOT NULL,          -- 更新日時
  PRIMARY KEY(`id`)
) ENGINE = InnoDB;
