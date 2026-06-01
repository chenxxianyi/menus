-- 修复 couple_bindings 表的 user_b_id 允许 NULL
-- 邀请阶段 user_b_id 还没有值，需要允许为空

ALTER TABLE `couple_bindings`
  MODIFY COLUMN `user_b_id` BIGINT DEFAULT NULL;

-- 删除旧的外键约束（名称可能不同，根据实际情况调整）
-- 先查看约束名称: SHOW CREATE TABLE couple_bindings;
-- 然后删除: ALTER TABLE couple_bindings DROP FOREIGN KEY fk_couple_bindings_user_b;

-- 如果 AutoMigrate 能自动处理，也可以直接删表重建：
-- DROP TABLE IF EXISTS `couple_orders`;
-- DROP TABLE IF EXISTS `couple_bindings`;
