CREATE TABLE IF NOT EXISTS `stores`.`product_reviews` (
  `id`         VARCHAR(22) NOT NULL,          -- 商品レビューID
  `product_id` VARCHAR(22) NOT NULL,          -- 商品ID
  `user_id`    VARCHAR(22) NOT NULL,          -- ユーザーID
  `rate`       BIGINT      NOT NULL,          -- 評価
  `title`      VARCHAR(64) NOT NULL,          -- タイトル
  `comment`    TEXT        NOT NULL,          -- コメント
  `created_at` DATETIME(3) NOT NULL,          -- 登録日時
  `updated_at` DATETIME(3) NOT NULL,          -- 更新日時
  `deleted_at` DATETIME(3) NULL DEFAULT NULL, -- 削除日時
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_product_reviews_product_id`
    FOREIGN KEY (`product_id`) REFERENCES `stores`.`products` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX `idx_product_reviews_product_id_rate` ON `stores`.`product_reviews` (`product_id`, `rate`) VISIBLE;
