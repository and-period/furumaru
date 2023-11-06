CREATE TABLE IF NOT EXISTS `stores`.`product_revisions` (
  `id`            BIGINT      NOT NULL AUTO_INCREMENT, -- 商品変更履歴ID
  `product_id`    VARCHAR(22) NOT NULL,                -- 商品ID
  `price`         BIGINT      NOT NULL,                -- 販売価格
  `cost`          BIGINT      NOT NULL,                -- 商品原価
  `created_at`    DATETIME    NOT NULL,                -- 登録日時
  `updated_at`    DATETIME    NOT NULL,                -- 更新日時
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_product_revisions_product_id`
    FOREIGN KEY (`product_id`) REFERENCES `stores`.`products` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
);

ALTER TABLE `stores`.`products` DROP COLUMN `price`;
ALTER TABLE `stores`.`products` DROP COLUMN `cost`;
