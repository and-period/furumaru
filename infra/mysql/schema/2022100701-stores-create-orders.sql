-- 注文履歴マスタテーブル
CREATE TABLE `stores`.`orders` (
  `id`                 VARCHAR(22)  NOT NULL,          -- 注文ID
  `user_id`            VARCHAR(22)  NOT NULL,          -- 利用者ID
  `payment_status`     INT          NOT NULL,          -- 支払いステータス
  `fulfillment_status` INT          NOT NULL,          -- 配送ステータス
  `cancel_type`        INT          NOT NULL,          -- 注文キャンセル種別
  `cancel_reason`      VARCHAR(256) NOT NULL,          -- 注文キャンセル理由
  `ordered_at`         DATETIME     NULL DEFAULT NULL, -- 注文要求日時
  `confirmed_at`       DATETIME     NULL DEFAULT NULL, -- 注文実行日時
  `captured_at`        DATETIME     NULL DEFAULT NULL, -- 注文確定日時
  `canceled_at`        DATETIME     NULL DEFAULT NULL, -- 注文キャンセル日時
  `delivered_at`       DATETIME     NULL DEFAULT NULL, -- 配送日時
  `created_at`         DATETIME     NOT NULL,          -- 登録日時
  `updated_at`         DATETIME     NOT NULL,          -- 更新日時
  PRIMARY KEY (`id`)
) ENGINE = InnoDB;

