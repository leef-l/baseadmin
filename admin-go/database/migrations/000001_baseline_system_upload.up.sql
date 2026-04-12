-- GBaseAdmin 精简初始化脚本
-- 仅保留当前仓库 `system` / `upload` 相关表和最小种子数据
-- 默认初始化账号: admin / admin123

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ============================================================
-- system: 部门
-- ============================================================
DROP TABLE IF EXISTS `system_dept`;
CREATE TABLE `system_dept` (
  `id` bigint unsigned NOT NULL COMMENT '部门ID（Snowflake）',
  `parent_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '上级部门ID，0 表示顶级部门',
  `title` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '部门名称',
  `username` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '部门负责人姓名',
  `email` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '负责人邮箱',
  `sort` int NOT NULL DEFAULT '0' COMMENT '排序（升序）',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态:0=关闭,1=开启',
  `created_by` bigint unsigned DEFAULT NULL COMMENT '创建人ID',
  `dept_id` bigint unsigned DEFAULT NULL COMMENT '所属部门ID',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间，非 NULL 表示已删除',
  PRIMARY KEY (`id`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='部门表';

INSERT INTO `system_dept` (
  `id`, `parent_id`, `title`, `username`, `email`, `sort`, `status`,
  `created_by`, `dept_id`, `created_at`, `updated_at`, `deleted_at`
) VALUES (
  1000000000000000001, 0, '总公司', 'admin', 'admin@example.com', 0, 1,
  0, 0, '2026-04-07 00:00:00', '2026-04-07 00:00:00', NULL
);

-- ============================================================
-- system: 角色
-- ============================================================
DROP TABLE IF EXISTS `system_role`;
CREATE TABLE `system_role` (
  `id` bigint unsigned NOT NULL COMMENT '角色ID（Snowflake）',
  `parent_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '上级角色ID，0 表示顶级角色',
  `title` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '角色名称',
  `data_scope` tinyint NOT NULL DEFAULT '1' COMMENT '数据范围:1=全部,2=本部门及以下,3=本部门,4=仅本人,5=自定义',
  `is_admin` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否超级管理员:0=否,1=是',
  `sort` int NOT NULL DEFAULT '0' COMMENT '排序（升序）',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态:0=关闭,1=开启',
  `created_by` bigint unsigned DEFAULT NULL COMMENT '创建人ID',
  `dept_id` bigint unsigned DEFAULT NULL COMMENT '所属部门ID',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间，非 NULL 表示已删除',
  PRIMARY KEY (`id`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色表';

INSERT INTO `system_role` (
  `id`, `parent_id`, `title`, `data_scope`, `is_admin`, `sort`, `status`,
  `created_by`, `dept_id`, `created_at`, `updated_at`, `deleted_at`
) VALUES (
  1000000000000000002, 0, '超级管理员', 1, 1, 0, 1,
  0, 1000000000000000001, '2026-04-07 00:00:00', '2026-04-07 00:00:00', NULL
);

-- ============================================================
-- system: 用户
-- ============================================================
DROP TABLE IF EXISTS `system_users`;
CREATE TABLE `system_users` (
  `id` bigint unsigned NOT NULL COMMENT '用户ID（Snowflake）',
  `username` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '登录用户名',
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '密码（bcrypt 加密）',
  `nickname` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '昵称/显示名',
  `email` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '邮箱地址',
  `avatar` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '头像图片 URL',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态:0=关闭,1=开启',
  `created_by` bigint unsigned DEFAULT NULL COMMENT '创建人ID',
  `dept_id` bigint unsigned DEFAULT NULL COMMENT '所属部门ID',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间，非 NULL 表示已删除',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`),
  KEY `idx_dept_id` (`dept_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

INSERT INTO `system_users` (
  `id`, `username`, `password`, `nickname`, `email`, `avatar`, `status`,
  `created_by`, `dept_id`, `created_at`, `updated_at`, `deleted_at`
) VALUES (
  1000000000000000003,
  'admin',
  '$2y$10$VWVHbVcWz5hfZlNf6Q.99ONLkvL5ktwKz9F70PM7dg9Nq5LXVT4gS',
  '超级管理员',
  'admin@example.com',
  '',
  1,
  0,
  1000000000000000001,
  '2026-04-07 00:00:00',
  '2026-04-07 00:00:00',
  NULL
);

-- ============================================================
-- system: 关联表
-- ============================================================
DROP TABLE IF EXISTS `system_user_role`;
CREATE TABLE `system_user_role` (
  `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
  `role_id` bigint unsigned NOT NULL COMMENT '角色ID',
  PRIMARY KEY (`user_id`,`role_id`),
  KEY `idx_role_id` (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户角色关联表';

INSERT INTO `system_user_role` (`user_id`, `role_id`) VALUES
  (1000000000000000003, 1000000000000000002);

DROP TABLE IF EXISTS `system_user_dept`;
CREATE TABLE `system_user_dept` (
  `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
  `dept_id` bigint unsigned NOT NULL COMMENT '部门ID',
  PRIMARY KEY (`user_id`,`dept_id`),
  KEY `idx_dept_id` (`dept_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户部门关联表';

INSERT INTO `system_user_dept` (`user_id`, `dept_id`) VALUES
  (1000000000000000003, 1000000000000000001);

DROP TABLE IF EXISTS `system_role_dept`;
CREATE TABLE `system_role_dept` (
  `role_id` bigint unsigned NOT NULL COMMENT '角色ID',
  `dept_id` bigint unsigned NOT NULL COMMENT '部门ID',
  PRIMARY KEY (`role_id`,`dept_id`),
  KEY `idx_dept_id` (`dept_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色自定义数据权限部门关联表';

-- ============================================================
-- system + upload: 菜单
-- ============================================================
DROP TABLE IF EXISTS `system_menu`;
CREATE TABLE `system_menu` (
  `id` bigint unsigned NOT NULL COMMENT '菜单ID（Snowflake）',
  `parent_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '上级菜单ID，0 表示顶级菜单',
  `title` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '菜单名称',
  `type` tinyint NOT NULL DEFAULT '1' COMMENT '类型:1=目录,2=菜单,3=按钮,4=外链,5=内链',
  `path` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '前端路由路径',
  `component` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '前端组件路径',
  `permission` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '权限标识（如 system:dept:list）',
  `icon` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '菜单图标（图标名称）',
  `sort` int NOT NULL DEFAULT '0' COMMENT '排序（升序）',
  `is_show` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否显示:0=隐藏,1=显示',
  `is_cache` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否缓存:0=不缓存,1=缓存',
  `link_url` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '外链/内链地址（type=4或5时有效）',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态:0=关闭,1=开启',
  `created_by` bigint unsigned DEFAULT NULL COMMENT '创建人ID',
  `dept_id` bigint unsigned DEFAULT NULL COMMENT '所属部门ID',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间，非 NULL 表示已删除',
  PRIMARY KEY (`id`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='菜单表';

INSERT INTO `system_menu` (
  `id`, `parent_id`, `title`, `type`, `path`, `component`, `permission`, `icon`,
  `sort`, `is_show`, `is_cache`, `link_url`, `status`, `created_by`, `dept_id`,
  `created_at`, `updated_at`, `deleted_at`
) VALUES
  (1000000000000000010, 0, '系统管理', 1, '/system', NULL, '', 'SettingOutlined', 10, 1, 0, NULL, 1, 0, 1000000000000000001, '2026-04-07 00:00:00', '2026-04-07 00:00:00', NULL),
  (1000000000000000011, 1000000000000000010, '部门管理', 2, '/system/dept', 'system/dept/index', 'system:dept:list', 'ApartmentOutlined', 1, 1, 0, NULL, 1, 0, 1000000000000000001, '2026-04-07 00:00:00', '2026-04-07 00:00:00', NULL),
  (1000000000000000021, 1000000000000000011, '部门新增', 3, NULL, NULL, 'system:dept:create', '', 1, 0, 0, NULL, 1, 0, 1000000000000000001, '2026-04-07 00:00:00', '2026-04-07 00:00:00', NULL),
  (1000000000000000022, 1000000000000000011, '部门修改', 3, NULL, NULL, 'system:dept:update', '', 2, 0, 0, NULL, 1, 0, 1000000000000000001, '2026-04-07 00:00:00', '2026-04-07 00:00:00', NULL),
  (1000000000000000023, 1000000000000000011, '部门删除', 3, NULL, NULL, 'system:dept:delete', '', 3, 0, 0, NULL, 1, 0, 1000000000000000001, '2026-04-07 00:00:00', '2026-04-07 00:00:00', NULL),
  (1000000000000000024, 1000000000000000011, '部门批量删除', 3, NULL, NULL, 'system:dept:batch-delete', '', 4, 0, 0, NULL, 1, 0, 1000000000000000001, '2026-04-07 00:00:00', '2026-04-07 00:00:00', NULL),
  (1000000000000000012, 1000000000000000010, '角色管理', 2, '/system/role', 'system/role/index', 'system:role:list', 'TeamOutlined', 2, 1, 0, NULL, 1, 0, 1000000000000000001, '2026-04-07 00:00:00', '2026-04-07 00:00:00', NULL),
  (1000000000000000031, 1000000000000000012, '角色新增', 3, NULL, NULL, 'system:role:create', '', 1, 0, 0, NULL, 1, 0, 1000000000000000001, '2026-04-07 00:00:00', '2026-04-07 00:00:00', NULL),
  (1000000000000000032, 1000000000000000012, '角色修改', 3, NULL, NULL, 'system:role:update', '', 2, 0, 0, NULL, 1, 0, 1000000000000000001, '2026-04-07 00:00:00', '2026-04-07 00:00:00', NULL),
  (1000000000000000033, 1000000000000000012, '角色删除', 3, NULL, NULL, 'system:role:delete', '', 3, 0, 0, NULL, 1, 0, 1000000000000000001, '2026-04-07 00:00:00', '2026-04-07 00:00:00', NULL),
  (1000000000000000034, 1000000000000000012, '资源授权', 3, NULL, NULL, 'system:role:grant:menu', '', 4, 0, 0, NULL, 1, 0, 1000000000000000001, '2026-04-07 00:00:00', '2026-04-07 00:00:00', NULL),
  (1000000000000000035, 1000000000000000012, '数据授权', 3, NULL, NULL, 'system:role:grant:dept', '', 5, 0, 0, NULL, 1, 0, 1000000000000000001, '2026-04-07 00:00:00', '2026-04-07 00:00:00', NULL),
  (1000000000000000036, 1000000000000000012, '角色批量删除', 3, NULL, NULL, 'system:role:batch-delete', '', 6, 0, 0, NULL, 1, 0, 1000000000000000001, '2026-04-07 00:00:00', '2026-04-07 00:00:00', NULL),
  (1000000000000000013, 1000000000000000010, '菜单管理', 2, '/system/menu', 'system/menu/index', 'system:menu:list', 'MenuOutlined', 3, 1, 0, NULL, 1, 0, 1000000000000000001, '2026-04-07 00:00:00', '2026-04-07 00:00:00', NULL),
  (1000000000000000041, 1000000000000000013, '菜单新增', 3, NULL, NULL, 'system:menu:create', '', 1, 0, 0, NULL, 1, 0, 1000000000000000001, '2026-04-07 00:00:00', '2026-04-07 00:00:00', NULL),
  (1000000000000000042, 1000000000000000013, '菜单修改', 3, NULL, NULL, 'system:menu:update', '', 2, 0, 0, NULL, 1, 0, 1000000000000000001, '2026-04-07 00:00:00', '2026-04-07 00:00:00', NULL),
  (1000000000000000043, 1000000000000000013, '菜单删除', 3, NULL, NULL, 'system:menu:delete', '', 3, 0, 0, NULL, 1, 0, 1000000000000000001, '2026-04-07 00:00:00', '2026-04-07 00:00:00', NULL),
  (1000000000000000044, 1000000000000000013, '菜单批量删除', 3, NULL, NULL, 'system:menu:batch-delete', '', 4, 0, 0, NULL, 1, 0, 1000000000000000001, '2026-04-07 00:00:00', '2026-04-07 00:00:00', NULL),
  (1000000000000000014, 1000000000000000010, '用户管理', 2, '/system/users', 'system/users/index', 'system:user:list', 'UserOutlined', 4, 1, 0, NULL, 1, 0, 1000000000000000001, '2026-04-07 00:00:00', '2026-04-07 00:00:00', NULL),
  (1000000000000000051, 1000000000000000014, '用户新增', 3, NULL, NULL, 'system:user:create', '', 1, 0, 0, NULL, 1, 0, 1000000000000000001, '2026-04-07 00:00:00', '2026-04-07 00:00:00', NULL),
  (1000000000000000052, 1000000000000000014, '用户修改', 3, NULL, NULL, 'system:user:update', '', 2, 0, 0, NULL, 1, 0, 1000000000000000001, '2026-04-07 00:00:00', '2026-04-07 00:00:00', NULL),
  (1000000000000000053, 1000000000000000014, '用户删除', 3, NULL, NULL, 'system:user:delete', '', 3, 0, 0, NULL, 1, 0, 1000000000000000001, '2026-04-07 00:00:00', '2026-04-07 00:00:00', NULL),
  (1000000000000000054, 1000000000000000014, '用户批量删除', 3, NULL, NULL, 'system:user:batch-delete', '', 4, 0, 0, NULL, 1, 0, 1000000000000000001, '2026-04-07 00:00:00', '2026-04-07 00:00:00', NULL),
  (314253730209861632, 0, '上传管理', 1, '/upload', NULL, '', 'CloudUploadOutlined', 20, 1, 0, NULL, 1, 0, 0, '2026-04-07 00:00:00', '2026-04-07 00:00:00', NULL),
  (314253730235027456, 314253730209861632, '上传配置', 2, '/upload/config', 'upload/config/index', 'upload:config:list', 'ControlOutlined', 1, 1, 0, NULL, 1, 0, 0, '2026-04-07 00:00:00', '2026-04-07 00:00:00', NULL),
  (314253730251804672, 314253730235027456, '上传配置新增', 3, NULL, NULL, 'upload:config:create', '', 1, 0, 0, NULL, 1, 0, 0, '2026-04-07 00:00:00', '2026-04-07 00:00:00', NULL),
  (314253730268581888, 314253730235027456, '上传配置修改', 3, NULL, NULL, 'upload:config:update', '', 2, 0, 0, NULL, 1, 0, 0, '2026-04-07 00:00:00', '2026-04-07 00:00:00', NULL),
  (314253730285359104, 314253730235027456, '上传配置删除', 3, NULL, NULL, 'upload:config:delete', '', 3, 0, 0, NULL, 1, 0, 0, '2026-04-07 00:00:00', '2026-04-07 00:00:00', NULL),
  (314253730302136320, 314253730235027456, '上传配置批量删除', 3, NULL, NULL, 'upload:config:batch-delete', '', 4, 0, 0, NULL, 1, 0, 0, '2026-04-07 00:00:00', '2026-04-07 00:00:00', NULL),
  (314253730344079360, 314253730209861632, '文件目录', 2, '/upload/dir', 'upload/dir/index', 'upload:dir:list', 'FolderOpenOutlined', 2, 1, 0, NULL, 1, 0, 0, '2026-04-07 00:00:00', '2026-04-07 00:00:00', NULL),
  (314253730365050880, 314253730344079360, '文件目录新增', 3, NULL, NULL, 'upload:dir:create', '', 1, 0, 0, NULL, 1, 0, 0, '2026-04-07 00:00:00', '2026-04-07 00:00:00', NULL),
  (314253730386022400, 314253730344079360, '文件目录修改', 3, NULL, NULL, 'upload:dir:update', '', 2, 0, 0, NULL, 1, 0, 0, '2026-04-07 00:00:00', '2026-04-07 00:00:00', NULL),
  (314253730415382528, 314253730344079360, '文件目录删除', 3, NULL, NULL, 'upload:dir:delete', '', 3, 0, 0, NULL, 1, 0, 0, '2026-04-07 00:00:00', '2026-04-07 00:00:00', NULL),
  (314253730436354048, 314253730344079360, '文件目录批量删除', 3, NULL, NULL, 'upload:dir:batch-delete', '', 4, 0, 0, NULL, 1, 0, 0, '2026-04-07 00:00:00', '2026-04-07 00:00:00', NULL),
  (314253730461519872, 314253730209861632, '文件目录规则', 2, '/upload/dir-rule', 'upload/dir_rule/index', 'upload:dir_rule:list', 'PartitionOutlined', 3, 1, 0, NULL, 1, 0, 0, '2026-04-07 00:00:00', '2026-04-07 00:00:00', NULL),
  (314253730478297088, 314253730461519872, '文件目录规则新增', 3, NULL, NULL, 'upload:dir_rule:create', '', 1, 0, 0, NULL, 1, 0, 0, '2026-04-07 00:00:00', '2026-04-07 00:00:00', NULL),
  (314253730490880000, 314253730461519872, '文件目录规则修改', 3, NULL, NULL, 'upload:dir_rule:update', '', 2, 0, 0, NULL, 1, 0, 0, '2026-04-07 00:00:00', '2026-04-07 00:00:00', NULL),
  (314253730503462912, 314253730461519872, '文件目录规则删除', 3, NULL, NULL, 'upload:dir_rule:delete', '', 3, 0, 0, NULL, 1, 0, 0, '2026-04-07 00:00:00', '2026-04-07 00:00:00', NULL),
  (314253730520240128, 314253730461519872, '文件目录规则批量删除', 3, NULL, NULL, 'upload:dir_rule:batch-delete', '', 4, 0, 0, NULL, 1, 0, 0, '2026-04-07 00:00:00', '2026-04-07 00:00:00', NULL),
  (314253730566377472, 314253730209861632, '文件记录', 2, '/upload/file', 'upload/file/index', 'upload:file:list', 'FileTextOutlined', 4, 1, 0, NULL, 1, 0, 0, '2026-04-07 00:00:00', '2026-04-07 00:00:00', NULL),
  (314253730583154688, 314253730566377472, '文件记录新增', 3, NULL, NULL, 'upload:file:create', '', 1, 0, 0, NULL, 1, 0, 0, '2026-04-07 00:00:00', '2026-04-07 00:00:00', NULL),
  (314253730595737600, 314253730566377472, '文件记录修改', 3, NULL, NULL, 'upload:file:update', '', 2, 0, 0, NULL, 1, 0, 0, '2026-04-07 00:00:00', '2026-04-07 00:00:00', NULL),
  (314253730620903424, 314253730566377472, '文件记录删除', 3, NULL, NULL, 'upload:file:delete', '', 3, 0, 0, NULL, 1, 0, 0, '2026-04-07 00:00:00', '2026-04-07 00:00:00', NULL),
  (314253730637680640, 314253730566377472, '文件记录批量删除', 3, NULL, NULL, 'upload:file:batch-delete', '', 4, 0, 0, NULL, 1, 0, 0, '2026-04-07 00:00:00', '2026-04-07 00:00:00', NULL);

DROP TABLE IF EXISTS `system_role_menu`;
CREATE TABLE `system_role_menu` (
  `role_id` bigint unsigned NOT NULL COMMENT '角色ID',
  `menu_id` bigint unsigned NOT NULL COMMENT '菜单ID',
  PRIMARY KEY (`role_id`,`menu_id`),
  KEY `idx_menu_id` (`menu_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色菜单权限关联表';

INSERT INTO `system_role_menu` (`role_id`, `menu_id`) VALUES
  (1000000000000000002, 1000000000000000010),
  (1000000000000000002, 1000000000000000011),
  (1000000000000000002, 1000000000000000021),
  (1000000000000000002, 1000000000000000022),
  (1000000000000000002, 1000000000000000023),
  (1000000000000000002, 1000000000000000024),
  (1000000000000000002, 1000000000000000012),
  (1000000000000000002, 1000000000000000031),
  (1000000000000000002, 1000000000000000032),
  (1000000000000000002, 1000000000000000033),
  (1000000000000000002, 1000000000000000034),
  (1000000000000000002, 1000000000000000035),
  (1000000000000000002, 1000000000000000036),
  (1000000000000000002, 1000000000000000013),
  (1000000000000000002, 1000000000000000041),
  (1000000000000000002, 1000000000000000042),
  (1000000000000000002, 1000000000000000043),
  (1000000000000000002, 1000000000000000044),
  (1000000000000000002, 1000000000000000014),
  (1000000000000000002, 1000000000000000051),
  (1000000000000000002, 1000000000000000052),
  (1000000000000000002, 1000000000000000053),
  (1000000000000000002, 1000000000000000054),
  (1000000000000000002, 314253730209861632),
  (1000000000000000002, 314253730235027456),
  (1000000000000000002, 314253730251804672),
  (1000000000000000002, 314253730268581888),
  (1000000000000000002, 314253730285359104),
  (1000000000000000002, 314253730302136320),
  (1000000000000000002, 314253730344079360),
  (1000000000000000002, 314253730365050880),
  (1000000000000000002, 314253730386022400),
  (1000000000000000002, 314253730415382528),
  (1000000000000000002, 314253730436354048),
  (1000000000000000002, 314253730461519872),
  (1000000000000000002, 314253730478297088),
  (1000000000000000002, 314253730490880000),
  (1000000000000000002, 314253730503462912),
  (1000000000000000002, 314253730520240128),
  (1000000000000000002, 314253730566377472),
  (1000000000000000002, 314253730583154688),
  (1000000000000000002, 314253730595737600),
  (1000000000000000002, 314253730620903424),
  (1000000000000000002, 314253730637680640);

-- ============================================================
-- upload: 上传配置
-- ============================================================
DROP TABLE IF EXISTS `upload_config`;
CREATE TABLE `upload_config` (
  `id` bigint unsigned NOT NULL COMMENT 'ID',
  `name` varchar(100) NOT NULL COMMENT '配置名称',
  `storage` tinyint(1) NOT NULL DEFAULT '1' COMMENT '存储类型:1=本地,2=阿里云OSS,3=腾讯云COS',
  `is_default` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否默认:0=否,1=是',
  `local_path` varchar(500) DEFAULT '' COMMENT '本地存储路径',
  `oss_endpoint` varchar(255) DEFAULT '' COMMENT 'OSS Endpoint',
  `oss_bucket` varchar(255) DEFAULT '' COMMENT 'OSS Bucket',
  `oss_access_key` varchar(255) DEFAULT '' COMMENT 'OSS AccessKey',
  `oss_secret_key` varchar(255) DEFAULT '' COMMENT 'OSS SecretKey',
  `cos_region` varchar(100) DEFAULT '' COMMENT 'COS Region',
  `cos_bucket` varchar(255) DEFAULT '' COMMENT 'COS Bucket',
  `cos_secret_id` varchar(255) DEFAULT '' COMMENT 'COS SecretId',
  `cos_secret_key` varchar(255) DEFAULT '' COMMENT 'COS SecretKey',
  `max_size` int DEFAULT '10' COMMENT '最大文件大小(MB)',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态:0=禁用,1=启用',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  `created_by` bigint unsigned DEFAULT NULL COMMENT '创建人',
  `dept_id` bigint unsigned DEFAULT NULL COMMENT '部门ID',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='上传配置';

INSERT INTO `upload_config` (
  `id`, `name`, `storage`, `is_default`, `local_path`,
  `oss_endpoint`, `oss_bucket`, `oss_access_key`, `oss_secret_key`,
  `cos_region`, `cos_bucket`, `cos_secret_id`, `cos_secret_key`,
  `max_size`, `status`, `created_at`, `updated_at`, `deleted_at`, `created_by`, `dept_id`
) VALUES (
  314309590294466560, '本地存储', 1, 1, 'resource/upload',
  '', '', '', '',
  '', '', '', '',
  10, 1, '2026-04-07 00:00:00', '2026-04-07 00:00:00', NULL, NULL, NULL
);

-- ============================================================
-- upload: 文件目录
-- ============================================================
DROP TABLE IF EXISTS `upload_dir`;
CREATE TABLE `upload_dir` (
  `id` bigint unsigned NOT NULL COMMENT 'ID',
  `parent_id` bigint unsigned DEFAULT '0' COMMENT '上级目录',
  `name` varchar(100) NOT NULL COMMENT '目录名称',
  `path` varchar(500) NOT NULL COMMENT '目录路径',
  `sort` int DEFAULT '0' COMMENT '排序',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态:0=禁用,1=启用',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  `created_by` bigint unsigned DEFAULT NULL COMMENT '创建人',
  `dept_id` bigint unsigned DEFAULT NULL COMMENT '部门ID',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文件目录';

INSERT INTO `upload_dir` (
  `id`, `parent_id`, `name`, `path`, `sort`, `status`,
  `created_at`, `updated_at`, `deleted_at`, `created_by`, `dept_id`
) VALUES (
  314696302266945536, 0, 'uploads', 'uploads', 0, 1,
  '2026-04-07 00:00:00', '2026-04-07 00:00:00', NULL, NULL, NULL
);

-- ============================================================
-- upload: 文件目录规则
-- ============================================================
DROP TABLE IF EXISTS `upload_dir_rule`;
CREATE TABLE `upload_dir_rule` (
  `id` bigint unsigned NOT NULL COMMENT 'ID',
  `dir_id` bigint unsigned NOT NULL COMMENT '目录ID',
  `category` tinyint(1) NOT NULL DEFAULT '1' COMMENT '类别:1=默认,2=类型,3=接口',
  `file_type` varchar(255) DEFAULT '' COMMENT '文件类型，多个用逗号分隔',
  `save_path` varchar(500) DEFAULT '' COMMENT '保存目录',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态:0=禁用,1=启用',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  `created_by` bigint unsigned DEFAULT NULL COMMENT '创建人',
  `dept_id` bigint unsigned DEFAULT NULL COMMENT '部门ID',
  PRIMARY KEY (`id`),
  KEY `idx_dir_id` (`dir_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文件目录规则';

INSERT INTO `upload_dir_rule` (
  `id`, `dir_id`, `category`, `file_type`, `save_path`, `status`,
  `created_at`, `updated_at`, `deleted_at`, `created_by`, `dept_id`
) VALUES (
  314720336681635840, 314696302266945536, 1, '', '{Y-m-d}', 1,
  '2026-04-07 00:00:00', '2026-04-07 00:00:00', NULL, NULL, NULL
);

-- ============================================================
-- upload: 文件记录
-- ============================================================
DROP TABLE IF EXISTS `upload_file`;
CREATE TABLE `upload_file` (
  `id` bigint unsigned NOT NULL COMMENT 'ID',
  `dir_id` bigint unsigned DEFAULT '0' COMMENT '所属目录',
  `name` varchar(255) NOT NULL COMMENT '文件名称',
  `url` varchar(500) NOT NULL COMMENT '文件地址',
  `ext` varchar(20) DEFAULT '' COMMENT '文件扩展名',
  `size` bigint unsigned DEFAULT '0' COMMENT '文件大小',
  `mime` varchar(100) DEFAULT '' COMMENT 'MIME类型',
  `storage` tinyint(1) NOT NULL DEFAULT '1' COMMENT '存储类型:1=本地,2=阿里云OSS,3=腾讯云COS',
  `is_image` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否图片:0=否,1=是',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  `created_by` bigint unsigned DEFAULT NULL COMMENT '创建人',
  `dept_id` bigint unsigned DEFAULT NULL COMMENT '部门ID',
  PRIMARY KEY (`id`),
  KEY `idx_dir_id` (`dir_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文件记录';

SET FOREIGN_KEY_CHECKS = 1;
