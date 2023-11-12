ALTER TABLE `stores`.`shipping_revisions` DROP COLUMN `box60_refrigerated`;
ALTER TABLE `stores`.`shipping_revisions` DROP COLUMN `box80_refrigerated`;
ALTER TABLE `stores`.`shipping_revisions` DROP COLUMN `box100_refrigerated`;

ALTER TABLE `stores`.`order_fulfillments` CHANGE COLUMN `shipping_method` `shipping_type` INT NOT NULL;
