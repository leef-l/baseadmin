-- 计费套餐体系

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

CREATE TABLE `system_plan` (
  `id` bigint unsigned NOT NULL COMMENT '套餐ID（Snowflake）',
  `name` varchar(80) NOT NULL COMMENT '套餐名称',
  `code` varchar(50) NOT NULL COMMENT '套餐编码',
  `description` varchar(500) DEFAULT '' COMMENT '套餐描述',
  `user_limit` int NOT NULL DEFAULT '-1' COMMENT '用户数量上限（-1=无限）',
  `storage_mb` int NOT NULL DEFAULT '-1' COMMENT '存储空间上限MB（-1=无限）',
  `features` text COMMENT '功能开关列表（JSON数组）',
  `price` int NOT NULL DEFAULT '0' COMMENT '价格（分）',
  `sort` int NOT NULL DEFAULT '0' COMMENT '排序（升序）',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态:0=关闭,1=开启',
  `created_by` bigint unsigned DEFAULT NULL COMMENT '创建人ID',
  `dept_id` bigint unsigned DEFAULT NULL COMMENT '所属部门ID',
  `tenant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '租户',
  `merchant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '商户',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='套餐表';

CREATE TABLE `system_tenant_plan` (
  `id` bigint unsigned NOT NULL COMMENT '订阅ID（Snowflake）',
  `tenant_id` bigint unsigned NOT NULL COMMENT '租户ID',
  `plan_id` bigint unsigned NOT NULL COMMENT '套餐ID',
  `start_at` datetime DEFAULT NULL COMMENT '生效时间',
  `expire_at` datetime DEFAULT NULL COMMENT '到期时间',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态:0=关闭,1=开启',
  `remark` varchar(500) DEFAULT '' COMMENT '备注',
  `created_by` bigint unsigned DEFAULT NULL COMMENT '创建人ID',
  `dept_id` bigint unsigned DEFAULT NULL COMMENT '所属部门ID',
  `merchant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '商户',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_tenant_id` (`tenant_id`),
  KEY `idx_plan_id` (`plan_id`),
  KEY `idx_expire_at` (`expire_at`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='租户套餐订阅表';

-- 默认免费套餐
INSERT INTO `system_plan` (`id`, `name`, `code`, `description`, `user_limit`, `storage_mb`, `features`, `price`, `sort`, `status`, `tenant_id`, `merchant_id`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
  2000000000000000001, '免费版', 'free', '默认免费套餐',
  10, 100, '["dept","role","user","upload"]', 0, 1, 1,
  0, 0, NOW(), NOW(), NULL
);

INSERT INTO `system_plan` (`id`, `name`, `code`, `description`, `user_limit`, `storage_mb`, `features`, `price`, `sort`, `status`, `tenant_id`, `merchant_id`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
  2000000000000000002, '专业版', 'pro', '适合中型企业',
  100, 10240, '["dept","role","user","upload","domain","merchant","export"]', 29900, 2, 1,
  0, 0, NOW(), NOW(), NULL
);

INSERT INTO `system_plan` (`id`, `name`, `code`, `description`, `user_limit`, `storage_mb`, `features`, `price`, `sort`, `status`, `tenant_id`, `merchant_id`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
  2000000000000000003, '企业版', 'enterprise', '适合大型企业',
  -1, -1, '["dept","role","user","upload","domain","merchant","export","ssl","cname"]', 99900, 3, 1,
  0, 0, NOW(), NOW(), NULL
);

SET FOREIGN_KEY_CHECKS = 1;
