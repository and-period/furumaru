CREATE TABLE IF NOT EXISTS `stores`.`shipping_revisions` (
  `id`                  BIGINT      NOT NULL AUTO_INCREMENT, -- 配送設定変更履歴ID
  `shipping_id`         VARCHAR(22) NOT NULL,                -- 配送設定ID
  `box60_rates`         JSON        NULL DEFAULT NULL,       -- 箱サイズ60の通常便配送料一覧
  `box60_refrigerated`  BIGINT      NOT NULL,                -- 箱サイズ60の冷蔵便追加配送料
  `box60_frozen`        BIGINT      NOT NULL,                -- 箱サイズ60の冷凍便追加配送料
  `box80_rates`         JSON        NULL DEFAULT NULL,       -- 箱サイズ60の通常便配送料一覧
  `box80_refrigerated`  BIGINT      NOT NULL,                -- 箱サイズ60の冷蔵便追加配送料
  `box80_frozen`        BIGINT      NOT NULL,                -- 箱サイズ60の冷凍便追加配送料
  `box100_rates`        JSON        NULL DEFAULT NULL,       -- 箱サイズ100の通常便配送料一覧
  `box100_refrigerated` BIGINT      NOT NULL,                -- 箱サイズ100の冷蔵便追加配送料
  `box100_frozen`       BIGINT      NOT NULL,                -- 箱サイズ100の冷凍便追加配送料
  `has_free_shipping`   TINYINT     NOT NULL,                -- 送料無料オプションの有無
  `free_shipping_rates` BIGINT      NOT NULL,                -- 送料無料になる金額
  `created_at`          DATETIME    NOT NULL,                -- 登録日時
  `updated_at`          DATETIME    NOT NULL,                -- 更新日時
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_shipping_revisions_shipping_id`
    FOREIGN KEY (`shipping_id`) REFERENCES `stores`.`shippings` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
);

ALTER TABLE `stores`.`shippings` DROP COLUMN `name`;
ALTER TABLE `stores`.`shippings` DROP COLUMN `is_default`;
ALTER TABLE `stores`.`shippings` DROP COLUMN `box60_rates`;
ALTER TABLE `stores`.`shippings` DROP COLUMN `box60_refrigerated`;
ALTER TABLE `stores`.`shippings` DROP COLUMN `box60_frozen`;
ALTER TABLE `stores`.`shippings` DROP COLUMN `box80_rates`;
ALTER TABLE `stores`.`shippings` DROP COLUMN `box80_refrigerated`;
ALTER TABLE `stores`.`shippings` DROP COLUMN `box80_frozen`;
ALTER TABLE `stores`.`shippings` DROP COLUMN `box100_rates`;
ALTER TABLE `stores`.`shippings` DROP COLUMN `box100_refrigerated`;
ALTER TABLE `stores`.`shippings` DROP COLUMN `box100_frozen`;
ALTER TABLE `stores`.`shippings` DROP COLUMN `has_free_shipping`;
ALTER TABLE `stores`.`shippings` DROP COLUMN `free_shipping_rates`;
