CREATE TABLE IF NOT EXISTS `stores`.`payment_systems` (
  `method_type` INT         NOT NULL, -- 決済種別
  `status`      INT         NOT NULL, -- 決済システム状態
  `created_at`  DATETIME(3) NOT NULL, -- 登録日時
  `updated_at`  DATETIME(3) NOT NULL, -- 更新日時
  PRIMARY KEY (`id`)
);
