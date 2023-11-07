DROP TABLE IF EXISTS `stores`.`payments`;
DROP TABLE IF EXISTS `stores`.`fulfillments`;
DROP TABLE IF EXISTS `stores`.`activities`;
DROP TABLE IF EXISTS `stores`.`order_items`;
DROP TABLE IF EXISTS `stores`.`orders`;

CREATE TABLE IF NOT EXISTS `stores`.`orders` (
  `id`             VARCHAR(22) NOT NULL,          -- 注文履歴ID
  `user_id`        VARCHAR(22) NOT NULL,          -- ユーザーID
  `coordinator_id` VARCHAR(22) NOT NULL,          -- コーディネータID
  `promotion_id`   VARCHAR(22) NULL DEFAULT NULL, -- プロモーションID
  `created_at`     DATETIME    NOT NULL,          -- 登録日時
  `updated_at`     DATETIME    NOT NULL,          -- 更新日時
  `deleted_at`     DATETIME    NULL DEFAULT NULL, -- 削除日時
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_orders_promotion_id`
    FOREIGN KEY (`promotion_id`) REFERENCES `stores`.`promotions` (`id`)
    ON DELETE SET NULL ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS `stores`.`order_payments` (
  `order_id`            VARCHAR(22)  NOT NULL,          -- 注文履歴ID
  `address_revision_id` BIGINT       NOT NULL,          -- 請求先情報ID
  `status`              INT          NOT NULL,          -- 決済状況
  `transaction_id`      VARCHAR(256) NOT NULL,          -- 決済ID
  `method_type`         INT          NOT NULL,          -- 決済手段
  `subtotal`            BIGINT       NOT NULL,          -- 購入金額
  `discount`            BIGINT       NOT NULL,          -- 割引金額
  `shipping_fee`        BIGINT       NOT NULL,          -- 配送手数料
  `tax`                 BIGINT       NOT NULL,          -- 消費税
  `total`               BIGINT       NOT NULL,          -- 合計金額
  `refund_total`        BIGINT       NOT NULL,          -- 返金金額
  `refund_type`         INT          NOT NULL,          -- 注文キャンセル種別
  `refund_reason`       TEXT         NOT NULL,          -- 注文キャンセル理由
  `ordered_at`          DATETIME     NULL DEFAULT NULL, -- 決済要求日時
  `paid_at`             DATETIME     NULL DEFAULT NULL, -- 決済承認日時
  `captured_at`         DATETIME     NULL DEFAULT NULL, -- 決済確定日時
  `failed_at`           DATETIME     NULL DEFAULT NULL, -- 決済失敗日時
  `refunded_at`         DATETIME     NULL DEFAULT NULL, -- 注文キャンセル日時
  `created_at`          DATETIME     NOT NULL,          -- 登録日時
  `updated_at`          DATETIME     NOT NULL,          -- 更新日時
  PRIMARY KEY (`order_id`),
  CONSTRAINT `fk_order_payments_order_id`
    FOREIGN KEY (`order_id`) REFERENCES `stores`.`orders` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS `stores`.`order_fulfillments` (
  `id`                  VARCHAR(22) NOT NULL,          -- 注文配送ID
  `order_id`            VARCHAR(22) NOT NULL,          -- 注文履歴ID
  `address_revision_id` BIGINT      NOT NULL,          -- 配送先情報ID
  `status`              INT         NOT NULL,          -- 配送状況
  `tracking_number`     VARCHAR(32) NULL DEFAULT NULL, -- 配送伝票番号
  `shipping_carrier`    INT         NOT NULL,          -- 配送会社
  `shipping_method`     INT         NOT NULL,          -- 配送方法
  `box_number`          BIGINT      NOT NULL,          -- 箱の通番
  `box_size`            INT         NOT NULL,          -- 箱の大きさ
  `shipped_at`          DATETIME    NULL DEFAULT NULL, -- 配送日時
  `created_at`          DATETIME    NOT NULL,          -- 登録日時
  `updated_at`          DATETIME    NOT NULL,          -- 更新日時
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_order_payments_order_id`
    FOREIGN KEY (`order_id`) REFERENCES `stores`.`orders` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS `stores`.`order_fulfillments` (
  `fulfillment_id`      VARCHAR(22) NOT NULL, -- 注文配送ID
  `product_revision_id` BIGINT      NOT NULL, -- 商品ID
  `order_id`            VARCHAR(22) NOT NULL, -- 注文履歴ID
  `quantity`            BIGINT      NOT NULL, -- 購入数量
  `created_at`          DATETIME    NOT NULL, -- 登録日時
  `updated_at`          DATETIME    NOT NULL, -- 更新日時
  PRIMARY KEY (`fulfillment_id`, `product_revision_id`),
  CONSTRAINT `fk_order_fulfillments_order_id`
    FOREIGN KEY (`order_id`) REFERENCES `stores`.`orders` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_order_fulfillments_fulfillment_id`
    FOREIGN KEY (`fulfillment_id`) REFERENCES `stores`.`order_fulfillments` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_order_fulfillments_product_revision_id`
    FOREIGN KEY (`product_revision_id`) REFERENCES `stores`.`product_revisions` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
);
