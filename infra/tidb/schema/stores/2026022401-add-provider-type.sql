ALTER TABLE `stores`.`payment_systems`
  ADD COLUMN `provider_type` int NOT NULL DEFAULT 1 AFTER `method_type`;

ALTER TABLE `stores`.`order_payments`
  ADD COLUMN `provider_type` int NOT NULL DEFAULT 1 AFTER `address_revision_id`;
