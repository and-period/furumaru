-- 管理者認証情報
CREATE TABLE `users`.`admins` (
  `id`             VARCHAR(22)  NOT NULL,          -- 管理者ID
  `cognito_id`     VARCHAR(36)  NOT NULL,          -- 管理者ID (Cognito用)
  `lastname`       VARCHAR(16)  NOT NULL,          -- 姓
  `firstname`      VARCHAR(16)  NOT NULL,          -- 名
  `lastname_kana`  VARCHAR(32)  NOT NULL,          -- 姓(かな)
  `firstname_kana` VARCHAR(32)  NOT NULL,          -- 名(かな)
  `role`           INT          NOT NULL,          -- 権限
  `email`          VARCHAR(256) NOT NULL,          -- メールアドレス
  `device`         VARCHAR(256) NULL DEFAULT NULL, -- デバイスID(通知用)
  `created_at`     DATETIME     NOT NULL,          -- 登録日時
  `updated_at`     DATETIME     NOT NULL,          -- 更新日時
  PRIMARY KEY(`id`)
) ENGINE = InnoDB;

CREATE UNIQUE INDEX `ui_admins_cognito_id` ON `users`.`admins` (`cognito_id` ASC) VISIBLE;
CREATE UNIQUE INDEX `ui_admins_email` ON `users`.`admins` (`email` ASC) VISIBLE;

ALTER TABLE `users`.`administrators` ADD COLUMN `admin_id` VARCHAR(22) NOT NULL AFTER `id`;
ALTER TABLE `users`.`coordinators` ADD COLUMN `admin_id` VARCHAR(22) NOT NULL AFTER `id`;
ALTER TABLE `users`.`producers` ADD COLUMN `admin_id` VARCHAR(22) NOT NULL AFTER `id`;
