-- プッシュ通知テンプレート管理テーブル
CREATE TABLE `messengers`.`push_templates` (
  `id`             VARCHAR(64) NOT NULL, -- テンプレートID
  `title_template` TEXT        NOT NULL, -- テンプレート(件名)
  `body_template`  TEXT        NOT NULL, -- テンプレート(内容)
  `image_url`      TEXT        NOT NULL, -- サムネイル画像URL
  `created_at`     DATETIME    NOT NULL, -- 登録日時
  `updated_at`     DATETIME    NOT NULL, -- 更新日時
  PRIMARY KEY(`id`)
) ENGINE = InnoDB;
