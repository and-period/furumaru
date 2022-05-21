CREATE TABLE `users`.`admins` (
  `id`             VARCHAR(22)  NOT NULL,          -- ユーザーID (Primary Key用)
  `cognito_id`     VARCHAR(36)  NOT NULL,          -- ユーザーID (Cognito用)
  `lastname`       VARCHAR(16)  NOT NULL,          -- 姓
  `firstname`      VARCHAR(16)  NOT NULL,          -- 名
  `lastname_kana`  VARCHAR(32)  NOT NULL,          -- 姓(かな)
  `firstname_kana` VARCHAR(32)  NOT NULL,          -- 名(かな)
  `store_name`     VARCHAR(64)  NOT NULL,          -- 店舗名
  `thumbnail_url`  TEXT         NOT NULL,          -- サムネイルURL
  `email`          VARCHAR(256) NOT NULL,          -- メールアドレス
  `phone_number`   VARCHAR(18)  NOT NULL,          -- 電話番号
  `postal_code`    VARCHAR(16)  NOT NULL,          -- 郵便番号
  `prefecture`     VARCHAR(32)  NOT NULL,          -- 都道府県
  `city`           VARCHAR(32)  NOT NULL,          -- 市区町村
  `address_line1`  VARCHAR(64)  NOT NULL,          -- 町名・番地
  `address_line2`  VARCHAR(64)  NOT NULL,          -- ビル名・号室など
  `role`           INT          NOT NULL,          -- 権限
  `created_at`     DATETIME     NOT NULL,          -- 登録日時
  `updated_at`     DATETIME     NOT NULL,          -- 更新日時
  `deleted_at`     DATETIME     NULL DEFAULT NULL, -- 退会日時
  PRIMARY KEY(`id`)
) ENGINE = InnoDB;

CREATE UNIQUE INDEX `ui_admins_cognito_id` ON `users`.`admins` (`cognito_id` ASC) VISIBLE;
CREATE UNIQUE INDEX `ui_admins_email` ON `users`.`admins` (`email` ASC) VISIBLE;
CREATE INDEX `idx_admins_role` ON `users`.`admins` (`role` ASC) VISIBLE;
