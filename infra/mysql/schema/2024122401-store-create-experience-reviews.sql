CREATE TABLE IF NOT EXISTS `stores`.`experience_reviews` (
  `id`            VARCHAR(22) NOT NULL,          -- 体験レビューID
  `experience_id` VARCHAR(22) NOT NULL,          -- 体験ID
  `user_id`       VARCHAR(22) NOT NULL,          -- ユーザーID
  `rate`          BIGINT      NOT NULL,          -- 評価
  `title`         VARCHAR(64) NOT NULL,          -- タイトル
  `comment`       TEXT        NOT NULL,          -- コメント
  `created_at`    DATETIME(3) NOT NULL,          -- 登録日時
  `updated_at`    DATETIME(3) NOT NULL,          -- 更新日時
  `deleted_at`    DATETIME(3) NULL DEFAULT NULL, -- 削除日時
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_experience_reviews_experience_id`
    FOREIGN KEY (`experience_id`) REFERENCES `stores`.`experiences` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX `idx_experience_reviews_experience_id_rate` ON `stores`.`experience_reviews` (`experience_id`, `rate`) VISIBLE;

CREATE TABLE IF NOT EXISTS `stores`.`experience_review_reactions` (
  `review_id`     VARCHAR(22) NOT NULL, -- 体験レビューID
  `user_id`       VARCHAR(22) NOT NULL, -- ユーザーID
  `reaction_type` INT         NOT NULL, -- リアクション種別
  `created_at`    DATETIME(3) NOT NULL, -- 登録日時
  `updated_at`    DATETIME(3) NOT NULL, -- 更新日時
  PRIMARY KEY (`review_id`, `user_id`),
  CONSTRAINT `fk_experience_review_reactions_review_id`
    FOREIGN KEY (`review_id`) REFERENCES `stores`.`experience_reviews` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
);
