CREATE TABLE IF NOT EXISTS `stores`.`shops` (
  `id`             VARCHAR(22) NOT NULL,          -- 店舗ID
  `coordinator_id` VARCHAR(22) NOT NULL,          -- コーディネーターID
  `name`           VARCHAR(64) NOT NULL,          -- 店舗名
  `activated`      TINYINT     NOT NULL,          -- 有効化フラグ
  `created_at`     DATETIME(3) NOT NULL,          -- 登録日時
  `updated_at`     DATETIME(3) NOT NULL,          -- 更新日時
  `deleted_at`     DATETIME(3) NULL DEFAULT NULL, -- 削除日時
  PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `stores`.`shop_producers` (
  `shop_id`     VARCHAR(22) NOT NULL, -- 店舗ID
  `producer_id` VARCHAR(22) NOT NULL, -- 生産者ID
  `created_at`  DATETIME(3) NOT NULL, -- 登録日時
  `updated_at`  DATETIME(3) NOT NULL, -- 更新日時
  PRIMARY KEY (`shop_id`, `producer_id`),
  CONSTRAINT `fk_shop_producers_shop_id`
    FOREIGN KEY (`shop_id`) REFERENCES `stores`.`shops` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
)

ALTER TABLE `stores`.`experiences` ADD COLUMN `shop_id` VARCHAR(22) NULL DEFAULT NULL;
ALTER TABLE `stores`.`orders` ADD COLUMN `shop_id` VARCHAR(22) NULL DEFAULT NULL;
ALTER TABLE `stores`.`products` ADD COLUMN `shop_id` VARCHAR(22) NULL DEFAULT NULL;
ALTER TABLE `stores`.`schedules` ADD COLUMN `shop_id` VARCHAR(22) NULL DEFAULT NULL;
ALTER TABLE `stores`.`shippings` ADD COLUMN `shop_id` VARCHAR(22) NULL DEFAULT NULL;
