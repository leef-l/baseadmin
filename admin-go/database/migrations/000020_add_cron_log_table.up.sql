-- 定时任务执行日志表

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

CREATE TABLE `system_cron_log` (
  `id` bigint unsigned NOT NULL COMMENT '日志ID（Snowflake）',
  `cron_id` bigint unsigned NOT NULL COMMENT '任务ID',
  `cron_name` varchar(100) NOT NULL COMMENT '任务名称',
  `start_at` datetime NOT NULL COMMENT '开始时间',
  `end_at` datetime DEFAULT NULL COMMENT '结束时间',
  `duration_ms` int DEFAULT '0' COMMENT '执行耗时（毫秒）',
  `result` varchar(50) NOT NULL DEFAULT 'pending' COMMENT '执行结果:success=成功,fail=失败,pending=执行中',
  `message` text COMMENT '执行消息/错误信息',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_cron_id` (`cron_id`),
  KEY `idx_result` (`result`),
  KEY `idx_start_at` (`start_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='定时任务执行日志表';

SET FOREIGN_KEY_CHECKS = 1;
