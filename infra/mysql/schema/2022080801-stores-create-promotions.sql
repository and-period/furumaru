-- プロモーション情報テーブル
CREATE TABLE `stores`.`promotions` (
  `id`            VARCHAR(22) NOT NULL, -- プロモーションID
  `title`         VARCHAR(64) NOT NULL, -- タイトル
  `description`   TEXT        NOT NULL, -- 内容
  `public`        TINYINT     NOT NULL, -- 公開フラグ
  `published_at`  DATETIME    NOT NULL, -- 公開日時
  `discount_type` INT         NOT NULL, -- 割引方法
  `discount_rate` BIGINT      NOT NULL, -- 割引額(%/円)
  `code`          VARCHAR(8)  NOT NULL, -- クーポンコード
  `code_type`     INT         NOT NULL, -- クーポン種別
  `start_at`      DATETIME    NOT NULL, -- クーポン使用可能日時(開始)
  `end_at`        DATETIME    NOT NULL, -- クーポン使用可能日時(終了)
  `created_at`    DATETIME    NOT NULL, -- 登録日時
  `updated_at`    DATETIME    NOT NULL, -- 更新日時
  PRIMARY KEY(`id`)
) ENGINE = InnoDB;

CREATE UNIQUE INDEX `ui_promotions_code` ON `stores`.`promotions` (`code` ASC) VISIBLE;
