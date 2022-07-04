-- 配送設定情報テーブル
CREATE TABLE `stores`.`shipping` (
  `id`                  VARCHAR(22) NOT NULL,          -- 配送設定ID
  `name`                VARCHAR(64) NOT NULL,          -- 配送設定名
  `box60_rates`         JSON        NULL DEFAULT NULL, -- 箱サイズ60の通常便配送料一覧
  `box60_refrigerated`  BIGINT      NOT NULL,          -- 箱サイズ60の冷蔵便追加配送料
  `box60_frozen`        BIGINT      NOT NULL,          -- 箱サイズ60の冷凍便追加配送料
  `box80_rates`         JSON        NULL DEFAULT NULL, -- 箱サイズ60の通常便配送料一覧
  `box80_refrigerated`  BIGINT      NOT NULL,          -- 箱サイズ60の冷蔵便追加配送料
  `box80_frozen`        BIGINT      NOT NULL,          -- 箱サイズ60の冷凍便追加配送料
  `box100_rates`        JSON        NULL DEFAULT NULL, -- 箱サイズ100の通常便配送料一覧
  `box100_refrigerated` BIGINT      NOT NULL,          -- 箱サイズ100の冷蔵便追加配送料
  `box100_frozen`       BIGINT      NOT NULL,          -- 箱サイズ100の冷凍便追加配送料
  `has_free_shipping`   TINYINT     NOT NULL,          -- 送料無料オプションの有無
  `free_shipping_rates` BIGINT      NOT NULL,          -- 送料無料になる金額
  `created_at`          DATETIME    NOT NULL,          -- 登録日時
  `updated_at`          DATETIME    NOT NULL,          -- 更新日時
  PRIMARY KEY(`id`)
) ENGINE = InnoDB;
