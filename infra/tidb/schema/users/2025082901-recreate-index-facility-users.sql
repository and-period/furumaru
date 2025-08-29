ALTER TABLE `users`.`facility_users` DROP UNIQUE KEY `ui_guests_email_producer_id`;

ALTER TABLE `users`.`facility_users` DROP FOREIGN KEY `fk_guests_user_id`;
ALTER TABLE `users`.`facility_users` DROP FOREIGN KEY `fk_guests_producer_id`;

ALTER TABLE `users`.`facility_users`
  ADD UNIQUE KEY `ui_facility_users_provider_type_external_id_producer_id`
    (`exists` DESC, `provider_type`, `external_id`, `producer_id`),
  ADD CONSTRAINT `fk_facility_users_user_id`
    FOREIGN KEY (`user_id`)
    REFERENCES `users`.`users` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `fk_facility_users_producer_id`
    FOREIGN KEY (`producer_id`)
    REFERENCES `users`.`producers` (`admin_id`)
    ON DELETE CASCADE ON UPDATE CASCADE;
