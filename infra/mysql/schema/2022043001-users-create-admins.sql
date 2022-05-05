CREATE TABLE `users`.`admins` (
  `id`            VARCHAR(22)  NOT NULL,          -- ユーザーID (Primary Key用)
  `cognito_id`    VARCHAR(36)  NOT NULL,          -- ユーザーID (Cognito)
  `email`         VARCHAR(256) NOT NULL,          -- メールアドレス
  `role`          INT          NOT NULL,          -- 権限
  `created_at`    DATETIME     NOT NULL,          -- 登録日時
  `updated_at`    DATETIME     NOT NULL,          -- 更新日時
  `deleted_at`    DATETIME     NULL DEFAULT NULL, -- 削除日時
  PRIMARY KEY(`id`)
) ENGINE = InnoDB;

CREATE UNIQUE INDEX `ui_admins_cognito_id` ON `users`.`admins` (`cognito_id` ASC) VISIBLE;
CREATE UNIQUE INDEX `ui_admins_email` ON `users`.`admins` (`email` ASC) VISIBLE;
