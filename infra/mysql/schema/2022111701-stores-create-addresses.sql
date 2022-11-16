-- 住所管理テーブル
CREATE TABLE `stores`.`addresses` (
  `id`              VARCHAR(22) NOT NULL,          -- 住所ID
  `user_id`         VARCHAR(22) NOT NULL,          -- ユーザーID
  `hash`            TEXT        NOT NULL,          -- 重複登録抑止用
  `is_default`      TINYINT     NOT NULL,          -- デフォルト設定フラグ
  `lastname`        VARCHAR(32) NOT NULL,          -- 姓
  `firstname`       VARCHAR(32) NOT NULL,          -- 名
  `postal_code`     VARCHAR(16) NOT NULL,          -- 郵便番号
  `prefecture`      VARCHAR(32) NOT NULL,          -- 都道府県
  `prefecture_code` INT         NOT NULL,          -- 都道府県コード
  `city`            VARCHAR(32) NOT NULL,          -- 市区町村
  `address_line1`   VARCHAR(64) NOT NULL,          -- 町名・番地
  `address_line2`   VARCHAR(64) NOT NULL,          -- ビル名・号室など
  `phone_number`    VARCHAR(18) NOT NULL,          -- 電話番号
  `exists`          TINYINT     NULL DEFAULT 1,    -- 存在検証用フラグ
  `created_at`      DATETIME    NOT NULL,          -- 登録日時
  `updated_at`      DATETIME    NOT NULL,          -- 更新日時
  `deleted_at`      DATETIME    NULL DEFAULT NULL, -- 削除日時
  PRIMARY KEY (`id`)
) ENGINE = InnoDB;

CREATE UNIQUE INDEX `ui_addresses_hash` (`exists` DESC, `hash` ASC) VISIBLE;
