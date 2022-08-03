-- メッセージテーブル
CREATE TABLE `messengers`.`messages` (
  `id`          VARCHAR(22)  NOT NULL, -- メッセージID
  `user_type`   INT          NOT NULL, -- ユーザー種別
  `user_id`     VARCHAR(22)  NOT NULL, -- ユーザーID
  `type`        INT          NOT NULL, -- メッセージ種別
  `title`       VARCHAR(256) NOT NULL, -- メッセージ件名
  `body`        TEXT         NOT NULL, -- メッセージ内容
  `link`        TEXT         NOT NULL, -- 遷移先リンク
  `read`        TINYINT      NOT NULL, -- 既読フラグ
  `received_at` DATETIME     NOT NULL, -- 受信日時
  `created_at`  DATETIME     NOT NULL, -- 登録日時
  `updated_at`  DATETIME     NOT NULL, -- 更新日時
  PRIMARY KEY(`id`)
) ENGINE = InnoDB;

CREATE INDEX `idx_messages_user_type_user_id` ON `messengers`.`messages` (`user_type`, `user_id`) VISIBLE;
