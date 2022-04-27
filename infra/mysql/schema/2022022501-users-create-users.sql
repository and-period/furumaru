CREATE SCHEMA IF NOT EXISTS `users` DEFAULT CHARACTER SET utf8mb4;

CREATE TABLE `users`.`users` (
  `id`            VARCHAR(22)  NOT NULL,          -- ユーザーID (Primary Key用)
  `cognito_id`    VARCHAR(36)  NOT NULL,          -- ユーザーID (Cognito)
  `provider_type` INT          NOT NULL,          -- 認証種別
  `email`         VARCHAR(256) NULL DEFAULT NULL, -- メールアドレス
  `phone_number`  VARCHAR(18)  NULL DEFAULT NULL, -- 電話番号 (国際番号(3桁)+国番号以下(15桁))
  `created_at`    DATETIME     NOT NULL,          -- 登録日時
  `updated_at`    DATETIME     NOT NULL,          -- 更新日時
  `verified_at`   DATETIME     NULL DEFAULT NULL, -- 本人確認完了日時
  `deleted_at`    DATETIME     NULL DEFAULT NULL, -- 退会日時
  PRIMARY KEY(`id`)
) ENGINE = InnoDB;

CREATE UNIQUE INDEX `ui_users_cognito_id` ON `users`.`users` (`cognito_id` ASC) VISIBLE;
CREATE UNIQUE INDEX `ui_users_email` ON `users`.`users` (`email` ASC) VISIBLE;
CREATE UNIQUE INDEX `ui_users_phone_number` ON `users`.`users` (`phone_number` ASC) VISIBLE;
