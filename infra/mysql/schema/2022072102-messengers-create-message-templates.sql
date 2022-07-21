-- メッセージテンプレート管理テーブル
CREATE TABLE `messengers`.`message_templates` (
  `id`         VARCHAR(64) NOT NULL, -- テンプレートID
  `template`   TEXT        NOT NULL, -- テンプレート
  `created_at` DATETIME    NOT NULL, -- 登録日時
  `updated_at` DATETIME    NOT NULL, -- 更新日時
  PRIMARY KEY(`id`)
) ENGINE = InnoDB;
