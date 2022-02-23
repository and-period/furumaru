CREATE SCHEMA IF NOT EXISTS `users` DEFAULT CHARACTER SET utf8mb4;

CREATE TABLE `users`.`users` (
  `id`              VARCHAR(22)  NOT NULL,              -- ユーザーID
  `email`           VARCHAR(256) NOT NULL,              -- メールアドレス
  `phone_number`    VARCHAR(18)  NOT NULL,              -- 電話番号 (国際番号(3桁)+国番号以下(15桁))
  `last_name`       VARCHAR(16)  NOT NULL,              -- 姓
  `first_name`      VARCHAR(16)  NOT NULL,              -- 名
  `last_name_kana`  VARCHAR(32)  NOT NULL,              -- 姓(かな)
  `first_name_kana` VARCHAR(32)  NOT NULL,              -- 名(かな)
  `created_at`      DATETIME     NOT NULL,              -- 登録日時
  `updated_at`      DATETIME     NOT NULL,              -- 更新日時
  `deleted_at`      DATETIME     NULL     DEFAULT NULL, -- 退会日時
  PRIMARY KEY(`id`)
) ENGINE = InnoDB;
