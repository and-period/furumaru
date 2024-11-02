CREATE TABLE IF NOT EXISTS `stores`.`order_experiences` (
  `order_id`                 VARCHAR(22) NOT NULL,          -- 注文ID
  `experience_revision_id`   BIGINT      NOT NULL,          -- 体験ID
  `adult_count`              BIGINT      NOT NULL,          -- 大人人数
  `junior_high_school_count` BIGINT      NOT NULL,          -- 中学生人数
  `elementary_school_count`  BIGINT      NOT NULL,          -- 小学生人数
  `preschool_count`          BIGINT      NOT NULL,          -- 幼児人数
  `senior_count`             BIGINT      NOT NULL,          -- シニア人数
  `remarks`                  JSON        NULL DEFAULT NULL, -- 備考
  `created_at`               DATETIME(3) NOT NULL,          -- 登録日時
  `updated_at`               DATETIME(3) NOT NULL,          -- 更新日時
  PRIMARY KEY (`order_id`, `experience_revision_id`),
  CONSTRAINT `fk_order_experiences_order_id`
    FOREIGN KEY (`order_id`) REFERENCES `stores`.`orders` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_order_experiences_experience_revision_id`
    FOREIGN KEY (`experience_revision_id`) REFERENCES `stores`.`experience_revisions` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
);
