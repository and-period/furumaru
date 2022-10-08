-- Stripeの顧客管理テーブル
CREATE TABLE `stores`.`stripe_users` (
  `id`         VARCHAR(22) NOT NULL, -- ユーザーID(Stripe用)
  `user_id`    VARCHAR(22) NOT NULL, -- ユーザーID
  `created_at` DATETIME    NOT NULL, -- 登録日時
  `updated_at` DATETIME    NOT NULL, -- 更新日時
  PRIMARY KEY (`id`)
) ENGINE = InnoDB;

CREATE UNIQUE INDEX `ui_stripe_users_user_id` ON `stores`.`stripe_users` (`user_id` ASC) VISIBLE;

-- クレジットカード情報管理テーブル
CREATE TABLE `stores`.`payment_cards` (
  `id`             VARCHAR(22) NOT NULL, -- 決済方法ID
  `user_id`        VARCHAR(22) NOT NULL, -- ユーザーID
  `stripe_user_id` VARCHAR(22) NOT NULL, -- ユーザーID(Stripe用)
  `is_default`     TINYINT     NOT NULL, -- デフォルト決済方法
  `card_brand`     VARCHAR(16) NOT NULL, -- クレジットカード会社
  `card_exp_year`  BIGINT      NOT NULL, -- クレジットカード有効期限(年)
  `card_exp_month` BIGINT      NOT NULL, -- クレジットカード有効期限(月)
  `card_last4`     BIGINT      NOT NULL, -- クレジットカード下４桁
  `created_at`     DATETIME    NOT NULL, -- 登録日時
  `updated_at`     DATETIME    NOT NULL, -- 更新日時
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_payment_cards_stripe_user_id`
    FOREIGN KEY (`stripe_user_id`) REFERENCES `stores`.`stripe_users` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB;
