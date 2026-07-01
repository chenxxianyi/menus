ALTER TABLE `menus`
  ADD INDEX `idx_menus_user_created` (`user_id`, `created_at`);
