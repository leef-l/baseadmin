-- 定时任务管理表（热更新：修改表数据后无需重启）

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

CREATE TABLE `system_cron` (
  `id` bigint unsigned NOT NULL COMMENT '任务ID（Snowflake）',
  `name` varchar(100) NOT NULL COMMENT '任务名称|search:like|keyword:on',
  `expression` varchar(50) NOT NULL COMMENT 'Cron表达式',
  `handler` varchar(100) NOT NULL COMMENT '处理器名称',
  `params` text COMMENT '参数（JSON）',
  `remark` varchar(500) DEFAULT '' COMMENT '备注',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态:0=关闭,1=开启|search:select',
  `last_run_at` datetime DEFAULT NULL COMMENT '上次执行时间',
  `last_result` varchar(255) DEFAULT '' COMMENT '上次执行结果',
  `created_by` bigint unsigned DEFAULT NULL COMMENT '创建人ID',
  `dept_id` bigint unsigned DEFAULT NULL COMMENT '所属部门ID',
  `tenant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '租户',
  `merchant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '商户',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_name` (`name`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='定时任务表';

-- 内置任务：每5分钟检查过期租户自动禁用
-- 注意：GoFrame gcron 使用 6 字段格式（秒 分 时 日 月 周），不是 5 字段 Unix cron
INSERT INTO `system_cron` (`id`, `name`, `expression`, `handler`, `params`, `remark`, `status`, `tenant_id`, `merchant_id`, `created_at`, `updated_at`, `deleted_at`)
VALUES (3000000000000000001, '租户过期检查', '0 */5 * * * *', 'tenant_expire_check', '{}', '每5分钟检查租户是否过期，过期自动禁用', 1, 0, 0, NOW(), NOW(), NULL);

SET FOREIGN_KEY_CHECKS = 1;
