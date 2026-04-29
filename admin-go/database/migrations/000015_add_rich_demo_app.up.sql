-- 丰富 Demo 体验应用
-- 覆盖 codegen 常用控件、搜索、树形、外键、导入导出、批量编辑、租户/商户/部门数据权限字段。

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- Demo 体验数据归属主体。只写入各环境共有字段，避免和历史线上库的 SaaS 迁移顺序耦合。
INSERT INTO `system_tenant` (`id`, `name`, `code`, `contact_name`, `contact_phone`, `domain`, `expire_at`, `status`, `remark`, `created_by`, `dept_id`, `tenant_id`, `merchant_id`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (324500000000000001, 'Demo体验租户', 'demo-tenant', 'Demo负责人', '13800000001', 'demo.easytestdev.online', '2036-12-31 23:59:59', 1, '用于 codegen 和后台组件体验的演示租户', 1000000000000000003, 1000000000000000001, 324500000000000001, 0, NOW(), NOW(), NULL)
ON DUPLICATE KEY UPDATE
  `name` = VALUES(`name`),
  `contact_name` = VALUES(`contact_name`),
  `contact_phone` = VALUES(`contact_phone`),
  `domain` = VALUES(`domain`),
  `expire_at` = VALUES(`expire_at`),
  `status` = VALUES(`status`),
  `remark` = VALUES(`remark`),
  `tenant_id` = VALUES(`tenant_id`),
  `merchant_id` = VALUES(`merchant_id`),
  `updated_at` = NOW(),
  `deleted_at` = NULL;

INSERT INTO `system_merchant` (`id`, `tenant_id`, `merchant_id`, `name`, `code`, `contact_name`, `contact_phone`, `address`, `status`, `remark`, `created_by`, `dept_id`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (324500000000000101, 324500000000000001, 324500000000000101, 'Demo体验商户', 'demo-merchant', 'Demo店长', '13800000002', '上海市浦东新区 Demo 路 88 号', 1, '用于 codegen 和后台组件体验的演示商户', 1000000000000000003, 1000000000000000001, NOW(), NOW(), NULL)
ON DUPLICATE KEY UPDATE
  `tenant_id` = VALUES(`tenant_id`),
  `merchant_id` = VALUES(`merchant_id`),
  `name` = VALUES(`name`),
  `contact_name` = VALUES(`contact_name`),
  `contact_phone` = VALUES(`contact_phone`),
  `address` = VALUES(`address`),
  `status` = VALUES(`status`),
  `remark` = VALUES(`remark`),
  `updated_at` = NOW(),
  `deleted_at` = NULL;

CREATE TABLE IF NOT EXISTS `demo_category` (
  `id` bigint unsigned NOT NULL COMMENT '分类ID（Snowflake）',
  `parent_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '父分类',
  `name` varchar(80) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '分类名称|search:like|keyword:on|priority:95',
  `icon` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '图标',
  `sort` int NOT NULL DEFAULT '0' COMMENT '排序（升序）',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态:0=禁用,1=启用',
  `tenant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '租户',
  `merchant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '商户',
  `created_by` bigint unsigned NOT NULL DEFAULT '0' COMMENT '创建人ID',
  `dept_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '所属部门ID',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间，非 NULL 表示已删除',
  PRIMARY KEY (`id`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_tenant_merchant` (`tenant_id`, `merchant_id`),
  KEY `idx_status` (`status`),
  KEY `idx_sort` (`sort`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='体验分类';

CREATE TABLE IF NOT EXISTS `demo_customer` (
  `id` bigint unsigned NOT NULL COMMENT '客户ID（Snowflake）',
  `avatar` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '头像',
  `name` varchar(80) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '客户名称|search:like|keyword:on|priority:95',
  `customer_no` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '客户编号|search:eq|priority:100',
  `phone` varchar(30) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '联系电话|search:like|keyword:on|priority:90',
  `email` varchar(120) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '邮箱|search:like|keyword:on|priority:90',
  `gender` tinyint NOT NULL DEFAULT '0' COMMENT '性别:0=未知,1=男,2=女',
  `level` tinyint NOT NULL DEFAULT '1' COMMENT '等级:1=普通,2=VIP,3=付费,4=冻结',
  `source_type` tinyint NOT NULL DEFAULT '1' COMMENT '来源:1=官网,2=小程序,3=线下,4=导入',
  `is_vip` tinyint NOT NULL DEFAULT '0' COMMENT '是否VIP:0=否,1=是',
  `registered_at` datetime DEFAULT NULL COMMENT '注册时间',
  `remark` text COLLATE utf8mb4_unicode_ci COMMENT '备注|search:like|keyword:only',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态:0=禁用,1=启用',
  `tenant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '租户',
  `merchant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '商户',
  `created_by` bigint unsigned NOT NULL DEFAULT '0' COMMENT '创建人ID',
  `dept_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '所属部门ID',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间，非 NULL 表示已删除',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_customer_no` (`customer_no`),
  KEY `idx_tenant_merchant` (`tenant_id`, `merchant_id`),
  KEY `idx_phone` (`phone`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='体验客户';

CREATE TABLE IF NOT EXISTS `demo_product` (
  `id` bigint unsigned NOT NULL COMMENT '商品ID（Snowflake）',
  `category_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '商品分类',
  `sku_no` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'SKU编号|search:eq|priority:100',
  `name` varchar(120) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '商品名称|search:like|keyword:on|priority:95',
  `cover` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '封面',
  `manual_file` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '说明书文件',
  `detail_content` text COLLATE utf8mb4_unicode_ci COMMENT '详情内容',
  `spec_json` text COLLATE utf8mb4_unicode_ci COMMENT '规格JSON',
  `website_url` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '官网URL',
  `type` tinyint NOT NULL DEFAULT '1' COMMENT '类型:1=普通,2=置顶,3=推荐,4=热门',
  `is_recommend` tinyint NOT NULL DEFAULT '0' COMMENT '是否推荐:0=否,1=是',
  `sale_price` int NOT NULL DEFAULT '0' COMMENT '销售价（分）',
  `stock_num` int NOT NULL DEFAULT '0' COMMENT '库存数量',
  `weight_num` int NOT NULL DEFAULT '0' COMMENT '重量（克）',
  `sort` int NOT NULL DEFAULT '0' COMMENT '排序（升序）',
  `icon` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '图标',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态:0=草稿,1=上架,2=下架',
  `tenant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '租户',
  `merchant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '商户',
  `created_by` bigint unsigned NOT NULL DEFAULT '0' COMMENT '创建人ID',
  `dept_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '所属部门ID',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间，非 NULL 表示已删除',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_sku_no` (`sku_no`),
  KEY `idx_category_id` (`category_id`),
  KEY `idx_tenant_merchant` (`tenant_id`, `merchant_id`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='体验商品';

CREATE TABLE IF NOT EXISTS `demo_campaign` (
  `id` bigint unsigned NOT NULL COMMENT '活动ID（Snowflake）',
  `campaign_no` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '活动编号|search:eq|priority:100',
  `title` varchar(120) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '活动标题|search:like|keyword:on|priority:95',
  `banner` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '横幅图',
  `type` tinyint NOT NULL DEFAULT '1' COMMENT '活动类型:1=免费,2=付费,3=公开,4=私密',
  `channel` tinyint NOT NULL DEFAULT '1' COMMENT '投放渠道:1=官网,2=小程序,3=短信,4=线下',
  `budget_amount` int NOT NULL DEFAULT '0' COMMENT '预算金额（分）',
  `landing_url` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '落地页URL',
  `rule_json` text COLLATE utf8mb4_unicode_ci COMMENT '规则JSON',
  `intro_content` text COLLATE utf8mb4_unicode_ci COMMENT '活动介绍',
  `start_at` datetime DEFAULT NULL COMMENT '开始时间',
  `end_at` datetime DEFAULT NULL COMMENT '结束时间',
  `is_public` tinyint NOT NULL DEFAULT '1' COMMENT '是否公开:0=否,1=是',
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '状态:0=草稿,1=已发布,2=已下架',
  `tenant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '租户',
  `merchant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '商户',
  `created_by` bigint unsigned NOT NULL DEFAULT '0' COMMENT '创建人ID',
  `dept_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '所属部门ID',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间，非 NULL 表示已删除',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_campaign_no` (`campaign_no`),
  KEY `idx_tenant_merchant` (`tenant_id`, `merchant_id`),
  KEY `idx_status` (`status`),
  KEY `idx_time` (`start_at`, `end_at`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='体验活动';

CREATE TABLE IF NOT EXISTS `demo_order` (
  `id` bigint unsigned NOT NULL COMMENT '订单ID（Snowflake）',
  `order_no` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '订单号|search:eq|priority:100',
  `customer_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '客户',
  `product_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '商品',
  `quantity` int NOT NULL DEFAULT '1' COMMENT '购买数量',
  `amount` int NOT NULL DEFAULT '0' COMMENT '订单金额（分）',
  `pay_status` tinyint NOT NULL DEFAULT '0' COMMENT '支付状态:0=待支付,1=已支付,2=已退款',
  `deliver_status` tinyint NOT NULL DEFAULT '0' COMMENT '发货状态:0=待发货,1=已发货,2=已签收',
  `paid_at` datetime DEFAULT NULL COMMENT '支付时间',
  `deliver_at` datetime DEFAULT NULL COMMENT '发货时间',
  `receiver_phone` varchar(30) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '收货电话',
  `address` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '收货地址|keyword:only',
  `remark` text COLLATE utf8mb4_unicode_ci COMMENT '备注|keyword:only',
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '状态:0=待确认,1=已确认,2=已取消',
  `tenant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '租户',
  `merchant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '商户',
  `created_by` bigint unsigned NOT NULL DEFAULT '0' COMMENT '创建人ID',
  `dept_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '所属部门ID',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间，非 NULL 表示已删除',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_order_no` (`order_no`),
  KEY `idx_customer_id` (`customer_id`),
  KEY `idx_product_id` (`product_id`),
  KEY `idx_tenant_merchant` (`tenant_id`, `merchant_id`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='体验订单';

CREATE TABLE IF NOT EXISTS `demo_work_order` (
  `id` bigint unsigned NOT NULL COMMENT '工单ID（Snowflake）',
  `ticket_no` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '工单号|search:eq|priority:100',
  `customer_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '客户',
  `product_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '商品',
  `order_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '订单',
  `title` varchar(120) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '工单标题|search:like|keyword:on|priority:95',
  `priority` tinyint NOT NULL DEFAULT '2' COMMENT '优先级:1=低,2=普通,3=高,4=紧急',
  `source_type` tinyint NOT NULL DEFAULT '1' COMMENT '来源:1=官网,2=电话,3=微信,4=后台',
  `description` text COLLATE utf8mb4_unicode_ci COMMENT '问题描述|search:like|keyword:only',
  `attachment_file` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '附件',
  `due_at` datetime DEFAULT NULL COMMENT '截止时间',
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '状态:0=待处理,1=进行中,2=已完成,3=已取消',
  `tenant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '租户',
  `merchant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '商户',
  `created_by` bigint unsigned NOT NULL DEFAULT '0' COMMENT '创建人ID',
  `dept_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '所属部门ID',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间，非 NULL 表示已删除',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_ticket_no` (`ticket_no`),
  KEY `idx_customer_id` (`customer_id`),
  KEY `idx_product_id` (`product_id`),
  KEY `idx_order_id` (`order_id`),
  KEY `idx_tenant_merchant` (`tenant_id`, `merchant_id`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='体验工单';

CREATE TABLE IF NOT EXISTS `demo_contract` (
  `id` bigint unsigned NOT NULL COMMENT '合同ID（Snowflake）',
  `contract_no` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '合同编号|search:eq|priority:100',
  `customer_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '客户',
  `order_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '订单',
  `title` varchar(120) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '合同标题|search:like|keyword:on|priority:95',
  `contract_file` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '合同文件',
  `sign_image` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '签章图片',
  `contract_amount` int NOT NULL DEFAULT '0' COMMENT '合同金额（分）',
  `sign_password` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '签署密码',
  `signed_at` datetime DEFAULT NULL COMMENT '签署时间',
  `expires_at` datetime DEFAULT NULL COMMENT '到期时间',
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '状态:0=待审核,1=已通过,2=已拒绝,3=已取消',
  `tenant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '租户',
  `merchant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '商户',
  `created_by` bigint unsigned NOT NULL DEFAULT '0' COMMENT '创建人ID',
  `dept_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '所属部门ID',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间，非 NULL 表示已删除',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_contract_no` (`contract_no`),
  KEY `idx_customer_id` (`customer_id`),
  KEY `idx_order_id` (`order_id`),
  KEY `idx_tenant_merchant` (`tenant_id`, `merchant_id`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='体验合同';

CREATE TABLE IF NOT EXISTS `demo_survey` (
  `id` bigint unsigned NOT NULL COMMENT '问卷ID（Snowflake）',
  `survey_no` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '问卷编号|search:eq|priority:100',
  `title` varchar(120) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '问卷标题|search:like|keyword:on|priority:95',
  `poster` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '海报',
  `question_json` text COLLATE utf8mb4_unicode_ci COMMENT '问题JSON',
  `intro_content` text COLLATE utf8mb4_unicode_ci COMMENT '问卷介绍',
  `publish_at` datetime DEFAULT NULL COMMENT '发布时间',
  `expire_at` datetime DEFAULT NULL COMMENT '过期时间',
  `is_anonymous` tinyint NOT NULL DEFAULT '1' COMMENT '是否匿名:0=否,1=是',
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '状态:0=草稿,1=已发布,2=已下架',
  `tenant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '租户',
  `merchant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '商户',
  `created_by` bigint unsigned NOT NULL DEFAULT '0' COMMENT '创建人ID',
  `dept_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '所属部门ID',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间，非 NULL 表示已删除',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_survey_no` (`survey_no`),
  KEY `idx_tenant_merchant` (`tenant_id`, `merchant_id`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='体验问卷';

CREATE TABLE IF NOT EXISTS `demo_appointment` (
  `id` bigint unsigned NOT NULL COMMENT '预约ID（Snowflake）',
  `appointment_no` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '预约编号|search:eq|priority:100',
  `customer_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '客户',
  `subject` varchar(120) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '预约主题|search:like|keyword:on|priority:95',
  `appointment_at` datetime DEFAULT NULL COMMENT '预约时间',
  `contact_phone` varchar(30) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '联系电话|search:like|keyword:on|priority:90',
  `address` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '预约地址|keyword:only',
  `remark` text COLLATE utf8mb4_unicode_ci COMMENT '备注|keyword:only',
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '状态:0=待确认,1=已确认,2=已完成,3=已取消',
  `tenant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '租户',
  `merchant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '商户',
  `created_by` bigint unsigned NOT NULL DEFAULT '0' COMMENT '创建人ID',
  `dept_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '所属部门ID',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间，非 NULL 表示已删除',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_appointment_no` (`appointment_no`),
  KEY `idx_customer_id` (`customer_id`),
  KEY `idx_tenant_merchant` (`tenant_id`, `merchant_id`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='体验预约';

CREATE TABLE IF NOT EXISTS `demo_audit_log` (
  `id` bigint unsigned NOT NULL COMMENT '审计日志ID（Snowflake）',
  `log_no` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '日志编号|search:eq|priority:100',
  `operator_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '操作人|ref:system_users.username',
  `action` tinyint NOT NULL DEFAULT '1' COMMENT '动作:1=创建,2=修改,3=删除,4=导出,5=导入',
  `target_type` tinyint NOT NULL DEFAULT '1' COMMENT '对象类型:1=客户,2=商品,3=订单,4=工单',
  `target_code` varchar(80) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '对象编号|search:eq|priority:88',
  `request_json` text COLLATE utf8mb4_unicode_ci COMMENT '请求JSON',
  `result` tinyint NOT NULL DEFAULT '1' COMMENT '结果:0=失败,1=成功',
  `client_ip` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '客户端IP|search:eq|priority:80',
  `occurred_at` datetime DEFAULT NULL COMMENT '发生时间',
  `remark` text COLLATE utf8mb4_unicode_ci COMMENT '备注|keyword:only',
  `tenant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '租户',
  `merchant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '商户',
  `created_by` bigint unsigned NOT NULL DEFAULT '0' COMMENT '创建人ID',
  `dept_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '所属部门ID',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间，非 NULL 表示已删除',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_log_no` (`log_no`),
  KEY `idx_operator_id` (`operator_id`),
  KEY `idx_tenant_merchant` (`tenant_id`, `merchant_id`),
  KEY `idx_occurred_at` (`occurred_at`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='体验审计日志';

INSERT INTO `demo_category` (`id`, `parent_id`, `name`, `icon`, `sort`, `status`, `tenant_id`, `merchant_id`, `created_by`, `dept_id`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (330000000000000001, 0, '产品线', 'AppstoreOutlined', 10, 1, 324500000000000001, 324500000000000101, 1000000000000000003, 1000000000000000001, NOW(), NOW(), NULL),
  (330000000000000002, 330000000000000001, '软件服务', 'CloudServerOutlined', 20, 1, 324500000000000001, 324500000000000101, 1000000000000000003, 1000000000000000001, NOW(), NOW(), NULL),
  (330000000000000003, 330000000000000001, '硬件设备', 'LaptopOutlined', 30, 1, 324500000000000001, 324500000000000101, 1000000000000000003, 1000000000000000001, NOW(), NOW(), NULL),
  (330000000000000004, 0, '营销活动', 'NotificationOutlined', 40, 1, 324500000000000001, 324500000000000101, 1000000000000000003, 1000000000000000001, NOW(), NOW(), NULL)
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`), `parent_id` = VALUES(`parent_id`), `status` = VALUES(`status`), `updated_at` = NOW(), `deleted_at` = NULL;

INSERT INTO `demo_customer` (`id`, `avatar`, `name`, `customer_no`, `phone`, `email`, `gender`, `level`, `source_type`, `is_vip`, `registered_at`, `remark`, `status`, `tenant_id`, `merchant_id`, `created_by`, `dept_id`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (330000000000001001, 'https://picsum.photos/seed/baseadmin-customer-a/96/96', '上海晴川科技', 'CUS-20260428-001', '13800001001', 'ops@example.com', 0, 2, 1, 1, '2026-04-20 09:30:00', '重点客户，关注续费和增购。', 1, 324500000000000001, 324500000000000101, 1000000000000000003, 1000000000000000001, NOW(), NOW(), NULL),
  (330000000000001002, 'https://picsum.photos/seed/baseadmin-customer-b/96/96', '杭州云帆贸易', 'CUS-20260428-002', '13800001002', 'buyer@example.com', 1, 1, 2, 0, '2026-04-21 10:15:00', '线下展会留资。', 1, 324500000000000001, 324500000000000101, 1000000000000000003, 1000000000000000001, NOW(), NOW(), NULL)
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`), `phone` = VALUES(`phone`), `email` = VALUES(`email`), `status` = VALUES(`status`), `updated_at` = NOW(), `deleted_at` = NULL;

INSERT INTO `demo_product` (`id`, `category_id`, `sku_no`, `name`, `cover`, `manual_file`, `detail_content`, `spec_json`, `website_url`, `type`, `is_recommend`, `sale_price`, `stock_num`, `weight_num`, `sort`, `icon`, `status`, `tenant_id`, `merchant_id`, `created_by`, `dept_id`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (330000000000002001, 330000000000000002, 'SKU-DEMO-SAAS-001', 'SaaS基础版', 'https://picsum.photos/seed/baseadmin-product-a/320/180', 'https://example.com/demo/manual-basic.pdf', '<p>SaaS基础版，适合中小团队快速上线。</p>', '{\"users\":20,\"storage\":\"50GB\",\"support\":\"workday\"}', 'https://baseadmin.easytestdev.online/admin/', 2, 1, 199900, 120, 0, 10, 'CloudOutlined', 1, 324500000000000001, 324500000000000101, 1000000000000000003, 1000000000000000001, NOW(), NOW(), NULL),
  (330000000000002002, 330000000000000003, 'SKU-DEMO-BOX-001', '边缘网关盒子', 'https://picsum.photos/seed/baseadmin-product-b/320/180', 'https://example.com/demo/manual-box.pdf', '<p>适配门店和设备采集场景。</p>', '{\"cpu\":\"4C\",\"memory\":\"8GB\",\"network\":\"dual-lan\"}', 'https://easytestdev.online/', 3, 0, 89900, 42, 680, 20, 'DeploymentUnitOutlined', 1, 324500000000000001, 324500000000000101, 1000000000000000003, 1000000000000000001, NOW(), NOW(), NULL)
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`), `category_id` = VALUES(`category_id`), `status` = VALUES(`status`), `updated_at` = NOW(), `deleted_at` = NULL;

INSERT INTO `demo_campaign` (`id`, `campaign_no`, `title`, `banner`, `type`, `channel`, `budget_amount`, `landing_url`, `rule_json`, `intro_content`, `start_at`, `end_at`, `is_public`, `status`, `tenant_id`, `merchant_id`, `created_by`, `dept_id`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (330000000000003001, 'CMP-20260428-001', '五一续费优惠', 'https://picsum.photos/seed/baseadmin-campaign-a/640/240', 2, 1, 500000, 'https://baseadmin.easytestdev.online/admin/', '{\"discount\":\"20%\",\"minAmount\":100000}', '<p>面向存量客户的续费优惠活动。</p>', '2026-04-28 00:00:00', '2026-05-08 23:59:59', 1, 1, 324500000000000001, 324500000000000101, 1000000000000000003, 1000000000000000001, NOW(), NOW(), NULL)
ON DUPLICATE KEY UPDATE `title` = VALUES(`title`), `status` = VALUES(`status`), `updated_at` = NOW(), `deleted_at` = NULL;

INSERT INTO `demo_order` (`id`, `order_no`, `customer_id`, `product_id`, `quantity`, `amount`, `pay_status`, `deliver_status`, `paid_at`, `deliver_at`, `receiver_phone`, `address`, `remark`, `status`, `tenant_id`, `merchant_id`, `created_by`, `dept_id`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (330000000000004001, 'ORD-20260428-001', 330000000000001001, 330000000000002001, 2, 399800, 1, 1, '2026-04-28 09:45:00', '2026-04-28 14:00:00', '13800001001', '上海市浦东新区 Demo 路 88 号', '优先安排客户成功回访。', 1, 324500000000000001, 324500000000000101, 1000000000000000003, 1000000000000000001, NOW(), NOW(), NULL),
  (330000000000004002, 'ORD-20260428-002', 330000000000001002, 330000000000002002, 1, 89900, 0, 0, NULL, NULL, '13800001002', '杭州市西湖区 Demo 街 18 号', '待客户付款。', 0, 324500000000000001, 324500000000000101, 1000000000000000003, 1000000000000000001, NOW(), NOW(), NULL)
ON DUPLICATE KEY UPDATE `customer_id` = VALUES(`customer_id`), `product_id` = VALUES(`product_id`), `status` = VALUES(`status`), `updated_at` = NOW(), `deleted_at` = NULL;

INSERT INTO `demo_work_order` (`id`, `ticket_no`, `customer_id`, `product_id`, `order_id`, `title`, `priority`, `source_type`, `description`, `attachment_file`, `due_at`, `status`, `tenant_id`, `merchant_id`, `created_by`, `dept_id`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (330000000000005001, 'TCK-20260428-001', 330000000000001001, 330000000000002001, 330000000000004001, '发票抬头需要变更', 2, 4, '客户希望将发票抬头改为集团公司名称。', 'https://example.com/demo/invoice-request.pdf', '2026-04-30 18:00:00', 1, 324500000000000001, 324500000000000101, 1000000000000000003, 1000000000000000001, NOW(), NOW(), NULL)
ON DUPLICATE KEY UPDATE `title` = VALUES(`title`), `status` = VALUES(`status`), `updated_at` = NOW(), `deleted_at` = NULL;

INSERT INTO `demo_contract` (`id`, `contract_no`, `customer_id`, `order_id`, `title`, `contract_file`, `sign_image`, `contract_amount`, `sign_password`, `signed_at`, `expires_at`, `status`, `tenant_id`, `merchant_id`, `created_by`, `dept_id`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (330000000000006001, 'CON-20260428-001', 330000000000001001, 330000000000004001, 'SaaS基础版年度服务合同', 'https://example.com/demo/contract.pdf', 'https://picsum.photos/seed/baseadmin-sign-a/160/80', 399800, '$2a$10$7EqJtq98hPqEX7fNZaFWoOhiR0ykfHXRh2QHjT8p4DDEpA3WJZT9G', '2026-04-28 11:00:00', '2027-04-27 23:59:59', 1, 324500000000000001, 324500000000000101, 1000000000000000003, 1000000000000000001, NOW(), NOW(), NULL)
ON DUPLICATE KEY UPDATE `title` = VALUES(`title`), `status` = VALUES(`status`), `updated_at` = NOW(), `deleted_at` = NULL;

INSERT INTO `demo_survey` (`id`, `survey_no`, `title`, `poster`, `question_json`, `intro_content`, `publish_at`, `expire_at`, `is_anonymous`, `status`, `tenant_id`, `merchant_id`, `created_by`, `dept_id`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (330000000000007001, 'SUR-20260428-001', '客户满意度调研', 'https://picsum.photos/seed/baseadmin-survey-a/320/180', '{\"questions\":[{\"type\":\"rate\",\"title\":\"服务满意度\"},{\"type\":\"text\",\"title\":\"改进建议\"}]}', '<p>用于收集客户对实施、培训、售后的反馈。</p>', '2026-04-28 08:00:00', '2026-05-31 23:59:59', 1, 1, 324500000000000001, 324500000000000101, 1000000000000000003, 1000000000000000001, NOW(), NOW(), NULL)
ON DUPLICATE KEY UPDATE `title` = VALUES(`title`), `status` = VALUES(`status`), `updated_at` = NOW(), `deleted_at` = NULL;

INSERT INTO `demo_appointment` (`id`, `appointment_no`, `customer_id`, `subject`, `appointment_at`, `contact_phone`, `address`, `remark`, `status`, `tenant_id`, `merchant_id`, `created_by`, `dept_id`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (330000000000008001, 'APT-20260428-001', 330000000000001001, '远程实施沟通', '2026-04-29 15:00:00', '13800001001', '线上会议', '确认初始化数据和权限边界。', 1, 324500000000000001, 324500000000000101, 1000000000000000003, 1000000000000000001, NOW(), NOW(), NULL)
ON DUPLICATE KEY UPDATE `subject` = VALUES(`subject`), `status` = VALUES(`status`), `updated_at` = NOW(), `deleted_at` = NULL;

INSERT INTO `demo_audit_log` (`id`, `log_no`, `operator_id`, `action`, `target_type`, `target_code`, `request_json`, `result`, `client_ip`, `occurred_at`, `remark`, `tenant_id`, `merchant_id`, `created_by`, `dept_id`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (330000000000009001, 'LOG-20260428-001', 1000000000000000003, 1, 3, 'ORD-20260428-001', '{\"module\":\"demo_order\",\"action\":\"create\"}', 1, '127.0.0.1', '2026-04-28 12:00:00', '演示审计日志。', 324500000000000001, 324500000000000101, 1000000000000000003, 1000000000000000001, NOW(), NOW(), NULL)
ON DUPLICATE KEY UPDATE `target_code` = VALUES(`target_code`), `result` = VALUES(`result`), `updated_at` = NOW(), `deleted_at` = NULL;

-- Demo 菜单固定种子，保证迁移后超级管理员可以直接看到动态路由。
INSERT INTO `system_menu` (`id`, `parent_id`, `title`, `type`, `path`, `component`, `permission`, `icon`, `sort`, `is_show`, `is_cache`, `link_url`, `status`, `created_by`, `dept_id`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (1000000000000000900, 0, 'Demo体验', 1, '/demo', NULL, '', 'ExperimentOutlined', 90, 1, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000910, 1000000000000000900, '体验分类', 2, '/demo/category', 'demo/category/index', 'demo:category:list', 'ApartmentOutlined', 10, 1, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000920, 1000000000000000900, '体验客户', 2, '/demo/customer', 'demo/customer/index', 'demo:customer:list', 'UserOutlined', 20, 1, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000930, 1000000000000000900, '体验商品', 2, '/demo/product', 'demo/product/index', 'demo:product:list', 'AppstoreOutlined', 30, 1, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000940, 1000000000000000900, '体验活动', 2, '/demo/campaign', 'demo/campaign/index', 'demo:campaign:list', 'NotificationOutlined', 40, 1, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000950, 1000000000000000900, '体验订单', 2, '/demo/order', 'demo/order/index', 'demo:order:list', 'ShoppingCartOutlined', 50, 1, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000960, 1000000000000000900, '体验工单', 2, '/demo/work-order', 'demo/work_order/index', 'demo:work_order:list', 'ToolOutlined', 60, 1, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000970, 1000000000000000900, '体验合同', 2, '/demo/contract', 'demo/contract/index', 'demo:contract:list', 'FileTextOutlined', 70, 1, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000980, 1000000000000000900, '体验问卷', 2, '/demo/survey', 'demo/survey/index', 'demo:survey:list', 'FormOutlined', 80, 1, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000990, 1000000000000000900, '体验预约', 2, '/demo/appointment', 'demo/appointment/index', 'demo:appointment:list', 'CalendarOutlined', 90, 1, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000001000, 1000000000000000900, '体验审计日志', 2, '/demo/audit-log', 'demo/audit_log/index', 'demo:audit_log:list', 'AuditOutlined', 100, 1, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL)
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
  `updated_at` = NOW(),
  `deleted_at` = NULL;

INSERT INTO `system_menu` (`id`, `parent_id`, `title`, `type`, `path`, `component`, `permission`, `icon`, `sort`, `is_show`, `is_cache`, `link_url`, `status`, `created_by`, `dept_id`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (1000000000000000911, 1000000000000000910, '体验分类新增', 3, NULL, NULL, 'demo:category:create', '', 1, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000912, 1000000000000000910, '体验分类修改', 3, NULL, NULL, 'demo:category:update', '', 2, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000913, 1000000000000000910, '体验分类删除', 3, NULL, NULL, 'demo:category:delete', '', 3, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000914, 1000000000000000910, '体验分类批量删除', 3, NULL, NULL, 'demo:category:batch-delete', '', 4, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000915, 1000000000000000910, '体验分类查看', 3, NULL, NULL, 'demo:category:detail', '', 5, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000916, 1000000000000000910, '体验分类导出', 3, NULL, NULL, 'demo:category:export', '', 6, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000917, 1000000000000000910, '体验分类批量编辑', 3, NULL, NULL, 'demo:category:batch-update', '', 8, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000921, 1000000000000000920, '体验客户新增', 3, NULL, NULL, 'demo:customer:create', '', 1, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000922, 1000000000000000920, '体验客户修改', 3, NULL, NULL, 'demo:customer:update', '', 2, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000923, 1000000000000000920, '体验客户删除', 3, NULL, NULL, 'demo:customer:delete', '', 3, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000924, 1000000000000000920, '体验客户批量删除', 3, NULL, NULL, 'demo:customer:batch-delete', '', 4, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000925, 1000000000000000920, '体验客户查看', 3, NULL, NULL, 'demo:customer:detail', '', 5, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000926, 1000000000000000920, '体验客户导出', 3, NULL, NULL, 'demo:customer:export', '', 6, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000927, 1000000000000000920, '体验客户导入', 3, NULL, NULL, 'demo:customer:import', '', 7, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000928, 1000000000000000920, '体验客户批量编辑', 3, NULL, NULL, 'demo:customer:batch-update', '', 8, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000931, 1000000000000000930, '体验商品新增', 3, NULL, NULL, 'demo:product:create', '', 1, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000932, 1000000000000000930, '体验商品修改', 3, NULL, NULL, 'demo:product:update', '', 2, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000933, 1000000000000000930, '体验商品删除', 3, NULL, NULL, 'demo:product:delete', '', 3, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000934, 1000000000000000930, '体验商品批量删除', 3, NULL, NULL, 'demo:product:batch-delete', '', 4, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000935, 1000000000000000930, '体验商品查看', 3, NULL, NULL, 'demo:product:detail', '', 5, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000936, 1000000000000000930, '体验商品导出', 3, NULL, NULL, 'demo:product:export', '', 6, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000937, 1000000000000000930, '体验商品导入', 3, NULL, NULL, 'demo:product:import', '', 7, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000938, 1000000000000000930, '体验商品批量编辑', 3, NULL, NULL, 'demo:product:batch-update', '', 8, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000941, 1000000000000000940, '体验活动新增', 3, NULL, NULL, 'demo:campaign:create', '', 1, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000942, 1000000000000000940, '体验活动修改', 3, NULL, NULL, 'demo:campaign:update', '', 2, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000943, 1000000000000000940, '体验活动删除', 3, NULL, NULL, 'demo:campaign:delete', '', 3, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000944, 1000000000000000940, '体验活动批量删除', 3, NULL, NULL, 'demo:campaign:batch-delete', '', 4, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000945, 1000000000000000940, '体验活动查看', 3, NULL, NULL, 'demo:campaign:detail', '', 5, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000946, 1000000000000000940, '体验活动导出', 3, NULL, NULL, 'demo:campaign:export', '', 6, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000947, 1000000000000000940, '体验活动导入', 3, NULL, NULL, 'demo:campaign:import', '', 7, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000948, 1000000000000000940, '体验活动批量编辑', 3, NULL, NULL, 'demo:campaign:batch-update', '', 8, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000951, 1000000000000000950, '体验订单新增', 3, NULL, NULL, 'demo:order:create', '', 1, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000952, 1000000000000000950, '体验订单修改', 3, NULL, NULL, 'demo:order:update', '', 2, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000953, 1000000000000000950, '体验订单删除', 3, NULL, NULL, 'demo:order:delete', '', 3, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000954, 1000000000000000950, '体验订单批量删除', 3, NULL, NULL, 'demo:order:batch-delete', '', 4, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000955, 1000000000000000950, '体验订单查看', 3, NULL, NULL, 'demo:order:detail', '', 5, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000956, 1000000000000000950, '体验订单导出', 3, NULL, NULL, 'demo:order:export', '', 6, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000957, 1000000000000000950, '体验订单导入', 3, NULL, NULL, 'demo:order:import', '', 7, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000958, 1000000000000000950, '体验订单批量编辑', 3, NULL, NULL, 'demo:order:batch-update', '', 8, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000961, 1000000000000000960, '体验工单新增', 3, NULL, NULL, 'demo:work_order:create', '', 1, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000962, 1000000000000000960, '体验工单修改', 3, NULL, NULL, 'demo:work_order:update', '', 2, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000963, 1000000000000000960, '体验工单删除', 3, NULL, NULL, 'demo:work_order:delete', '', 3, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000964, 1000000000000000960, '体验工单批量删除', 3, NULL, NULL, 'demo:work_order:batch-delete', '', 4, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000965, 1000000000000000960, '体验工单查看', 3, NULL, NULL, 'demo:work_order:detail', '', 5, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000966, 1000000000000000960, '体验工单导出', 3, NULL, NULL, 'demo:work_order:export', '', 6, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000967, 1000000000000000960, '体验工单导入', 3, NULL, NULL, 'demo:work_order:import', '', 7, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000968, 1000000000000000960, '体验工单批量编辑', 3, NULL, NULL, 'demo:work_order:batch-update', '', 8, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000971, 1000000000000000970, '体验合同新增', 3, NULL, NULL, 'demo:contract:create', '', 1, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000972, 1000000000000000970, '体验合同修改', 3, NULL, NULL, 'demo:contract:update', '', 2, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000973, 1000000000000000970, '体验合同删除', 3, NULL, NULL, 'demo:contract:delete', '', 3, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000974, 1000000000000000970, '体验合同批量删除', 3, NULL, NULL, 'demo:contract:batch-delete', '', 4, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000975, 1000000000000000970, '体验合同查看', 3, NULL, NULL, 'demo:contract:detail', '', 5, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000976, 1000000000000000970, '体验合同导出', 3, NULL, NULL, 'demo:contract:export', '', 6, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000977, 1000000000000000970, '体验合同导入', 3, NULL, NULL, 'demo:contract:import', '', 7, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000978, 1000000000000000970, '体验合同批量编辑', 3, NULL, NULL, 'demo:contract:batch-update', '', 8, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000981, 1000000000000000980, '体验问卷新增', 3, NULL, NULL, 'demo:survey:create', '', 1, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000982, 1000000000000000980, '体验问卷修改', 3, NULL, NULL, 'demo:survey:update', '', 2, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000983, 1000000000000000980, '体验问卷删除', 3, NULL, NULL, 'demo:survey:delete', '', 3, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000984, 1000000000000000980, '体验问卷批量删除', 3, NULL, NULL, 'demo:survey:batch-delete', '', 4, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000985, 1000000000000000980, '体验问卷查看', 3, NULL, NULL, 'demo:survey:detail', '', 5, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000986, 1000000000000000980, '体验问卷导出', 3, NULL, NULL, 'demo:survey:export', '', 6, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000987, 1000000000000000980, '体验问卷导入', 3, NULL, NULL, 'demo:survey:import', '', 7, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000988, 1000000000000000980, '体验问卷批量编辑', 3, NULL, NULL, 'demo:survey:batch-update', '', 8, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000991, 1000000000000000990, '体验预约新增', 3, NULL, NULL, 'demo:appointment:create', '', 1, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000992, 1000000000000000990, '体验预约修改', 3, NULL, NULL, 'demo:appointment:update', '', 2, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000993, 1000000000000000990, '体验预约删除', 3, NULL, NULL, 'demo:appointment:delete', '', 3, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000994, 1000000000000000990, '体验预约批量删除', 3, NULL, NULL, 'demo:appointment:batch-delete', '', 4, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000995, 1000000000000000990, '体验预约查看', 3, NULL, NULL, 'demo:appointment:detail', '', 5, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000996, 1000000000000000990, '体验预约导出', 3, NULL, NULL, 'demo:appointment:export', '', 6, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000997, 1000000000000000990, '体验预约导入', 3, NULL, NULL, 'demo:appointment:import', '', 7, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000000998, 1000000000000000990, '体验预约批量编辑', 3, NULL, NULL, 'demo:appointment:batch-update', '', 8, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000001001, 1000000000000001000, '体验审计日志新增', 3, NULL, NULL, 'demo:audit_log:create', '', 1, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000001002, 1000000000000001000, '体验审计日志修改', 3, NULL, NULL, 'demo:audit_log:update', '', 2, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000001003, 1000000000000001000, '体验审计日志删除', 3, NULL, NULL, 'demo:audit_log:delete', '', 3, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000001004, 1000000000000001000, '体验审计日志批量删除', 3, NULL, NULL, 'demo:audit_log:batch-delete', '', 4, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000001005, 1000000000000001000, '体验审计日志查看', 3, NULL, NULL, 'demo:audit_log:detail', '', 5, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000001006, 1000000000000001000, '体验审计日志导出', 3, NULL, NULL, 'demo:audit_log:export', '', 6, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000001007, 1000000000000001000, '体验审计日志导入', 3, NULL, NULL, 'demo:audit_log:import', '', 7, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL)
ON DUPLICATE KEY UPDATE
  `parent_id` = VALUES(`parent_id`),
  `title` = VALUES(`title`),
  `permission` = VALUES(`permission`),
  `sort` = VALUES(`sort`),
  `status` = VALUES(`status`),
  `updated_at` = NOW(),
  `deleted_at` = NULL;

INSERT INTO `system_role_menu` (`role_id`, `menu_id`)
SELECT 1000000000000000002, `id`
FROM `system_menu`
WHERE (`path` LIKE '/demo%' OR `permission` LIKE 'demo:%') AND `deleted_at` IS NULL
ON DUPLICATE KEY UPDATE `menu_id` = VALUES(`menu_id`);

SET FOREIGN_KEY_CHECKS = 1;
