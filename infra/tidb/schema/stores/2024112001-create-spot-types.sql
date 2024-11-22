-- スポット種別
CREATE TABLE IF NOT EXISTS `stores`.`spot_types` (
  `id`         VARCHAR(22)  NOT NULL, -- スポット種別ID
  `name`       VARCHAR(32)  NOT NULL, -- スポット種別名
  `created_at` DATETIME(3)  NOT NULL, -- 登録日時
  `updated_at` DATETIME(3)  NOT NULL, -- 更新日時
  PRIMARY KEY (`id`)
);

CREATE UNIQUE INDEX `ui_spot_types_name` ON `stores`.`spot_types` (`name`) VISIBLE;

-- 外部キー制約の追加（既存レコードを考慮してNULLABLEにする）
ALTER TABLE `stores`.`spots` ADD COLUMN `spot_type_id` VARCHAR(22) NULL;

ALTER TABLE `stores`.`spots` ADD FOREIGN KEY `stores`.`spot_types` (`spot_type_id`)
  REFERENCES `stores`.`spot_types` (`id`);
  ON DELETE SET NULL ON UPDATE CASCADE;
