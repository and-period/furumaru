-- 開催スケジュールテーブルの修正
ALTER TABLE `stores`.`schedules` ADD COLUMN `coordinator_id` VARCHAR(22) NOT NULL AFTER `id`;

-- ライブ配信テーブルの修正
ALTER TABLE `stores`.`lives` DROP COLUMN `recommends`;
ALTER TABLE `stores`.`lives` ADD COLUMN `published` TINYINT NOT NULL AFTER `end_at`;

-- ライブ関連商品テーブル
CREATE TABLE `stores`.`live_products` (
  `live_id`    VARCHAR(22) NOT NULL,          -- ライブ配信ID
  `product_id` VARCHAR(22) NOT NULL,          -- 商品ID
  `created_at` DATETIME    NOT NULL,          -- 登録日時
  `updated_at` DATETIME    NOT NULL,          -- 更新日時
  `deleted_at` DATETIME    NULL DEFAULT NULL, -- 削除日時
  PRIMARY KEY(`live_id`, `product_id`),
  CONSTRAINT `fk_live_products_live_id`
    FOREIGN KEY (`live_id`) REFERENCES `stores`.`lives` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_live_products_product_id`
    FOREIGN KEY (`product_id`) REFERENCES `stores`.`products` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB;
