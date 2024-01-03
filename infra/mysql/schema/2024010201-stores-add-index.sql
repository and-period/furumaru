CREATE FULLTEXT INDEX `ftx_categories` ON `stores`.`categories` (`name`) WITH PARSER ngram;
CREATE FULLTEXT INDEX `ftx_products` ON `stores`.`products` (`name`, `description`) WITH PARSER ngram;
CREATE FULLTEXT INDEX `ftx_product_tags` ON `stores`.`product_tags` (`name`) WITH PARSER ngram;
CREATE FULLTEXT INDEX `ftx_product_types` ON `stores`.`product_types` (`name`) WITH PARSER ngram;
CREATE FULLTEXT INDEX `ftx_promotions` ON `stores`.`promotions` (`title`) WITH PARSER ngram;
