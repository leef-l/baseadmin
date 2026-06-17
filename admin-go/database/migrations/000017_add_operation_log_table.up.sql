-- 操作日志表：补齐 tenant_id / user_id / username 审计维度

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS `system_operation_log`;
CREATE TABLE `system_operation_log` (
  `id` bigint unsigned NOT NULL COMMENT '日志ID（Snowflake）',
  `module` varchar(64) NOT NULL COMMENT '模块名',
  `action` varchar(32) NOT NULL COMMENT '操作类型',
  `target_id` varchar(100) DEFAULT '' COMMENT '操作目标ID',
  `detail` text COMMENT '操作详情',
  `tenant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '租户ID',
  `user_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '操作人ID',
  `username` varchar(64) DEFAULT '' COMMENT '操作人用户名',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_module_action` (`module`, `action`),
  KEY `idx_tenant_id` (`tenant_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='操作日志表';

SET FOREIGN_KEY_CHECKS = 1;
