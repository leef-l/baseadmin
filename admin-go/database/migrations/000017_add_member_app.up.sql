-- 会员分销系统
-- 包含：会员管理、等级管理、钱包、商城、仓库寄售、团队数据导出

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ================================================================
-- 1. 会员等级配置表
-- ================================================================
CREATE TABLE IF NOT EXISTS `member_level` (
  `id` bigint unsigned NOT NULL COMMENT '等级ID（Snowflake）',
  `name` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '等级名称|search:like|keyword:on|priority:100',
  `level_no` int unsigned NOT NULL DEFAULT '0' COMMENT '等级编号（越大越高）|search:eq',
  `icon` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '等级图标',
  `duration_days` int unsigned NOT NULL DEFAULT '0' COMMENT '有效天数（0=永久）',
  `need_active_count` int unsigned NOT NULL DEFAULT '0' COMMENT '升级所需有效用户数',
  `need_team_turnover` bigint unsigned NOT NULL DEFAULT '0' COMMENT '升级所需团队营业额（分）',
  `is_top` tinyint NOT NULL DEFAULT '0' COMMENT '是否最高等级:0=否,1=是|search:select',
  `auto_deploy` tinyint NOT NULL DEFAULT '0' COMMENT '到达后自动部署站点:0=否,1=是',
  `remark` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '等级说明|search:off',
  `sort` int NOT NULL DEFAULT '0' COMMENT '排序（升序）',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态:0=关闭,1=开启|search:select',
  `tenant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '租户',
  `merchant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '商户',
  `created_by` bigint unsigned NOT NULL DEFAULT '0' COMMENT '创建人ID',
  `dept_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '所属部门ID',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间，非 NULL 表示已删除',
  PRIMARY KEY (`id`),
  KEY `idx_level_no` (`level_no`),
  KEY `idx_tenant_merchant` (`tenant_id`, `merchant_id`),
  KEY `idx_status` (`status`),
  KEY `idx_sort` (`sort`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='会员等级配置';

-- ================================================================
-- 2. 会员用户表（核心，parent_id 树形）
-- ================================================================
CREATE TABLE IF NOT EXISTS `member_user` (
  `id` bigint unsigned NOT NULL COMMENT '会员ID（Snowflake）',
  `parent_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '上级会员',
  `username` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户名（登录账号）|search:eq|keyword:on|priority:100',
  `password` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '密码（bcrypt加密）',
  `nickname` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '昵称|search:like|keyword:on|priority:95',
  `phone` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '手机号|search:eq|keyword:on|priority:90',
  `avatar` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '头像',
  `real_name` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '真实姓名|search:like|keyword:on',
  `level_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '当前等级|ref:member_level.name|search:select',
  `level_expire_at` datetime DEFAULT NULL COMMENT '等级到期时间',
  `team_count` int unsigned NOT NULL DEFAULT '0' COMMENT '团队总人数',
  `direct_count` int unsigned NOT NULL DEFAULT '0' COMMENT '直推人数',
  `active_count` int unsigned NOT NULL DEFAULT '0' COMMENT '有效用户数',
  `team_turnover` bigint unsigned NOT NULL DEFAULT '0' COMMENT '团队总营业额（分）',
  `is_active` tinyint NOT NULL DEFAULT '1' COMMENT '是否激活:0=未激活,1=已激活|search:select',
  `is_qualified` tinyint NOT NULL DEFAULT '1' COMMENT '仓库资格:0=已失效,1=有效|search:select',
  `invite_code` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '邀请码|search:eq',
  `register_ip` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '注册IP',
  `last_login_at` datetime DEFAULT NULL COMMENT '最后登录时间',
  `remark` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '备注|search:off',
  `sort` int NOT NULL DEFAULT '0' COMMENT '排序（升序）',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态:0=冻结,1=正常|search:select',
  `tenant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '租户',
  `merchant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '商户',
  `created_by` bigint unsigned NOT NULL DEFAULT '0' COMMENT '创建人ID',
  `dept_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '所属部门ID',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间，非 NULL 表示已删除',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`),
  UNIQUE KEY `uk_invite_code` (`invite_code`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_level_id` (`level_id`),
  KEY `idx_phone` (`phone`),
  KEY `idx_tenant_merchant` (`tenant_id`, `merchant_id`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='会员用户';

-- ================================================================
-- 3. 等级变更日志
-- ================================================================
CREATE TABLE IF NOT EXISTS `member_level_log` (
  `id` bigint unsigned NOT NULL COMMENT 'ID（Snowflake）',
  `user_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '会员|ref:member_user.nickname|search:select',
  `old_level_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '变更前等级|ref:member_level.name',
  `new_level_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '变更后等级|ref:member_level.name',
  `change_type` tinyint NOT NULL DEFAULT '1' COMMENT '变更类型:1=自动升级,2=后台调整,3=过期降级|search:select',
  `expire_at` datetime DEFAULT NULL COMMENT '新等级到期时间',
  `remark` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '变更说明|search:off',
  `tenant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '租户',
  `merchant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '商户',
  `created_by` bigint unsigned NOT NULL DEFAULT '0' COMMENT '创建人ID',
  `dept_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '所属部门ID',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间，非 NULL 表示已删除',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_tenant_merchant` (`tenant_id`, `merchant_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='等级变更日志';

-- ================================================================
-- 4. 会员钱包表（每人3条：优惠券/奖金/推广奖）
-- ================================================================
CREATE TABLE IF NOT EXISTS `member_wallet` (
  `id` bigint unsigned NOT NULL COMMENT 'ID（Snowflake）',
  `user_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '会员|ref:member_user.nickname|search:select',
  `wallet_type` tinyint NOT NULL DEFAULT '1' COMMENT '钱包类型:1=优惠券余额,2=奖金余额,3=推广奖余额|search:select',
  `balance` bigint NOT NULL DEFAULT '0' COMMENT '当前余额（分）',
  `total_income` bigint unsigned NOT NULL DEFAULT '0' COMMENT '累计收入（分）',
  `total_expense` bigint unsigned NOT NULL DEFAULT '0' COMMENT '累计支出（分）',
  `frozen_amount` bigint unsigned NOT NULL DEFAULT '0' COMMENT '冻结金额（分）',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态:0=冻结,1=正常|search:select',
  `tenant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '租户',
  `merchant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '商户',
  `created_by` bigint unsigned NOT NULL DEFAULT '0' COMMENT '创建人ID',
  `dept_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '所属部门ID',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间，非 NULL 表示已删除',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_type` (`user_id`, `wallet_type`),
  KEY `idx_tenant_merchant` (`tenant_id`, `merchant_id`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='会员钱包';

-- ================================================================
-- 5. 钱包流水记录
-- ================================================================
CREATE TABLE IF NOT EXISTS `member_wallet_log` (
  `id` bigint unsigned NOT NULL COMMENT 'ID（Snowflake）',
  `user_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '会员|ref:member_user.nickname|search:select',
  `wallet_type` tinyint NOT NULL DEFAULT '1' COMMENT '钱包类型:1=优惠券余额,2=奖金余额,3=推广奖余额|search:select',
  `change_type` tinyint NOT NULL DEFAULT '1' COMMENT '变动类型:1=充值,2=消费,3=推广奖,4=仓库卖出收入,5=平台扣除,6=后台调整|search:select',
  `change_amount` bigint NOT NULL DEFAULT '0' COMMENT '变动金额（分，正增负减）',
  `before_balance` bigint NOT NULL DEFAULT '0' COMMENT '变动前余额（分）',
  `after_balance` bigint NOT NULL DEFAULT '0' COMMENT '变动后余额（分）',
  `related_order_no` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '关联单号|search:eq|keyword:on',
  `remark` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '备注说明|search:off',
  `tenant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '租户',
  `merchant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '商户',
  `created_by` bigint unsigned NOT NULL DEFAULT '0' COMMENT '创建人ID',
  `dept_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '所属部门ID',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间，非 NULL 表示已删除',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_wallet_type` (`wallet_type`),
  KEY `idx_related_order_no` (`related_order_no`),
  KEY `idx_tenant_merchant` (`tenant_id`, `merchant_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='钱包流水记录';

-- ================================================================
-- 6. 换绑上级日志
-- ================================================================
CREATE TABLE IF NOT EXISTS `member_rebind_log` (
  `id` bigint unsigned NOT NULL COMMENT 'ID（Snowflake）',
  `user_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '会员|ref:member_user.nickname|search:select',
  `old_parent_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '原上级|ref:member_user.nickname',
  `new_parent_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '新上级|ref:member_user.nickname',
  `reason` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '换绑原因|search:off',
  `operator_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '操作人|ref:system_users.username',
  `tenant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '租户',
  `merchant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '商户',
  `created_by` bigint unsigned NOT NULL DEFAULT '0' COMMENT '创建人ID',
  `dept_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '所属部门ID',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间，非 NULL 表示已删除',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_tenant_merchant` (`tenant_id`, `merchant_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='换绑上级日志';

-- ================================================================
-- 7. 商城商品分类（树形 parent_id）
-- ================================================================
CREATE TABLE IF NOT EXISTS `member_shop_category` (
  `id` bigint unsigned NOT NULL COMMENT 'ID（Snowflake）',
  `parent_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '上级分类',
  `name` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '分类名称|search:like|keyword:on|priority:100',
  `icon` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '分类图标',
  `sort` int NOT NULL DEFAULT '0' COMMENT '排序（升序）',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态:0=关闭,1=开启|search:select',
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商城商品分类';

-- ================================================================
-- 8. 商城商品
-- ================================================================
CREATE TABLE IF NOT EXISTS `member_shop_goods` (
  `id` bigint unsigned NOT NULL COMMENT 'ID（Snowflake）',
  `category_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '商品分类|ref:member_shop_category.name|search:select',
  `title` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '商品名称|search:like|keyword:on|priority:100',
  `cover` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '封面图',
  `images` text COLLATE utf8mb4_unicode_ci COMMENT '商品图片（JSON数组）|search:off',
  `price` bigint unsigned NOT NULL DEFAULT '0' COMMENT '售价（分，优惠券余额支付）',
  `original_price` bigint unsigned NOT NULL DEFAULT '0' COMMENT '原价（分）',
  `stock` int unsigned NOT NULL DEFAULT '0' COMMENT '库存',
  `sales` int unsigned NOT NULL DEFAULT '0' COMMENT '销量',
  `content` text COLLATE utf8mb4_unicode_ci COMMENT '商品详情|search:off',
  `sort` int NOT NULL DEFAULT '0' COMMENT '排序（升序）',
  `is_recommend` tinyint NOT NULL DEFAULT '0' COMMENT '是否推荐:0=否,1=是|search:select',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态:0=下架,1=上架|search:select',
  `tenant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '租户',
  `merchant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '商户',
  `created_by` bigint unsigned NOT NULL DEFAULT '0' COMMENT '创建人ID',
  `dept_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '所属部门ID',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间，非 NULL 表示已删除',
  PRIMARY KEY (`id`),
  KEY `idx_category_id` (`category_id`),
  KEY `idx_tenant_merchant` (`tenant_id`, `merchant_id`),
  KEY `idx_status` (`status`),
  KEY `idx_sort` (`sort`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商城商品';

-- ================================================================
-- 9. 商城订单
-- ================================================================
CREATE TABLE IF NOT EXISTS `member_shop_order` (
  `id` bigint unsigned NOT NULL COMMENT 'ID（Snowflake）',
  `order_no` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '订单号|search:eq|keyword:on|priority:100',
  `user_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '购买会员|ref:member_user.nickname|search:select',
  `goods_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '商品|ref:member_shop_goods.title|search:select',
  `goods_title` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '商品名称（快照）',
  `goods_cover` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '商品封面（快照）',
  `quantity` int unsigned NOT NULL DEFAULT '1' COMMENT '购买数量',
  `total_price` bigint unsigned NOT NULL DEFAULT '0' COMMENT '订单总价（分）',
  `pay_wallet` tinyint NOT NULL DEFAULT '1' COMMENT '支付钱包:1=优惠券余额',
  `order_status` tinyint NOT NULL DEFAULT '1' COMMENT '订单状态:1=已完成,2=已取消|search:select',
  `remark` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '订单备注|search:off',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态:0=关闭,1=开启|search:select',
  `tenant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '租户',
  `merchant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '商户',
  `created_by` bigint unsigned NOT NULL DEFAULT '0' COMMENT '创建人ID',
  `dept_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '所属部门ID',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间，非 NULL 表示已删除',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_order_no` (`order_no`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_goods_id` (`goods_id`),
  KEY `idx_tenant_merchant` (`tenant_id`, `merchant_id`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商城订单';

-- ================================================================
-- 10. 仓库商品（寄售）
-- ================================================================
CREATE TABLE IF NOT EXISTS `member_warehouse_goods` (
  `id` bigint unsigned NOT NULL COMMENT 'ID（Snowflake）',
  `goods_no` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '商品编号|search:eq|keyword:on|priority:100',
  `title` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '商品名称|search:like|keyword:on|priority:95',
  `cover` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '商品封面',
  `init_price` bigint unsigned NOT NULL DEFAULT '0' COMMENT '初始价格（分）',
  `current_price` bigint unsigned NOT NULL DEFAULT '0' COMMENT '当前价格（分）',
  `price_rise_rate` int unsigned NOT NULL DEFAULT '0' COMMENT '每次加价比例（百分比，如10=10%）',
  `platform_fee_rate` int unsigned NOT NULL DEFAULT '0' COMMENT '平台扣除比例（百分比，如5=5%）',
  `owner_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '当前持有人|ref:member_user.nickname|search:select',
  `trade_count` int unsigned NOT NULL DEFAULT '0' COMMENT '流转次数',
  `goods_status` tinyint NOT NULL DEFAULT '1' COMMENT '商品状态:1=持有中,2=挂卖中,3=交易中|search:select',
  `remark` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '备注|search:off',
  `sort` int NOT NULL DEFAULT '0' COMMENT '排序（升序）',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态:0=关闭,1=开启|search:select',
  `tenant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '租户',
  `merchant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '商户',
  `created_by` bigint unsigned NOT NULL DEFAULT '0' COMMENT '创建人ID',
  `dept_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '所属部门ID',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间，非 NULL 表示已删除',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_goods_no` (`goods_no`),
  KEY `idx_owner_id` (`owner_id`),
  KEY `idx_goods_status` (`goods_status`),
  KEY `idx_tenant_merchant` (`tenant_id`, `merchant_id`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='仓库商品';

-- ================================================================
-- 11. 仓库挂卖记录
-- ================================================================
CREATE TABLE IF NOT EXISTS `member_warehouse_listing` (
  `id` bigint unsigned NOT NULL COMMENT 'ID（Snowflake）',
  `goods_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '仓库商品|ref:member_warehouse_goods.title|search:select',
  `seller_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '卖家|ref:member_user.nickname|search:select',
  `listing_price` bigint unsigned NOT NULL DEFAULT '0' COMMENT '挂卖价格（分，自动加价后）',
  `listing_status` tinyint NOT NULL DEFAULT '1' COMMENT '挂卖状态:1=挂卖中,2=已售出,3=已取消|search:select',
  `listed_at` datetime DEFAULT NULL COMMENT '挂卖时间',
  `sold_at` datetime DEFAULT NULL COMMENT '售出时间',
  `remark` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '备注|search:off',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态:0=关闭,1=开启|search:select',
  `tenant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '租户',
  `merchant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '商户',
  `created_by` bigint unsigned NOT NULL DEFAULT '0' COMMENT '创建人ID',
  `dept_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '所属部门ID',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间，非 NULL 表示已删除',
  PRIMARY KEY (`id`),
  KEY `idx_goods_id` (`goods_id`),
  KEY `idx_seller_id` (`seller_id`),
  KEY `idx_listing_status` (`listing_status`),
  KEY `idx_tenant_merchant` (`tenant_id`, `merchant_id`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='仓库挂卖记录';

-- ================================================================
-- 12. 仓库交易记录
-- ================================================================
CREATE TABLE IF NOT EXISTS `member_warehouse_trade` (
  `id` bigint unsigned NOT NULL COMMENT 'ID（Snowflake）',
  `trade_no` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '交易编号|search:eq|keyword:on|priority:100',
  `goods_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '仓库商品|ref:member_warehouse_goods.title|search:select',
  `listing_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '挂卖记录|ref:member_warehouse_listing.id',
  `seller_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '卖家|ref:member_user.nickname|search:select',
  `buyer_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '买家|ref:member_user.nickname|search:select',
  `trade_price` bigint unsigned NOT NULL DEFAULT '0' COMMENT '成交价格（分）',
  `platform_fee` bigint unsigned NOT NULL DEFAULT '0' COMMENT '平台扣除费用（分）',
  `seller_income` bigint unsigned NOT NULL DEFAULT '0' COMMENT '卖家实收（分）',
  `trade_status` tinyint NOT NULL DEFAULT '1' COMMENT '交易状态:1=待卖家确认,2=已确认完成,3=已取消|search:select',
  `confirmed_at` datetime DEFAULT NULL COMMENT '确认时间',
  `remark` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '备注|search:off',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态:0=关闭,1=开启|search:select',
  `tenant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '租户',
  `merchant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '商户',
  `created_by` bigint unsigned NOT NULL DEFAULT '0' COMMENT '创建人ID',
  `dept_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '所属部门ID',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间，非 NULL 表示已删除',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_trade_no` (`trade_no`),
  KEY `idx_goods_id` (`goods_id`),
  KEY `idx_seller_id` (`seller_id`),
  KEY `idx_buyer_id` (`buyer_id`),
  KEY `idx_trade_status` (`trade_status`),
  KEY `idx_tenant_merchant` (`tenant_id`, `merchant_id`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='仓库交易记录';

-- ================================================================
-- 13. 团队数据导出记录
-- ================================================================
CREATE TABLE IF NOT EXISTS `member_team_export` (
  `id` bigint unsigned NOT NULL COMMENT 'ID（Snowflake）',
  `user_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '目标会员|ref:member_user.nickname|search:select',
  `team_member_count` int unsigned NOT NULL DEFAULT '0' COMMENT '团队成员数',
  `export_type` tinyint NOT NULL DEFAULT '1' COMMENT '导出类型:1=手动导出,2=自动升级导出|search:select',
  `file_url` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '导出文件地址',
  `file_size` bigint unsigned NOT NULL DEFAULT '0' COMMENT '文件大小（字节）',
  `deploy_status` tinyint NOT NULL DEFAULT '0' COMMENT '部署状态:0=未部署,1=部署中,2=已部署,3=部署失败|search:select',
  `deploy_domain` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '部署域名|search:like',
  `deployed_at` datetime DEFAULT NULL COMMENT '部署完成时间',
  `remark` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '备注|search:off',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态:0=关闭,1=开启|search:select',
  `tenant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '租户',
  `merchant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '商户',
  `created_by` bigint unsigned NOT NULL DEFAULT '0' COMMENT '创建人ID',
  `dept_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '所属部门ID',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间，非 NULL 表示已删除',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_deploy_status` (`deploy_status`),
  KEY `idx_tenant_merchant` (`tenant_id`, `merchant_id`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='团队数据导出';

-- ================================================================
-- 菜单种子数据（member 一级目录 + 13 个二级页面 + 操作权限按钮）
-- ================================================================

-- 一级目录：会员管理
INSERT INTO `system_menu` (`id`, `parent_id`, `title`, `type`, `path`, `component`, `permission`, `icon`, `sort`, `is_show`, `is_cache`, `link_url`, `status`, `created_by`, `dept_id`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (1000000000000002000, 0, '会员管理', 1, '/member', NULL, '', 'TeamOutlined', 30, 1, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL)
ON DUPLICATE KEY UPDATE `title` = VALUES(`title`), `icon` = VALUES(`icon`), `sort` = VALUES(`sort`), `status` = VALUES(`status`), `updated_at` = NOW(), `deleted_at` = NULL;

-- 二级页面
INSERT INTO `system_menu` (`id`, `parent_id`, `title`, `type`, `path`, `component`, `permission`, `icon`, `sort`, `is_show`, `is_cache`, `link_url`, `status`, `created_by`, `dept_id`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (1000000000000002010, 1000000000000002000, '会员列表',   2, '/member/user',              'member/user/index',              'member:user:list',              'UserOutlined',         10, 1, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000002020, 1000000000000002000, '等级配置',   2, '/member/level',             'member/level/index',             'member:level:list',             'CrownOutlined',        20, 1, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000002030, 1000000000000002000, '等级日志',   2, '/member/level-log',         'member/level_log/index',         'member:level_log:list',         'HistoryOutlined',      25, 1, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000002040, 1000000000000002000, '钱包管理',   2, '/member/wallet',            'member/wallet/index',            'member:wallet:list',            'WalletOutlined',       30, 1, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000002050, 1000000000000002000, '钱包流水',   2, '/member/wallet-log',        'member/wallet_log/index',        'member:wallet_log:list',        'TransactionOutlined',  35, 1, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000002060, 1000000000000002000, '换绑日志',   2, '/member/rebind-log',        'member/rebind_log/index',        'member:rebind_log:list',        'SwapOutlined',         40, 1, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000002070, 1000000000000002000, '商城分类',   2, '/member/shop-category',     'member/shop_category/index',     'member:shop_category:list',     'AppstoreOutlined',     50, 1, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000002080, 1000000000000002000, '商城商品',   2, '/member/shop-goods',        'member/shop_goods/index',        'member:shop_goods:list',        'ShoppingOutlined',     55, 1, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000002090, 1000000000000002000, '商城订单',   2, '/member/shop-order',        'member/shop_order/index',        'member:shop_order:list',        'ShoppingCartOutlined', 60, 1, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000002100, 1000000000000002000, '仓库商品',   2, '/member/warehouse-goods',   'member/warehouse_goods/index',   'member:warehouse_goods:list',   'InboxOutlined',        70, 1, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000002110, 1000000000000002000, '挂卖记录',   2, '/member/warehouse-listing', 'member/warehouse_listing/index', 'member:warehouse_listing:list', 'TagsOutlined',         75, 1, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000002120, 1000000000000002000, '交易记录',   2, '/member/warehouse-trade',   'member/warehouse_trade/index',   'member:warehouse_trade:list',   'SwapOutlined',         80, 1, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000002130, 1000000000000002000, '团队导出',   2, '/member/team-export',       'member/team_export/index',       'member:team_export:list',       'ExportOutlined',       90, 1, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL)
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

-- 操作权限按钮（每个模块：新增/修改/删除/批量删除/查看/导出/批量编辑）
INSERT INTO `system_menu` (`id`, `parent_id`, `title`, `type`, `path`, `component`, `permission`, `icon`, `sort`, `is_show`, `is_cache`, `link_url`, `status`, `created_by`, `dept_id`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  -- 会员列表
  (1000000000000002011, 1000000000000002010, '会员新增',       3, NULL, NULL, 'member:user:create',       '', 1, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000002012, 1000000000000002010, '会员修改',       3, NULL, NULL, 'member:user:update',       '', 2, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000002013, 1000000000000002010, '会员删除',       3, NULL, NULL, 'member:user:delete',       '', 3, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000002014, 1000000000000002010, '会员批量删除',   3, NULL, NULL, 'member:user:batch-delete', '', 4, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000002015, 1000000000000002010, '会员查看',       3, NULL, NULL, 'member:user:detail',       '', 5, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000002016, 1000000000000002010, '会员导出',       3, NULL, NULL, 'member:user:export',       '', 6, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000002017, 1000000000000002010, '会员批量编辑',   3, NULL, NULL, 'member:user:batch-update', '', 8, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  -- 等级配置
  (1000000000000002021, 1000000000000002020, '等级新增',       3, NULL, NULL, 'member:level:create',       '', 1, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000002022, 1000000000000002020, '等级修改',       3, NULL, NULL, 'member:level:update',       '', 2, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000002023, 1000000000000002020, '等级删除',       3, NULL, NULL, 'member:level:delete',       '', 3, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000002024, 1000000000000002020, '等级批量删除',   3, NULL, NULL, 'member:level:batch-delete', '', 4, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000002025, 1000000000000002020, '等级查看',       3, NULL, NULL, 'member:level:detail',       '', 5, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000002026, 1000000000000002020, '等级导出',       3, NULL, NULL, 'member:level:export',       '', 6, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  -- 等级日志
  (1000000000000002031, 1000000000000002030, '等级日志查看',   3, NULL, NULL, 'member:level_log:detail',   '', 5, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000002036, 1000000000000002030, '等级日志导出',   3, NULL, NULL, 'member:level_log:export',   '', 6, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  -- 钱包管理
  (1000000000000002041, 1000000000000002040, '钱包查看',       3, NULL, NULL, 'member:wallet:detail',       '', 5, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000002042, 1000000000000002040, '钱包调整',       3, NULL, NULL, 'member:wallet:update',       '', 2, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000002046, 1000000000000002040, '钱包导出',       3, NULL, NULL, 'member:wallet:export',       '', 6, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  -- 钱包流水
  (1000000000000002051, 1000000000000002050, '钱包流水查看',   3, NULL, NULL, 'member:wallet_log:detail',   '', 5, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000002056, 1000000000000002050, '钱包流水导出',   3, NULL, NULL, 'member:wallet_log:export',   '', 6, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  -- 换绑日志
  (1000000000000002061, 1000000000000002060, '换绑日志查看',   3, NULL, NULL, 'member:rebind_log:detail',   '', 5, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000002066, 1000000000000002060, '换绑日志导出',   3, NULL, NULL, 'member:rebind_log:export',   '', 6, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  -- 商城分类
  (1000000000000002071, 1000000000000002070, '商城分类新增',   3, NULL, NULL, 'member:shop_category:create',       '', 1, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000002072, 1000000000000002070, '商城分类修改',   3, NULL, NULL, 'member:shop_category:update',       '', 2, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000002073, 1000000000000002070, '商城分类删除',   3, NULL, NULL, 'member:shop_category:delete',       '', 3, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000002074, 1000000000000002070, '商城分类批量删除', 3, NULL, NULL, 'member:shop_category:batch-delete', '', 4, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000002075, 1000000000000002070, '商城分类查看',   3, NULL, NULL, 'member:shop_category:detail',       '', 5, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000002076, 1000000000000002070, '商城分类导出',   3, NULL, NULL, 'member:shop_category:export',       '', 6, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  -- 商城商品
  (1000000000000002081, 1000000000000002080, '商城商品新增',   3, NULL, NULL, 'member:shop_goods:create',       '', 1, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000002082, 1000000000000002080, '商城商品修改',   3, NULL, NULL, 'member:shop_goods:update',       '', 2, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000002083, 1000000000000002080, '商城商品删除',   3, NULL, NULL, 'member:shop_goods:delete',       '', 3, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000002084, 1000000000000002080, '商城商品批量删除', 3, NULL, NULL, 'member:shop_goods:batch-delete', '', 4, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000002085, 1000000000000002080, '商城商品查看',   3, NULL, NULL, 'member:shop_goods:detail',       '', 5, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000002086, 1000000000000002080, '商城商品导出',   3, NULL, NULL, 'member:shop_goods:export',       '', 6, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  -- 商城订单
  (1000000000000002091, 1000000000000002090, '商城订单查看',   3, NULL, NULL, 'member:shop_order:detail',       '', 5, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000002096, 1000000000000002090, '商城订单导出',   3, NULL, NULL, 'member:shop_order:export',       '', 6, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  -- 仓库商品
  (1000000000000002101, 1000000000000002100, '仓库商品新增',   3, NULL, NULL, 'member:warehouse_goods:create',       '', 1, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000002102, 1000000000000002100, '仓库商品修改',   3, NULL, NULL, 'member:warehouse_goods:update',       '', 2, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000002103, 1000000000000002100, '仓库商品删除',   3, NULL, NULL, 'member:warehouse_goods:delete',       '', 3, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000002104, 1000000000000002100, '仓库商品批量删除', 3, NULL, NULL, 'member:warehouse_goods:batch-delete', '', 4, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000002105, 1000000000000002100, '仓库商品查看',   3, NULL, NULL, 'member:warehouse_goods:detail',       '', 5, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000002106, 1000000000000002100, '仓库商品导出',   3, NULL, NULL, 'member:warehouse_goods:export',       '', 6, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  -- 挂卖记录
  (1000000000000002111, 1000000000000002110, '挂卖记录查看',   3, NULL, NULL, 'member:warehouse_listing:detail',   '', 5, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000002116, 1000000000000002110, '挂卖记录导出',   3, NULL, NULL, 'member:warehouse_listing:export',   '', 6, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  -- 交易记录
  (1000000000000002121, 1000000000000002120, '交易记录查看',   3, NULL, NULL, 'member:warehouse_trade:detail',     '', 5, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000002126, 1000000000000002120, '交易记录导出',   3, NULL, NULL, 'member:warehouse_trade:export',     '', 6, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  -- 团队导出
  (1000000000000002131, 1000000000000002130, '团队导出新增',   3, NULL, NULL, 'member:team_export:create',         '', 1, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000002135, 1000000000000002130, '团队导出查看',   3, NULL, NULL, 'member:team_export:detail',         '', 5, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL),
  (1000000000000002136, 1000000000000002130, '团队导出导出',   3, NULL, NULL, 'member:team_export:export',         '', 6, 0, 0, NULL, 1, 0, 0, NOW(), NOW(), NULL)
ON DUPLICATE KEY UPDATE
  `parent_id` = VALUES(`parent_id`),
  `title` = VALUES(`title`),
  `permission` = VALUES(`permission`),
  `sort` = VALUES(`sort`),
  `status` = VALUES(`status`),
  `updated_at` = NOW(),
  `deleted_at` = NULL;

-- 超级管理员角色绑定所有 member 菜单
INSERT INTO `system_role_menu` (`role_id`, `menu_id`)
SELECT 1000000000000000002, `id`
FROM `system_menu`
WHERE (`path` LIKE '/member%' OR `permission` LIKE 'member:%') AND `deleted_at` IS NULL
ON DUPLICATE KEY UPDATE `menu_id` = VALUES(`menu_id`);

SET FOREIGN_KEY_CHECKS = 1;
