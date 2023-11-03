ALTER TABLE `stores`.`payments` DROP FOREIGN KEY `fk_payments_address_id`;
ALTER TABLE `stores`.`fulfillments` DROP FOREIGN KEY `fk_fulfillments_address_id`;

DROP TABLE IF EXISTS `stores`.`addresses`;
