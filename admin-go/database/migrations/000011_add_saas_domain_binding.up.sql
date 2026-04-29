-- SaaS 自定义域名绑定

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

CREATE TABLE IF NOT EXISTS `system_domain` (
  `id` bigint unsigned NOT NULL COMMENT '域名ID（Snowflake）',
  `domain` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '绑定域名',
  `owner_type` tinyint(1) NOT NULL DEFAULT '1' COMMENT '主体类型:1=租户,2=商户',
  `tenant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '租户',
  `merchant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '商户',
  `app_code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'admin' COMMENT '应用编码：admin/upload/shop',
  `verify_token` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '域名校验令牌',
  `verify_status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '校验状态:0=未校验,1=已校验',
  `ssl_status` tinyint(1) NOT NULL DEFAULT '0' COMMENT 'SSL状态:0=未配置,1=已配置',
  `nginx_status` tinyint(1) NOT NULL DEFAULT '0' COMMENT 'Nginx配置状态:0=未应用,1=已应用',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态:0=关闭,1=开启',
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '备注',
  `created_by` bigint unsigned DEFAULT NULL COMMENT '创建人ID',
  `dept_id` bigint unsigned DEFAULT NULL COMMENT '所属部门ID',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间，非 NULL 表示已删除',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_domain_app` (`domain`, `app_code`),
  KEY `idx_tenant_merchant` (`tenant_id`, `merchant_id`),
  KEY `idx_owner_type` (`owner_type`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='域名绑定表';

INSERT INTO `system_menu` (`id`, `parent_id`, `title`, `type`, `path`, `component`, `permission`, `icon`, `sort`, `is_show`, `is_cache`, `link_url`, `status`, `created_by`, `dept_id`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (1000000000000000100, 1000000000000000010, '域名管理', 2, '/system/domain', 'system/domain/index', 'system:domain:list', 'GlobalOutlined', 7, 1, 0, NULL, 1, 0, 1000000000000000001, '2026-04-28 00:00:00', '2026-04-28 00:00:00', NULL),
  (1000000000000000101, 1000000000000000100, '域名新增', 3, NULL, NULL, 'system:domain:create', '', 1, 0, 0, NULL, 1, 0, 1000000000000000001, '2026-04-28 00:00:00', '2026-04-28 00:00:00', NULL),
  (1000000000000000102, 1000000000000000100, '域名修改', 3, NULL, NULL, 'system:domain:update', '', 2, 0, 0, NULL, 1, 0, 1000000000000000001, '2026-04-28 00:00:00', '2026-04-28 00:00:00', NULL),
  (1000000000000000103, 1000000000000000100, '域名删除', 3, NULL, NULL, 'system:domain:delete', '', 3, 0, 0, NULL, 1, 0, 1000000000000000001, '2026-04-28 00:00:00', '2026-04-28 00:00:00', NULL),
  (1000000000000000104, 1000000000000000100, '域名批量删除', 3, NULL, NULL, 'system:domain:batch-delete', '', 4, 0, 0, NULL, 1, 0, 1000000000000000001, '2026-04-28 00:00:00', '2026-04-28 00:00:00', NULL),
  (1000000000000000105, 1000000000000000100, '应用Nginx配置', 3, NULL, NULL, 'system:domain:apply', '', 5, 0, 0, NULL, 1, 0, 1000000000000000001, '2026-04-28 00:00:00', '2026-04-28 00:00:00', NULL)
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
  (1000000000000000002, 1000000000000000100),
  (1000000000000000002, 1000000000000000101),
  (1000000000000000002, 1000000000000000102),
  (1000000000000000002, 1000000000000000103),
  (1000000000000000002, 1000000000000000104),
  (1000000000000000002, 1000000000000000105)
ON DUPLICATE KEY UPDATE `menu_id` = VALUES(`menu_id`);

SET FOREIGN_KEY_CHECKS = 1;
