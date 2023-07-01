-- 商品タグテーブル
CREATE TABLE `stores`.`product_tags` (
  `id`          VARCHAR(22) NOT NULL, -- 商品タグID
  `name`        VARCHAR(32) NOT NULL, -- 商品タグ名
  `created_at`  DATETIME    NOT NULL, -- 登録日時
  `updated_at`  DATETIME    NOT NULL, -- 更新日時
  PRIMARY KEY(`id`)
) ENGINE = InnoDB;

CREATE UNIQUE INDEX `ui_product_tags_name` ON `stores`.`product_tags` (`name` ASC) VISIBLE;
