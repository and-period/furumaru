CREATE TABLE IF NOT EXISTS `media`.`broadcast_viewer_logs` (
  `broadcast_id` VARCHAR(22) NOT NULL,          -- ライブ配信ID
  `session_id`   VARCHAR(22) NOT NULL,          -- セッションID
  `created_at`   DATETIME(3) NOT NULL,          -- 登録日時
  `user_id`      VARCHAR(22) NULL DEFAULT NULL, -- ユーザーID
  `user_agent`   TEXT        NOT NULL,          -- ユーザーエージェント
  `client_ip`    VARCHAR(15) NOT NULL,          -- 接続元IPアドレス
  `updated_at`   DATETIME(3) NOT NULL,          -- 更新日時
  PRIMARY KEY (`broadcast_id`, `session_id`, `created_at` DESC),
  CONSTRAINT `fk_broadcast_viewer_logs_broadcast_id`
    FOREIGN KEY (`broadcast_id`) REFERENCES `media`.`broadcasts` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
);
