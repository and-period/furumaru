CREATE TABLE `users`.`shops` (
  `id`            VARCHAR(22)  NOT NULL,          -- ユーザーID (Primary Key用)
  `cognito_id`    VARCHAR(36)  NOT NULL,          -- ユーザーID (Cognito)
  `name`          VARCHAR(64)  NOT NULL,          -- ユーザー名 (表示用)
  `email`         VARCHAR(256) NOT NULL,          -- メールアドレス
  `created_at`    DATETIME     NOT NULL,          -- 登録日時
  `updated_at`    DATETIME     NOT NULL,          -- 更新日時
  `deleted_at`    DATETIME     NULL DEFAULT NULL, -- 削除日時
  PRIMARY KEY(`id`)
) ENGINE = InnoDB;

CREATE UNIQUE INDEX `ui_shops_cognito_id` ON `users`.`shops` (`cognito_id` ASC) VISIBLE;
CREATE UNIQUE INDEX `ui_shops_email` ON `users`.`shops` (`email` ASC) VISIBLE;
