CREATE SCHEMA IF NOT EXISTS `stores` DEFAULT CHARACTER SET utf8mb4;

-- 商品種別テーブル
CREATE TABLE `stores`.`categories` (
  `id`         VARCHAR(22) NOT NULL, -- カテゴリID
  `name`       VARCHAR(32) NOT NULL, -- カテゴリ名
  `created_at` DATETIME    NOT NULL, -- 登録日時
  `updated_at` DATETIME    NOT NULL, -- 更新日時
  PRIMARY KEY(`id`)
) ENGINE = InnoDB;

CREATE UNIQUE INDEX `ui_categories_name` ON `stores`.`categories` (`name` ASC) VISIBLE;

-- 品目テーブル
CREATE TABLE `stores`.`product_types` (
  `id`          VARCHAR(22) NOT NULL, -- 品目ID
  `category_id` VARCHAR(22) NOT NULL, -- カテゴリID
  `name`        VARCHAR(32) NOT NULL, -- 品目名
  `created_at`  DATETIME    NOT NULL, -- 登録日時
  `updated_at`  DATETIME    NOT NULL, -- 更新日時
  PRIMARY KEY(`id`),
  CONSTRAINT `fk_product_types_category_id`
    FOREIGN KEY (`category_id`) REFERENCES `stores`.`categories` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB;

CREATE UNIQUE INDEX `ui_product_types_category_id_name` ON `stores`.`product_types` (`category_id` ASC, `name` ASC) VISIBLE;

-- 商品テーブル
CREATE TABLE `stores`.`products` (
  `id`                VARCHAR(22)  NOT NULL,          -- 商品ID
  `producer_id`       VARCHAR(22)  NOT NULL,          -- 生産者ID
  `category_id`       VARCHAR(22)  NULL DEFAULT NULL, -- 商品種別ID
  `product_type_id`   VARCHAR(22)  NULL DEFAULT NULL, -- 品目ID
  `name`              VARCHAR(128) NOT NULL,          -- 商品名
  `description`       TEXT         NOT NULL,          -- 商品説明
  `public`            TINYINT      NOT NULL,          -- 公開フラグ(0:下書き, 1:公開)
  `inventory`         BIGINT       NOT NULL,          -- 在庫数
  `weight`            BIGINT       NOT NULL,          -- 重量(数値)
  `weight_unit`       INT          NOT NULL,          -- 重量単位
  `item`              BIGINT       NOT NULL,          -- 数量(数値)
  `item_unit`         VARCHAR(16)  NOT NULL,          -- 数量単位
  `item_description`  VARCHAR(64)  NOT NULL,          -- 数量単位(説明)
  `media`             JSON         NULL DEFAULT NULL, -- メディア一覧(URL)
  `price`             BIGINT       NOT NULL,          -- 販売価格
  `box60_rate`        BIGINT       NOT NULL,          -- 箱の占有率(サイズ:60)
  `box80_rate`        BIGINT       NOT NULL,          -- 箱の占有率(サイズ:80)
  `box100_rate`       BIGINT       NOT NULL,          -- 箱の占有率(サイズ:100)
  `origin_prefecture` VARCHAR(32)  NOT NULL,          -- 原産地(都道府県)
  `origin_city`       VARCHAR(32)  NOT NULL,          -- 原産地(市区町村)
  `created_at`        DATETIME     NOT NULL,          -- 登録日時
  `updated_at`        DATETIME     NOT NULL,          -- 更新日時
  `deleted_at`        DATETIME     NULL DEFAULT NULL, -- 削除日時
  PRIMARY KEY(`id`),
  CONSTRAINT `fk_products_caterogy_id`
    FOREIGN KEY (`category_id`) REFERENCES `stores`.`categories` (`id`)
    ON DELETE SET NULL ON UPDATE CASCADE,
  CONSTRAINT `fk_products_product_type_id`
    FOREIGN KEY (`product_type_id`) REFERENCES `stores`.`product_types` (`id`)
    ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE = InnoDB;