-- 支払い管理テーブル
CREATE TABLE `stores`.`order_payments` (
  `id`                VARCHAR(22) NOT NULL,          -- 支払いID
  `transaction_id`    VARCHAR(22) NOT NULL,          -- 支払いID(Stripe用)
  `order_id`          VARCHAR(22) NOT NULL,          -- 注文履歴ID
  `promotion_id`      VARCHAR(22) NULL DEFAULT NULL, -- プロモーションID
  `payment_id`        VARCHAR(22) NULL DEFAULT NULL, -- 決済手段ID
  `payment_type`      INT         NOT NULL,          -- 決済手段
  `subtotal`          BIGINT      NOT NULL,          -- 購入金額
  `discount`          BIGINT      NOT NULL,          -- 割引額
  `shipping_charge`   BIGINT      NOT NULL,          -- 配送料
  `tax`               BIGINT      NOT NULL,          -- 消費税
  `total`             BIGINT      NOT NULL,          -- 支払い合計金額
  `lastname`          VARCHAR(32) NOT NULL,          -- 請求先情報 姓
  `firstname`         VARCHAR(32) NOT NULL,          -- 請求先情報 名
  `postal_code`       VARCHAR(16) NOT NULL,          -- 請求先情報 郵便番号
  `prefecture`        VARCHAR(32) NOT NULL,          -- 請求先情報 都道府県
  `city`              VARCHAR(32) NOT NULL,          -- 請求先情報 市区町村
  `address_line1`     VARCHAR(64) NOT NULL,          -- 請求先情報 町名・番地
  `address_line2`     VARCHAR(64) NOT NULL,          -- 請求先情報 ビル名・号室など
  `phone_number`      VARCHAR(18) NOT NULL,          -- 請求先情報 電話番号
  `created_at`        DATETIME    NOT NULL,          -- 登録日時
  `updated_at`        DATETIME    NOT NULL,          -- 更新日時
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_order_payments_order_id`
    FOREIGN KEY (`order_id`) REFERENCES `stores`.`orders` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_order_payments_promotion_id`
    FOREIGN KEY (`promotion_id`) REFERENCES `stores`.`promotions` (`id`)
    ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE = InnoDB;

CREATE UNIQUE INDEX `ui_order_payments_transaction_id` ON `stores`.`order_payments` (`transaction_id` ASC) VISIBLE;

-- 配送管理テーブル
CREATE TABLE `stores`.`order_fulfillments` (
  `id`               VARCHAR(22) NOT NULL,          -- 配送ID
  `order_id`         VARCHAR(22) NOT NULL,          -- 注文履歴ID
  `shipping_id`      VARCHAR(22) NOT NULL,          -- 配送設定ID
  `tracking_number`  VARCHAR(32) NULL DEFAULT NULL, -- 伝票番号
  `shipping_carrier` INT         NOT NULL,          -- 配送会社
  `shipping_method`  INT         NOT NULL,          -- 配送方法
  `box_size`         INT         NOT NULL,          -- 箱の大きさ
  `box_count`        BIGINT      NOT NULL,          -- 箱の個数
  `weight_total`     BIGINT      NOT NULL,          -- 総重量(g)
  `lastname`         VARCHAR(32) NOT NULL,          -- 配送先情報 姓
  `firstname`        VARCHAR(32) NOT NULL,          -- 配送先情報 名
  `postal_code`      VARCHAR(16) NOT NULL,          -- 配送先情報 郵便番号
  `prefecture`       VARCHAR(32) NOT NULL,          -- 配送先情報 都道府県
  `city`             VARCHAR(32) NOT NULL,          -- 配送先情報 市区町村
  `address_line1`    VARCHAR(64) NOT NULL,          -- 配送先情報 町名・番地
  `address_line2`    VARCHAR(64) NOT NULL,          -- 配送先情報 ビル名・号室など
  `phone_number`     VARCHAR(18) NOT NULL,          -- 配送先情報 電話番号
  `created_at`       DATETIME    NOT NULL,          -- 登録日時
  `updated_at`       DATETIME    NOT NULL,          -- 更新日時
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_order_fulfillments_order_id`
    FOREIGN KEY (`order_id`) REFERENCES `stores`.`orders` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_order_fulfillments_shipping_id`
    FOREIGN KEY (`shipping_id`) REFERENCES `stores`.`shippings` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB;

CREATE UNIQUE INDEX `ui_order_fulfillments_tracking_number` ON `stores`.`order_fulfillments` (`tracking_number` ASC) VISIBLE;

-- 注文商品管理テーブル
CREATE TABLE `stores`.`order_items` (
  `id`          VARCHAR(22) NOT NULL, -- 注文詳細ID
  `order_id`    VARCHAR(22) NOT NULL, -- 注文履歴ID
  `product_id`  VARCHAR(22) NOT NULL, -- 商品ID
  `price`       BIGINT      NOT NULL, -- 購入価格
  `quantity`    BIGINT      NOT NULL, -- 購入数量
  `weight`      BIGINT      NOT NULL, -- 商品重量
  `weight_unit` INT         NOT NULL, -- 商品重量単位
  `created_at`  DATETIME    NOT NULL, -- 登録日時
  `updated_at`  DATETIME    NOT NULL, -- 更新日時
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_order_items_order_id`
    FOREIGN KEY (`order_id`) REFERENCES `stores`.`orders` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_order_items_product_id`
    FOREIGN KEY (`product_id`) REFERENCES `stores`.`products` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB;

CREATE UNIQUE INDEX `ui_order_items_order_id_product_id` ON `stores`.`order_items` (`order_id` ASC, `product_id` ASC) VISIBLE;

-- 注文イベントログ管理テーブル
CREATE TABLE `stores`.`order_activities` (
  `id`         VARCHAR(22) NOT NULL,          -- イベントログID
  `order_id`   VARCHAR(22) NOT NULL,          -- 注文履歴ID
  `user_id`    VARCHAR(22) NOT NULL,          -- ユーザーID
  `event_type` INT         NOT NULL,          -- イベントログ種別
  `detail`     INT         NOT NULL,          -- イベントログ詳細
  `metadata`   JSON        NULL DEFAULT NULL, -- メタデータ
  `created_at` DATETIME    NOT NULL,          -- 登録日時
  `updated_at` DATETIME    NOT NULL,          -- 更新日時
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_order_activities_order_id`
    FOREIGN KEY (`order_id`) REFERENCES `stores`.`orders` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB;
