-- お問い合わせ情報テーブル
CREATE TABLE `messengers`.`contacts` (
  `id`           VARCHAR(22)  NOT NULL, -- お問い合わせID
  `title`        VARCHAR(64)  NOT NULL, -- 件名
  `content`      TEXT         NOT NULL, -- 内容
  `username`     VARCHAR(64)  NOT NULL, -- 氏名
  `email`        VARCHAR(255) NOT NULL, -- メールアドレス
  `phone_number` VARCHAR(18)  NOT NULL, -- 電話番号
  `priority`     INT          NOT NULL, -- 優先度
  `status`       INT          NOT NULL, -- ステータス
  `note`         TEXT         NOT NULL, -- 対応者メモ
  `created_at`   DATETIME     NOT NULL, -- 登録日時
  `updated_at`   DATETIME     NOT NULL, -- 更新日時
  PRIMARY KEY(`id`)
) ENGINE = InnoDB;
