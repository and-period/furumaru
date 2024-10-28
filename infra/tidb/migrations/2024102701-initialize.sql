CREATE SCHEMA IF NOT EXISTS `migrations` DEFAULT CHARACTER SET utf8mb4;

CREATE TABLE `schemas` (
  `database` varchar(256) NOT NULL,
  `version` varchar(10) NOT NULL,
  `filename` varchar(256) NOT NULL,
  `created_at` int NOT NULL,
  `updated_at` int NOT NULL,
  PRIMARY KEY (`database`, `version`)
);
