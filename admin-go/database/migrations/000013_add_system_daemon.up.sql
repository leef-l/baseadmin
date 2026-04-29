-- 超级管理员守护进程管理

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

CREATE TABLE IF NOT EXISTS `system_daemon` (
  `id` bigint unsigned NOT NULL COMMENT '守护进程ID（Snowflake）',
  `name` varchar(80) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '显示名称',
  `program` varchar(80) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Supervisor进程名',
  `command` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '启动命令',
  `directory` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '运行目录',
  `run_user` varchar(80) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'root' COMMENT '运行用户',
  `numprocs` int unsigned NOT NULL DEFAULT '1' COMMENT '进程数量',
  `priority` int unsigned NOT NULL DEFAULT '999' COMMENT '启动优先级',
  `autostart` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否随Supervisor启动',
  `autorestart` tinyint(1) NOT NULL DEFAULT '1' COMMENT '异常退出是否自动重启',
  `startsecs` int unsigned NOT NULL DEFAULT '3' COMMENT '启动稳定秒数',
  `startretries` int unsigned NOT NULL DEFAULT '3' COMMENT '启动重试次数',
  `stop_signal` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'QUIT' COMMENT '停止信号',
  `environment` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '环境变量，Supervisor environment格式',
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '备注',
  `created_by` bigint unsigned DEFAULT NULL COMMENT '创建人ID',
  `dept_id` bigint unsigned DEFAULT NULL COMMENT '所属部门ID',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间，非 NULL 表示已删除',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_program` (`program`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='守护进程配置表';

INSERT INTO `system_menu` (`id`, `parent_id`, `title`, `type`, `path`, `component`, `permission`, `icon`, `sort`, `is_show`, `is_cache`, `link_url`, `status`, `created_by`, `dept_id`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (1000000000000000110, 1000000000000000010, '守护进程', 2, '/system/daemon', 'system/daemon/index', 'system:daemon:list', 'ControlOutlined', 8, 1, 0, NULL, 1, 0, 1000000000000000001, '2026-04-28 00:00:00', '2026-04-28 00:00:00', NULL),
  (1000000000000000111, 1000000000000000110, '守护进程新增', 3, NULL, NULL, 'system:daemon:create', '', 1, 0, 0, NULL, 1, 0, 1000000000000000001, '2026-04-28 00:00:00', '2026-04-28 00:00:00', NULL),
  (1000000000000000112, 1000000000000000110, '守护进程修改', 3, NULL, NULL, 'system:daemon:update', '', 2, 0, 0, NULL, 1, 0, 1000000000000000001, '2026-04-28 00:00:00', '2026-04-28 00:00:00', NULL),
  (1000000000000000113, 1000000000000000110, '守护进程删除', 3, NULL, NULL, 'system:daemon:delete', '', 3, 0, 0, NULL, 1, 0, 1000000000000000001, '2026-04-28 00:00:00', '2026-04-28 00:00:00', NULL),
  (1000000000000000114, 1000000000000000110, '守护进程批量删除', 3, NULL, NULL, 'system:daemon:batch-delete', '', 4, 0, 0, NULL, 1, 0, 1000000000000000001, '2026-04-28 00:00:00', '2026-04-28 00:00:00', NULL),
  (1000000000000000115, 1000000000000000110, '守护进程重启', 3, NULL, NULL, 'system:daemon:restart', '', 5, 0, 0, NULL, 1, 0, 1000000000000000001, '2026-04-28 00:00:00', '2026-04-28 00:00:00', NULL),
  (1000000000000000116, 1000000000000000110, '守护进程暂停', 3, NULL, NULL, 'system:daemon:stop', '', 6, 0, 0, NULL, 1, 0, 1000000000000000001, '2026-04-28 00:00:00', '2026-04-28 00:00:00', NULL),
  (1000000000000000117, 1000000000000000110, '守护进程查看', 3, NULL, NULL, 'system:daemon:view', '', 7, 0, 0, NULL, 1, 0, 1000000000000000001, '2026-04-28 00:00:00', '2026-04-28 00:00:00', NULL)
ON DUPLICATE KEY UPDATE
  `parent_id` = VALUES(`parent_id`),
  `title` = VALUES(`title`),
  `type` = VALUES(`type`),
  `path` = VALUES(`path`),
  `component` = VALUES(`component`),
  `permission` = VALUES(`permission`),
  `icon` = VALUES(`icon`),
  `sort` = VALUES(`sort`),
  `is_show` = VALUES(`is_show`),
  `status` = VALUES(`status`),
  `deleted_at` = NULL;

INSERT INTO `system_role_menu` (`role_id`, `menu_id`) VALUES
  (1000000000000000002, 1000000000000000110),
  (1000000000000000002, 1000000000000000111),
  (1000000000000000002, 1000000000000000112),
  (1000000000000000002, 1000000000000000113),
  (1000000000000000002, 1000000000000000114),
  (1000000000000000002, 1000000000000000115),
  (1000000000000000002, 1000000000000000116),
  (1000000000000000002, 1000000000000000117)
ON DUPLICATE KEY UPDATE `menu_id` = VALUES(`menu_id`);

SET FOREIGN_KEY_CHECKS = 1;
