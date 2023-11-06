CREATE TABLE IF NOT EXISTS `users`.`address_revisions` (
  `id`            BIGINT      NOT NULL AUTO_INCREMENT, -- 住所変更履歴ID
  `address_id`    VARCHAR(22) NOT NULL,                -- 住所ID
  `lastname`      VARCHAR(32) NOT NULL,                -- 姓
  `firstname`     VARCHAR(32) NOT NULL,                -- 名
  `postal_code`   VARCHAR(16) NOT NULL,                -- 郵便番号
  `prefecture`    INT         NOT NULL,                -- 都道府県コード
  `city`          VARCHAR(32) NOT NULL,                -- 市区町村
  `address_line1` VARCHAR(64) NOT NULL,                -- 町名・番地
  `address_line2` VARCHAR(64) NOT NULL,                -- ビル名・号室など
  `phone_number`  VARCHAR(18) NOT NULL,                -- 電話番号
  `created_at`    DATETIME    NOT NULL,                -- 登録日時
  `updated_at`    DATETIME    NOT NULL,                -- 更新日時
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_address_revisions_address_id`
    FOREIGN KEY (`address_id`) REFERENCES `users`.`addresses` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
);

ALTER TABLE `users`.`addresses` DROP INDEX `ui_addresses_hash`;

ALTER TABLE `users`.`addresses` DROP COLUMN `hash`;
ALTER TABLE `users`.`addresses` DROP COLUMN `lastname`;
ALTER TABLE `users`.`addresses` DROP COLUMN `firstname`;
ALTER TABLE `users`.`addresses` DROP COLUMN `postal_code`;
ALTER TABLE `users`.`addresses` DROP COLUMN `prefecture`;
ALTER TABLE `users`.`addresses` DROP COLUMN `city`;
ALTER TABLE `users`.`addresses` DROP COLUMN `address_line1`;
ALTER TABLE `users`.`addresses` DROP COLUMN `address_line2`;
ALTER TABLE `users`.`addresses` DROP COLUMN `phone_number`;
ALTER TABLE `users`.`addresses` DROP COLUMN `exists`;
