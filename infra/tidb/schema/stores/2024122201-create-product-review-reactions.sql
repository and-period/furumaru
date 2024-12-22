CREATE TABLE IF NOT EXISTS `stores`.`product_review_reactions` (
  `review_id`     VARCHAR(22) NOT NULL, -- 商品レビューID
  `user_id`       VARCHAR(22) NOT NULL, -- ユーザーID
  `reaction_type` INT         NOT NULL, -- リアクション種別
  `created_at`    DATETIME(3) NOT NULL, -- 登録日時
  `updated_at`    DATETIME(3) NOT NULL, -- 更新日時
  PRIMARY KEY (`review_id`, `user_id`),
  CONSTRAINT `fk_product_review_reactions_review_id`
    FOREIGN KEY (`review_id`) REFERENCES `stores`.`product_reviews` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
);
