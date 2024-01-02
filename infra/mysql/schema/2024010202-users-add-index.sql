CREATE FULLTEXT INDEX `ftx_coordinators` ON `users`.`coordinators` (`username`, `marche_name`, `profile`) WITH PARSER ngram;
CREATE FULLTEXT INDEX `ftx_producers` ON `users`.`producers` (`username`, `profile`) WITH PARSER ngram;
