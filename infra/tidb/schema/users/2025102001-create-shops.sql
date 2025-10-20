CREATE TABLE IF NOT EXISTS `users`.`shops` (
  `id`               VARCHAR(22) NOT NULL,          -- 店舗ID
  `coordinator_id`   VARCHAR(22) NOT NULL,          -- コーディネーターID
  `name`             VARCHAR(64) NOT NULL,          -- 店舗名
  `activated`        TINYINT     NOT NULL,          -- 有効化フラグ
  `product_type_ids` JSON        NULL DEFAULT NULL, -- 取扱商品タイプID一覧
  `business_days`    JSON        NULL DEFAULT NULL, -- 営業日
  `created_at`       DATETIME(3) NOT NULL,          -- 登録日時
  `updated_at`       DATETIME(3) NOT NULL,          -- 更新日時
  `deleted_at`       DATETIME(3) NULL DEFAULT NULL, -- 削除日時
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_shops_coordinator_id`
    FOREIGN KEY (`coordinator_id`) REFERENCES `users`.`coordinators` (`admin_id`)
    ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS `users`.`shop_producers` (
  `shop_id`     VARCHAR(22) NOT NULL, -- 店舗ID
  `producer_id` VARCHAR(22) NOT NULL, -- 生産者ID
  `created_at`  DATETIME(3) NOT NULL, -- 登録日時
  `updated_at`  DATETIME(3) NOT NULL, -- 更新日時
  PRIMARY KEY (`shop_id`, `producer_id`),
  CONSTRAINT `fk_shop_producers_shop_id`
    FOREIGN KEY (`shop_id`) REFERENCES `users`.`shops` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_shop_producers_producer_id`
    FOREIGN KEY (`producer_id`) REFERENCES `users`.`producers` (`admin_id`)
    ON DELETE CASCADE ON UPDATE CASCADE
);
