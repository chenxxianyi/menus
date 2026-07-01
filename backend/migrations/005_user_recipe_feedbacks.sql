CREATE TABLE IF NOT EXISTS user_recipe_feedbacks (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  user_id BIGINT UNSIGNED NOT NULL,
  recipe_id BIGINT UNSIGNED NOT NULL,
  feedback_type VARCHAR(20) NOT NULL,
  source VARCHAR(30) DEFAULT '',
  created_at DATETIME(3) DEFAULT NULL,
  updated_at DATETIME(3) DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY idx_user_recipe_feedback (user_id, recipe_id, feedback_type),
  KEY idx_user_recipe_feedback_user_id (user_id),
  KEY idx_user_recipe_feedback_recipe_id (recipe_id),
  KEY idx_user_recipe_feedback_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
