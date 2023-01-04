ALTER TABLE `stores`.`payment_cards` CHANGE COLUMN `card_brand` `brand` VARCHAR(16) NOT NULL;
ALTER TABLE `stores`.`payment_cards` CHANGE COLUMN `card_exp_year` `exp_year` BIGINT NOT NULL;
ALTER TABLE `stores`.`payment_cards` CHANGE COLUMN `card_exp_month` `exp_month` BIGINT NOT NULL;
ALTER TABLE `stores`.`payment_cards` CHANGE COLUMN `card_last4` `last4` BIGINT NOT NULL;
