CREATE SCHEMA IF NOT EXISTS `shops` DEFAULT CHARACTER SET utf8mb4;

-- 店舗情報テーブル
CREATE TABLE `shops`.`stores` (
  `id`            BIGINT       NOT NULL AUTO_INCREMENT, -- 商品ID
  `name`          VARCHAR(64)  NOT NULL,                -- 店舗名
  `thumbnail_url` TEXT         NOT NULL,                -- サムネイルURL
  `created_at`    DATETIME     NOT NULL,                -- 登録日時
  `updated_at`    DATETIME     NOT NULL,                -- 更新日時
  `deleted_at`    DATETIME     NULL DEFAULT NULL,       -- 削除日時
  PRIMARY KEY (`id`)
) ENGINE = InnoDB;

CREATE UNIQUE INDEX `ui_stores_name` ON `shops`.`stores` (`name` ASC) VISIBLE;

-- 店舗スタッフ情報テーブル
CREATE TABLE `shops`.`staffs` (
  `store_id`     BIGINT       NOT NULL, -- 店舗ID
  `user_id`      VARCHAR(22)  NOT NULL, -- 販売者ID
  `role`         INT          NOT NULL, -- 権限
  `created_at`   DATETIME     NOT NULL, -- 登録日時
  `updated_at`   DATETIME     NOT NULL, -- 更新日時
  PRIMARY KEY (`store_id`, `user_id`),
  CONSTRAINT `fk_staffs_store_id`
    FOREIGN KEY (`store_id`) REFERENCES `shops`.`stores` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
);
