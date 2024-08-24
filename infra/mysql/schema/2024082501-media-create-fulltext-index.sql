CREATE FULLTEXT INDEX `ftx_videos` ON `media`.`videos` (`title`, `description`) WITH PARSER ngram;
