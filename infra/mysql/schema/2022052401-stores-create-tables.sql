CREATE SCHEMA IF NOT EXISTS `stores` DEFAULT CHARACTER SET utf8mb4;

-- 商品情報テーブル
CREATE TABLE `stores`.`products` (
  `id`           VARCHAR(22)  NOT NULL,          -- 商品ID
  `name`         VARCHAR(128) NOT NULL,          -- 商品名
  `description`  TEXT         NOT NULL,          -- 説明
  `public`       TINYINT      NOT NULL,          -- 公開フラグ (0:下書き, 1:公開)
  `type`         VARCHAR(64)  NOT NULL,          -- 商品タイプ (検索用)
  `media`        JSON,                           -- メディア一覧(URL)
  `price`        BIGINT       NOT NULL,          -- 販売価格
  `discount`     BIGINT       NOT NULL,          -- 割引価格
  `tax_included` TINYINT      NOT NULL,          -- 税込み
  `unit_cost`    BIGINT       NOT NULL,          -- 仕入れ価格
  `sku`          VARCHAR(40),                    -- 在庫SKU
  `barcode`      VARCHAR(16),                    -- 在庫バーコード (ISBN、UPC、GTINなど)
  `created_at`   DATETIME     NOT NULL,          -- 登録日時
  `updated_at`   DATETIME     NOT NULL,          -- 更新日時
  `deleted_at`   DATETIME     NULL DEFAULT NULL, -- 削除日時
  PRIMARY KEY(`id`)
) ENGINE = InnoDB;

-- 商品タグ情報テーブル
CREATE TABLE `stores`.`product_tags` (
  `product_id` VARCHAR(22) NOT NULL, -- 商品ID
  `name`       VARCHAR(32) NOT NULL, -- タグ名
  `created_at` DATETIME    NOT NULL, -- 登録日時
  `updated_at` DATETIME    NOT NULL, -- 更新日時
  PRIMARY KEY(`product_id`, `name`),
  CONSTRAINT `fk_product_tags_product_id`
    FOREIGN KEY (`product_id`) REFERENCES `stores`.`products` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB;

CREATE UNIQUE INDEX `ui_product_tags_product_id_name` ON `stores`.`product_tags` (`product_id` ASC, `name` ASC) VISIBLE;
CREATE INDEX `idx_product_tags_product_id` ON `stores`.`product_tags` (`product_id` ASC) VISIBLE;
CREATE INDEX `idx_product_tags_name` ON `stores`.`product_tags` (`name` ASC) VISIBLE;
