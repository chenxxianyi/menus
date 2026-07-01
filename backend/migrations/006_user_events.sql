CREATE TABLE IF NOT EXISTS user_events (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  user_id BIGINT UNSIGNED NOT NULL,
  event_name VARCHAR(80) NOT NULL,
  entity_type VARCHAR(50) DEFAULT '',
  entity_id BIGINT UNSIGNED DEFAULT 0,
  payload_json JSON DEFAULT NULL,
  created_at DATETIME(3) DEFAULT NULL,
  PRIMARY KEY (id),
  KEY idx_user_events_user_id (user_id),
  KEY idx_user_events_event_name (event_name),
  KEY idx_user_events_entity_type (entity_type),
  KEY idx_user_events_entity_id (entity_id),
  KEY idx_user_events_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
