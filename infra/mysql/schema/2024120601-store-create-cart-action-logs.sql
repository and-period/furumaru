CREATE TABLE IF NOT EXISTS `stores`.`cart_action_logs` (
  `session_id`   VARCHAR(22) NOT NULL,          -- セッションID
  `created_at`   DATETIME(3) NOT NULL,          -- 登録日時
  `type`         INT         NOT NULL,          -- カート操作種別
  `user_id`      VARCHAR(22) NULL DEFAULT NULL, -- ユーザーID
  `user_agent`   TEXT        NOT NULL,          -- ユーザーエージェント
  `client_ip`    VARCHAR(15) NOT NULL,          -- 接続元IPアドレス
  `product_id`   VARCHAR(22) NULL DEFAULT NULL, -- 商品ID
  `updated_at`   DATETIME(3) NOT NULL,          -- 更新日時
  PRIMARY KEY (`session_id`, `created_at` DESC, `type`),
  CONSTRAINT `fk_cart_action_logs_product_id`
    FOREIGN KEY (`product_id`) REFERENCES `stores`.`products` (`id`)
    ON DELETE SET NULL ON UPDATE CASCADE
);
