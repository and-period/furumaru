DROP TABLE IF EXISTS `stores`.`order_items`;
DROP TABLE IF EXISTS `stores`.`order_payments`;
DROP TABLE IF EXISTS `stores`.`order_fulfillments`;
DROP TABLE IF EXISTS `stores`.`order_activities`;
DROP TABLE IF EXISTS `stores`.`orders`;

-- 注文履歴管理テーブル
CREATE TABLE IF NOT EXISTS `stores`.`orders` (
  `id`                 VARCHAR(22) NOT NULL,          -- 注文履歴ID
  `user_id`            VARCHAR(22) NOT NULL,          -- ユーザーID
  `coordinator_id`     VARCHAR(22) NOT NULL,          -- 配送担当者ID
  `schedule_id`        VARCHAR(22) NULL,              -- マルシェ開催スケジュールID
  `promotion_id`       VARCHAR(22) NULL,              -- プロモーションID
  `payment_status`     INT         NOT NULL,          -- 支払い状況
  `fulfillment_status` INT         NOT NULL,          -- 配送状況
  `cancel_type`        INT         NOT NULL,          -- キャンセル種別
  `cancel_reason`      TEXT        NOT NULL,          -- キャンセル理由
  `created_at`         DATETIME    NOT NULL,          -- 登録日時
  `updated_at`         DATETIME    NOT NULL,          -- 更新日時
  `ordered_at`         DATETIME    NULL DEFAULT NULL, -- 注文日時
  `paid_at`            DATETIME    NULL DEFAULT NULL, -- 仮売上日時
  `captured_at`        DATETIME    NULL DEFAULT NULL, -- 売上確定日時
  `failed_at`          DATETIME    NULL DEFAULT NULL, -- 支払い失敗日時
  `refunded_at`        DATETIME    NULL DEFAULT NULL, -- キャンセル日時
  `shipped_at`         DATETIME    NULL DEFAULT NULL, -- 配送日時
  `deleted_at`         DATETIME    NULL DEFAULT NULL, -- 削除日時
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_orders_schedule_id`
    FOREIGN KEY (`schedule_id`) REFERENCES `stores`.`schedules` (`id`)
    ON DELETE SET NULL ON UPDATE CASCADE,
  CONSTRAINT `fk_orders_promotion_id`
    FOREIGN KEY (`promotion_id`) REFERENCES `stores`.`promotions` (`id`)
    ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE = InnoDB;

-- 注文イベントログ管理テーブル
CREATE TABLE IF NOT EXISTS `stores`.`activities` (
  `id`         VARCHAR(22) NOT NULL,          -- イベントログID
  `order_id`   VARCHAR(22) NOT NULL,          -- 注文履歴ID
  `user_id`    VARCHAR(22) NOT NULL,          -- ユーザーID
  `event_type` INT         NOT NULL,          -- イベントログ種別
  `detail`     TEXT        NOT NULL,          -- イベントログ詳細
  `metadata`   JSON        NULL DEFAULT NULL, -- メタデータ
  `created_at` DATETIME    NOT NULL,          -- 登録日時
  `updated_at` DATETIME    NOT NULL,          -- 更新日時
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_activities_order_id`
    FOREIGN KEY (`order_id`) REFERENCES `stores`.`orders` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB;

-- 注文商品管理テーブル
CREATE TABLE IF NOT EXISTS `stores`.`order_items` (
  `order_id`   VARCHAR(22) NOT NULL, -- 注文履歴ID
  `product_id` VARCHAR(22) NOT NULL, -- 商品ID
  `price`      BIGINT      NOT NULL, -- 購入時価格
  `quantity`   BIGINT      NOT NULL, -- 購入数
  `created_at` DATETIME    NOT NULL, -- 登録日時
  `updated_at` DATETIME    NOT NULL, -- 更新日時
  PRIMARY KEY (`order_id`, `product_id`),
  CONSTRAINT `fk_order_items_order_id`
    FOREIGN KEY (`order_id`) REFERENCES `stores`.`orders` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_order_items_product_id`
    FOREIGN KEY (`product_id`) REFERENCES `stores`.`products` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB;

CREATE UNIQUE INDEX `ui_order_items_order_id_product_id`
  ON `stores`.`order_items` (`order_id` ASC, `product_id` ASC) VISIBLE;

-- 注文支払い管理テーブル
CREATE TABLE IF NOT EXISTS `stores`.`payments` (
  `order_id`        VARCHAR(22) NOT NULL,          -- 注文履歴ID
  `address_id`      VARCHAR(22) NOT NULL,          -- 請求先情報ID
  `transaction_id`  VARCHAR(22) NOT NULL,          -- 取引管理番号
  `method_type`     INT         NOT NULL,          -- 支払い種別
  `method_id`       VARCHAR(22) NULL DEFAULT NULL, -- 支払いID
  `subtotal`        BIGINT      NOT NULL,          -- 購入金額
  `discount`        BIGINT      NOT NULL,          -- 割引金額
  `shipping_fee`    BIGINT      NOT NULL,          -- 配送手数料
  `tax`             BIGINT      NOT NULL,          -- 消費税
  `total`           BIGINT      NOT NULL,          -- 合計金額
  `refund_total`    BIGINT      NOT NULL,          -- 返金金額
  `created_at`      DATETIME    NOT NULL,          -- 登録日時
  `updated_at`      DATETIME    NOT NULL,          -- 更新日時
  `deleted_at`      DATETIME    NULL DEFAULT NULL, -- 削除日時
  PRIMARY KEY (`order_id`),
  CONSTRAINT `fk_payments_order_id`
    FOREIGN KEY (`order_id`) REFERENCES `stores`.`orders` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_payments_address_id`
    FOREIGN KEY (`address_id`) REFERENCES `stores`.`addresses` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB;

CREATE UNIQUE INDEX `ui_payments_transaction_id`
  ON `stores`.`payments` (`transaction_id` ASC) VISIBLE;

-- 注文配送履歴テーブル
CREATE TABLE IF NOT EXISTS `stores`.`fulfillments` (
  `order_id`         VARCHAR(22) NOT NULL,          -- 注文履歴ID
  `address_id`       VARCHAR(22) NOT NULL,          -- 配送先情報ID
  `tracking_number`  VARCHAR(32) NULL DEFAULT NULL, -- 配送追跡番号
  `shipping_carrier` INT         NOT NULL,          -- 配送会社
  `shipping_method`  INT         NOT NULL,          -- 配送方法
  `box_size`         INT         NOT NULL,          -- 配送時の箱サイズ
  `created_at`       DATETIME    NOT NULL,          -- 登録日時
  `updated_at`       DATETIME    NOT NULL,          -- 更新日時
  `deleted_at`       DATETIME    NULL DEFAULT NULL, -- 削除日時
  PRIMARY KEY (`order_id`),
  CONSTRAINT `fk_fulfillments_order_id`
    FOREIGN KEY (`order_id`) REFERENCES `stores`.`orders` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fulfillments_address_id`
    FOREIGN KEY (`address_id`) REFERENCES `stores`.`addresses` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB;

CREATE UNIQUE INDEX `ui_fulfillments_shipping_carrier_tracking_number`
  ON `stores`.`fulfillments` (`shipping_carrier` ASC, `tracking_number` ASC) VISIBLE;

-- カートの箱管理テーブル
CREATE TABLE IF NOT EXISTS `stores`.`carts` (
  `id`          VARCHAR(22) NOT NULL,          -- カートID
  `user_id`     VARCHAR(22) NOT NULL,          -- ユーザーID
  `schedule_id` VARCHAR(22) NULL DEFAULT NULL, -- マルシェ開催スケジュールID
  `box_number`  BIGINT      NOT NULL,          -- 箱の通し番号
  `box_type`    INT         NOT NULL,          -- 箱の種別
  `box_size`    INT         NOT NULL,          -- 箱のサイズ
  `created_at`  DATETIME    NOT NULL,          -- 登録日時
  `updated_at`  DATETIME    NOT NULL,          -- 更新日時
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_carts_schedule_id`
    FOREIGN KEY (`schedule_id`) REFERENCES `stores`.`schedules` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB;

CREATE UNIQUE INDEX `ui_carts_user_id_schedule_id_box_number`
  ON `stores`.`carts` (`user_id` ASC, `schedule_id` ASC, `box_number` ASC) VISIBLE;

-- カートの商品管理テーブル
CREATE TABLE IF NOT EXISTS `stores`.`cart_items` (
  `cart_id`    VARCHAR(22) NOT NULL, -- カートID
  `product_id` VARCHAR(22) NOT NULL, -- 商品ID
  `quantity`   BIGINT      NOT NULL, -- 数量
  `created_at` DATETIME    NOT NULL, -- 登録日時
  `updated_at` DATETIME    NOT NULL, -- 更新日時
  PRIMARY KEY (`cart_id`, `product_id`),
  CONSTRAINT `fk_cart_items_cart_id`
    FOREIGN KEY (`cart_id`) REFERENCES `stores`.`carts` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_cart_items_product_id`
    FOREIGN KEY (`product_id`) REFERENCES `stores`.`products` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB;
