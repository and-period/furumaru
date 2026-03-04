ALTER TABLE `stores`.`ai_chat_sessions`
  MODIFY COLUMN `title` VARCHAR(256) NOT NULL DEFAULT '';

ALTER TABLE `stores`.`ai_chat_messages`
  MODIFY COLUMN `content` LONGTEXT NOT NULL;
