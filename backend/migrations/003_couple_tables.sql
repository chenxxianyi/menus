-- 情侣绑定表
CREATE TABLE IF NOT EXISTS `couple_bindings` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `user_a_id` BIGINT NOT NULL COMMENT '发起方',
  `user_b_id` BIGINT DEFAULT 0 COMMENT '接受方',
  `couple_name` VARCHAR(50) DEFAULT '' COMMENT '情侣昵称',
  `invite_code` VARCHAR(20) NOT NULL COMMENT '邀请码',
  `status` TINYINT DEFAULT 0 COMMENT '0-待绑定 1-已绑定 2-已解绑',
  `bound_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE INDEX `idx_invite_code` (`invite_code`),
  INDEX `idx_user_a` (`user_a_id`),
  INDEX `idx_user_b` (`user_b_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 情侣点餐表
CREATE TABLE IF NOT EXISTS `couple_orders` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `couple_id` BIGINT NOT NULL COMMENT '情侣关系ID',
  `user_id` BIGINT NOT NULL COMMENT '点餐人',
  `recipe_id` BIGINT COMMENT '关联菜谱ID',
  `dish_name` VARCHAR(100) NOT NULL COMMENT '菜品名称',
  `meal_type` VARCHAR(20) COMMENT 'breakfast/lunch/dinner/snack',
  `meal_date` VARCHAR(10) COMMENT '想吃日期 YYYY-MM-DD',
  `note` VARCHAR(200) DEFAULT '' COMMENT '备注',
  `status` TINYINT DEFAULT 0 COMMENT '0-待确认 1-已确认 2-已取消',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (`couple_id`) REFERENCES `couple_bindings`(`id`),
  INDEX `idx_couple_id` (`couple_id`),
  INDEX `idx_meal_date` (`meal_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
