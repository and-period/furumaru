-- お問い合わせ種別テーブル
CREATE TABLE `messengers`.`contact_categories` (
  `id`           VARCHAR(22)  NOT NULL, -- お問い合わせ種別ID
  `title`        VARCHAR(64)  NOT NULL, -- 件名
  `created_at`   DATETIME     NOT NULL, -- 登録日時
  `updated_at`   DATETIME     NOT NULL, -- 更新日時
  PRIMARY KEY(`id`)
) ENGINE = InnoDB;

-- お問い合わせ情報テーブル
CREATE TABLE `messengers`.`contacts` (
  `id`           VARCHAR(22)  NOT NULL,          -- お問い合わせID
  `category_id`  VARCHAR(22)  NOT NULL,          -- お問い合わせ種別ID
  `title`        VARCHAR(64)  NOT NULL,          -- 件名
  `content`      TEXT         NOT NULL,          -- 内容
  `username`     VARCHAR(64)  NOT NULL,          -- 氏名
  `email`        VARCHAR(255) NOT NULL,          -- メールアドレス
  `phone_number` VARCHAR(18)  NOT NULL,          -- 電話番号
  `status`       INT          NOT NULL,          -- ステータス(不明:0, 未着手:1, 進行中:2, 対応不要:3, 対応完了:4)
  `responder_id` VARCHAR(22)  NULL DEFAULT NULL  -- 対応者ID
  `note`         TEXT         NOT NULL,          -- 対応者メモ
  `created_at`   DATETIME     NOT NULL,          -- 登録日時
  `updated_at`   DATETIME     NOT NULL,          -- 更新日時
  `deleted_at`   DATETIME     NULL DEFAULT NULL, -- 削除日時
  PRIMARY KEY(`id`),
  CONSTRAINT `fk_contacts_caterogy_id`
    FOREIGN KEY (`category_id`) REFERENCES `messengers`.`contact_categories` (`id`)
    ON DELETE SET NULL ON UPDATE CASCADE,
) ENGINE = InnoDB;

-- 会話内容テーブル
CREATE TABLE `messengers`.`threads` (
  `id`           VARCHAR(22)  NOT NULL,          -- 会話内容ID
  `contact_id`   VARCHAR(22)  NOT NULL,          -- お問い合わせID
  `user_id`      VARCHAR(22)  NULL DEFAULT NULL, -- 送信者ID
  `user_type`    INT          NOT NULL,          -- 送信者の種類(不明:0, admin:1, user:2, guest:3)
  `content`      TEXT         NOT NULL,          -- 内容
  `created_at`   DATETIME     NOT NULL,          -- 登録日時
  `updated_at`   DATETIME     NOT NULL,          -- 更新日時
  `deleted_at`   DATETIME     NULL DEFAULT NULL, -- 削除日時
  PRIMARY KEY(`id`),
  CONSTRAINT `fk_threads_contact_id`
    FOREIGN KEY (`contact_id`) REFERENCES `messengers`.`contacts` (`id`)
    ON DELETE SET NULL ON UPDATE CASCADE,
) ENGINE = InnoDB;

-- 既読管理テーブル
CREATE TABLE `messengers`.`contact_reads` (
  `id`           VARCHAR(22)  NOT NULL,          -- 既読管理ID
  `contact_id`   VARCHAR(22)  NOT NULL,          -- お問い合わせID
  `user_id`      VARCHAR(22)  NULL DEFAULT NULL, -- 送信者ID
  `user_type`    INT          NOT NULL,          -- 送信者の種類(不明:0, admin:1, user:2, guest:3)
  `read_flag`    TINYINT      NOT NULL,          -- 既読フラグ
  `created_at`   DATETIME     NOT NULL,          -- 登録日時
  `updated_at`   DATETIME     NOT NULL,          -- 更新日時
  `deleted_at`   DATETIME     NULL DEFAULT NULL, -- 削除日時
  PRIMARY KEY(`id`),
  CONSTRAINT `fk_contact_reads_contact_id`
    FOREIGN KEY (`contact_id`) REFERENCES `messengers`.`contacts` (`id`)
    ON DELETE SET NULL ON UPDATE CASCADE,
) ENGINE = InnoDB;
