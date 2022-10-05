DROP TABLE IF EXISTS `users`.`users`;

-- 利用者情報テーブル
CREATE TABLE `users`.`users` (
  `id`         VARCHAR(22)  NOT NULL, -- 利用者ID
  `registered` TINYINT      NOT NULL, -- 会員登録フラグ
  `device`     VARCHAR(256) NOT NULL, -- デバイスID(通知用)
  `created_at` DATETIME     NOT NULL, -- 登録日時
  `updated_at` DATETIME     NOT NULL, -- 更新日時
  PRIMARY KEY (`id`)
) ENGINE = InnoDB;

-- 会員情報テーブル
CREATE TABLE IF NOT EXISTS `users`.`members` (
  `user_id`       VARCHAR(22)  NOT NULL,          -- ユーザーID
  `cognito_id`    VARCHAR(36)  NOT NULL,          -- ユーザーID(Cognito用)
  `account_id`    VARCHAR(32)  NULL DEFAULT NULL, -- ユーザーID(検索用)
  `username`      VARCHAR(32)  NOT NULL,          -- 表示名
  `provider_type` INT          NOT NULL,          -- 認証方法
  `email`         VARCHAR(256) NULL DEFAULT NULL, -- メールアドレス
  `phone_number`  VARCHAR(18)  NULL DEFAULT NULL, -- 電話番号
  `thumbnail_url` TEXT NOT     NULL,              -- サムネイルURL
  `created_at`    DATETIME     NOT NULL,          -- 登録日時
  `updated_at`    DATETIME     NOT NULL,          -- 更新日時
  `verified_at`   DATETIME     NULL DEFAULT NULL, -- 認証日時
  `deleted_at`    DATETIME     NULL DEFAULT NULL, -- 退会日時
  PRIMARY KEY (`user_id`),
  CONSTRAINT `fk_accounts_user_id`
    FOREIGN KEY (`user_id`) REFERENCES `users`.`users` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB;

CREATE UNIQUE INDEX `ui_users_cognito_id` ON `users`.`members` (`cognito_id` ASC) VISIBLE;
CREATE UNIQUE INDEX `ui_users_account_id` ON `users`.`members` (`account_id` ASC) VISIBLE;
CREATE UNIQUE INDEX `ui_users_provider_type_email` ON `users`.`members` (`provider_type` ASC, `email` ASC) VISIBLE;
CREATE UNIQUE INDEX `ui_users_provider_type_phone_number` ON `users`.`members` (`provider_type` ASC, `phone_number` ASC) VISIBLE;

-- ゲスト情報テーブル
CREATE TABLE IF NOT EXISTS `users`.`guests` (
  `user_id`      VARCHAR(22)  NOT NULL, -- ユーザーID
  `email`        VARCHAR(256) NOT NULL, -- メールアドレス
  `phone_number` VARCHAR(18)  NOT NULL, -- 電話番号
  `created_at`   DATETIME     NOT NULL, -- 登録日時
  `updated_at`   DATETIME     NOT NULL, -- 更新日時
  PRIMARY KEY (`user_id`),
  CONSTRAINT `fk_guests_user_id`
    FOREIGN KEY (`user_id`) REFERENCES `users`.`users` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB;

CREATE UNIQUE INDEX `ui_guests_email_phone_number` ON `users`.`guests` (`email` ASC, `phone_number` ASC) VISIBLE;

-- 購入者情報テーブル
CREATE TABLE IF NOT EXISTS `users`.`customers` (
  `user_id`        VARCHAR(22) NOT NULL, -- ユーザーID
  `lastname`       VARCHAR(16) NOT NULL, -- 性
  `firstname`      VARCHAR(16) NOT NULL, -- 名
  `lastname_kana`  VARCHAR(32) NOT NULL, -- 性(かな)
  `firstname_kana` VARCHAR(45) NOT NULL, -- 名
  `postal_code`    VARCHAR(16) NOT NULL, -- 郵便番号
  `prefecture`     VARCHAR(32) NOT NULL, -- 都道府県
  `city`           VARCHAR(32) NOT NULL, -- 市区町村
  `address_line1`  VARCHAR(64) NOT NULL, -- 町名・番地
  `address_line2`  VARCHAR(64) NOT NULL, -- ビル名・号室など
  `created_at`     DATETIME    NOT NULL, -- 登録日時
  `updated_at`     DATETIME    NOT NULL, -- 更新日時
  PRIMARY KEY (`user_id`),
  CONSTRAINT `fk_customer_user_id`
    FOREIGN KEY (`user_id`) REFERENCES `users`.`users` (`id`)
    ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE = InnoDB;
