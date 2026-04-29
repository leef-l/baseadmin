-- SaaS 多租户 / 多商户基础结构

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

CREATE TABLE `system_tenant` (
  `id` bigint unsigned NOT NULL COMMENT '租户ID（Snowflake）',
  `name` varchar(80) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '租户名称',
  `code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '租户编码',
  `contact_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '联系人',
  `contact_phone` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '联系电话',
  `domain` varchar(120) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '租户域名',
  `expire_at` datetime DEFAULT NULL COMMENT '到期时间',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态:0=关闭,1=开启',
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '备注',
  `created_by` bigint unsigned DEFAULT NULL COMMENT '创建人ID',
  `dept_id` bigint unsigned DEFAULT NULL COMMENT '所属部门ID',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间，非 NULL 表示已删除',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='租户表';

CREATE TABLE `system_merchant` (
  `id` bigint unsigned NOT NULL COMMENT '商户ID（Snowflake）',
  `tenant_id` bigint unsigned NOT NULL COMMENT '租户',
  `name` varchar(80) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '商户名称',
  `code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '商户编码',
  `contact_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '联系人',
  `contact_phone` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '联系电话',
  `address` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '商户地址',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态:0=关闭,1=开启',
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '备注',
  `created_by` bigint unsigned DEFAULT NULL COMMENT '创建人ID',
  `dept_id` bigint unsigned DEFAULT NULL COMMENT '所属部门ID',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间，非 NULL 表示已删除',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_tenant_code` (`tenant_id`, `code`),
  KEY `idx_tenant_id` (`tenant_id`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商户表';

ALTER TABLE `system_dept`
  ADD COLUMN `tenant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '租户' AFTER `dept_id`,
  ADD COLUMN `merchant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '商户' AFTER `tenant_id`,
  ADD KEY `idx_tenant_merchant` (`tenant_id`, `merchant_id`);

ALTER TABLE `system_role`
  ADD COLUMN `tenant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '租户' AFTER `dept_id`,
  ADD COLUMN `merchant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '商户' AFTER `tenant_id`,
  ADD KEY `idx_tenant_merchant` (`tenant_id`, `merchant_id`);

ALTER TABLE `system_users`
  ADD COLUMN `tenant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '租户' AFTER `dept_id`,
  ADD COLUMN `merchant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '商户' AFTER `tenant_id`,
  ADD KEY `idx_tenant_merchant` (`tenant_id`, `merchant_id`);

INSERT INTO `system_menu` (`id`, `parent_id`, `title`, `type`, `path`, `component`, `permission`, `icon`, `sort`, `is_show`, `is_cache`, `link_url`, `status`, `created_by`, `dept_id`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (1000000000000000080, 1000000000000000010, '租户管理', 2, '/system/tenant', 'system/tenant/index', 'system:tenant:list', 'ClusterOutlined', 5, 1, 0, NULL, 1, 0, 1000000000000000001, '2026-04-27 00:00:00', '2026-04-27 00:00:00', NULL),
  (1000000000000000081, 1000000000000000080, '租户新增', 3, NULL, NULL, 'system:tenant:create', '', 1, 0, 0, NULL, 1, 0, 1000000000000000001, '2026-04-27 00:00:00', '2026-04-27 00:00:00', NULL),
  (1000000000000000082, 1000000000000000080, '租户修改', 3, NULL, NULL, 'system:tenant:update', '', 2, 0, 0, NULL, 1, 0, 1000000000000000001, '2026-04-27 00:00:00', '2026-04-27 00:00:00', NULL),
  (1000000000000000083, 1000000000000000080, '租户删除', 3, NULL, NULL, 'system:tenant:delete', '', 3, 0, 0, NULL, 1, 0, 1000000000000000001, '2026-04-27 00:00:00', '2026-04-27 00:00:00', NULL),
  (1000000000000000084, 1000000000000000080, '租户批量删除', 3, NULL, NULL, 'system:tenant:batch-delete', '', 4, 0, 0, NULL, 1, 0, 1000000000000000001, '2026-04-27 00:00:00', '2026-04-27 00:00:00', NULL),
  (1000000000000000090, 1000000000000000010, '商户管理', 2, '/system/merchant', 'system/merchant/index', 'system:merchant:list', 'ShopOutlined', 6, 1, 0, NULL, 1, 0, 1000000000000000001, '2026-04-27 00:00:00', '2026-04-27 00:00:00', NULL),
  (1000000000000000091, 1000000000000000090, '商户新增', 3, NULL, NULL, 'system:merchant:create', '', 1, 0, 0, NULL, 1, 0, 1000000000000000001, '2026-04-27 00:00:00', '2026-04-27 00:00:00', NULL),
  (1000000000000000092, 1000000000000000090, '商户修改', 3, NULL, NULL, 'system:merchant:update', '', 2, 0, 0, NULL, 1, 0, 1000000000000000001, '2026-04-27 00:00:00', '2026-04-27 00:00:00', NULL),
  (1000000000000000093, 1000000000000000090, '商户删除', 3, NULL, NULL, 'system:merchant:delete', '', 3, 0, 0, NULL, 1, 0, 1000000000000000001, '2026-04-27 00:00:00', '2026-04-27 00:00:00', NULL),
  (1000000000000000094, 1000000000000000090, '商户批量删除', 3, NULL, NULL, 'system:merchant:batch-delete', '', 4, 0, 0, NULL, 1, 0, 1000000000000000001, '2026-04-27 00:00:00', '2026-04-27 00:00:00', NULL);

INSERT INTO `system_role_menu` (`role_id`, `menu_id`) VALUES
  (1000000000000000002, 1000000000000000080),
  (1000000000000000002, 1000000000000000081),
  (1000000000000000002, 1000000000000000082),
  (1000000000000000002, 1000000000000000083),
  (1000000000000000002, 1000000000000000084),
  (1000000000000000002, 1000000000000000090),
  (1000000000000000002, 1000000000000000091),
  (1000000000000000002, 1000000000000000092),
  (1000000000000000002, 1000000000000000093),
  (1000000000000000002, 1000000000000000094);

SET FOREIGN_KEY_CHECKS = 1;
