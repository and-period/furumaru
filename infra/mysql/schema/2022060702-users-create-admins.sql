DROP TABLE IF EXISTS `users`.`admins`;

-- 管理者認証情報
CREATE TABLE `users`.`admin_auths` (
  `admin_id`   VARCHAR(22) NOT NULL, -- 管理者ID
  `cognito_id` VARCHAR(22) NOT NULL, -- 管理者ID (Cognito用)
  `role`       INT         NOT NULL, -- 権限
  `created_at` DATETIME    NOT NULL, -- 登録日時
  `updated_at` DATETIME    NOT NULL, -- 更新日時
  PRIMARY KEY(`admin_id`)
) ENGINE = InnoDB;

CREATE UNIQUE INDEX `ui_admin_auths_cognito_id` ON `users`.`admin_auths` (`cognito_id` ASC) VISIBLE;

-- システム管理者情報
CREATE TABLE `users`.`administrators` (
  `id`             VARCHAR(22)  NOT NULL,          -- システム管理者ID
  `lastname`       VARCHAR(16)  NOT NULL,          -- 姓
  `firstname`      VARCHAR(16)  NOT NULL,          -- 名
  `lastname_kana`  VARCHAR(32)  NOT NULL,          -- 姓(かな)
  `firstname_kana` VARCHAR(32)  NOT NULL,          -- 名(かな)
  `email`          VARCHAR(256) NOT NULL,          -- メールアドレス
  `phone_number`   VARCHAR(18)  NOT NULL,          -- 電話番号
  `created_at`     DATETIME     NOT NULL,          -- 登録日時
  `updated_at`     DATETIME     NOT NULL,          -- 更新日時
  `deleted_at`     DATETIME     NULL DEFAULT NULL, -- 退会日時
  PRIMARY KEY(`id`)
) ENGINE = InnoDB;

CREATE UNIQUE INDEX `ui_administrators_email` ON `users`.`administrators` (`email` ASC) VISIBLE;

-- コーディネータ情報
CREATE TABLE `users`.`coordinators` (
  `id`                VARCHAR(22)  NOT NULL,          -- コーディネータID
  `lastname`          VARCHAR(16)  NOT NULL,          -- 姓
  `firstname`         VARCHAR(16)  NOT NULL,          -- 名
  `lastname_kana`     VARCHAR(32)  NOT NULL,          -- 姓(かな)
  `firstname_kana`    VARCHAR(32)  NOT NULL,          -- 名(かな)
  `company_name`      VARCHAR(64)  NOT NULL,          -- 会社名
  `store_name`        VARCHAR(64)  NOT NULL,          -- 店舗名
  `thumbnail_url`     TEXT         NOT NULL,          -- サムネイルURL
  `header_url`        TEXT         NOT NULL,          -- ヘッダー画像URL
  `twitter_account`   VARCHAR(15)  NOT NULL,          -- SNS(Twitter)アカウント名
  `instagram_account` VARCHAR(30)  NOT NULL,          -- SNS(Instagram)アカウント名
  `facebook_account`  VARCHAR(50)  NOT NULL,          -- SNS(Facebook)アカウント名
  `email`             VARCHAR(256) NOT NULL,          -- メールアドレス
  `phone_number`      VARCHAR(18)  NOT NULL,          -- 電話番号
  `postal_code`       VARCHAR(16)  NOT NULL,          -- 郵便番号
  `prefecture`        VARCHAR(32)  NOT NULL,          -- 都道府県
  `city`              VARCHAR(32)  NOT NULL,          -- 市区町村
  `address_line1`     VARCHAR(64)  NOT NULL,          -- 町名・番地
  `address_line2`     VARCHAR(64)  NOT NULL,          -- ビル名・号室など
  `created_at`        DATETIME     NOT NULL,          -- 登録日時
  `updated_at`        DATETIME     NOT NULL,          -- 更新日時
  `deleted_at`        DATETIME     NULL DEFAULT NULL, -- 退会日時
  PRIMARY KEY(`id`)
) ENGINE = InnoDB;

CREATE UNIQUE INDEX `ui_coordinators_email` ON `users`.`coordinators` (`email` ASC) VISIBLE;

-- 生産者情報
CREATE TABLE `users`.`producers` (
  `id`             VARCHAR(22)  NOT NULL,          -- 生産者ID
  `lastname`       VARCHAR(16)  NOT NULL,          -- 姓
  `firstname`      VARCHAR(16)  NOT NULL,          -- 名
  `lastname_kana`  VARCHAR(32)  NOT NULL,          -- 姓(かな)
  `firstname_kana` VARCHAR(32)  NOT NULL,          -- 名(かな)
  `store_name`     VARCHAR(64)  NOT NULL,          -- 店舗名
  `thumbnail_url`  TEXT         NOT NULL,          -- サムネイルURL
  `header_url`     TEXT         NOT NULL,          -- ヘッダー画像URL
  `email`          VARCHAR(256) NOT NULL,          -- メールアドレス
  `phone_number`   VARCHAR(18)  NOT NULL,          -- 電話番号
  `postal_code`    VARCHAR(16)  NOT NULL,          -- 郵便番号
  `prefecture`     VARCHAR(32)  NOT NULL,          -- 都道府県
  `city`           VARCHAR(32)  NOT NULL,          -- 市区町村
  `address_line1`  VARCHAR(64)  NOT NULL,          -- 町名・番地
  `address_line2`  VARCHAR(64)  NOT NULL,          -- ビル名・号室など
  `created_at`     DATETIME     NOT NULL,          -- 登録日時
  `updated_at`     DATETIME     NOT NULL,          -- 更新日時
  `deleted_at`     DATETIME     NULL DEFAULT NULL, -- 退会日時
  PRIMARY KEY(`id`)
) ENGINE = InnoDB;

CREATE UNIQUE INDEX `ui_producers_email` ON `users`.`producers` (`email` ASC) VISIBLE;
CREATE UNIQUE INDEX `ui_producers_store_name` ON `users`.`producers` (`store_name` ASC) VISIBLE;
