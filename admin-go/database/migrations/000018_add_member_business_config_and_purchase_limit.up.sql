-- ================================================================
-- 000018: 业务配置表 + 会员/等级限购字段 + 等级数据重置
-- 业务背景：
--   1. 进货时间窗（10:00-10:30）/ 寄售时间窗（14:30~）/ 工作日设置 等运营参数全部走配置表
--   2. 自购阶梯返佣 / 自购 1% 返奖励 / 直推 0.4% 返推广 等比例配置
--   3. 每日限购按等级控制；会员表加每日已购计数，跨日 cron 重置
--   4. 重置 member_level 为 普通=1单 / 高级=2单 / 核心=5单
--
-- 幂等性：通过 INFORMATION_SCHEMA + PREPARE/EXECUTE 实现 ADD COLUMN IF NOT EXISTS。
-- 不使用 DELIMITER（兼容 mysql CLI 直接 source 加载）。
-- ================================================================

SET NAMES utf8mb4;

-- ----------------------------------------------------------------
-- 1. 会员业务配置（singleton，整表只一行 config_key='global'）
-- ----------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `member_business_config` (
  `id` bigint unsigned NOT NULL COMMENT '配置ID（Snowflake）',
  `config_key` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'global' COMMENT '配置键|search:eq',
  `payload` json NOT NULL COMMENT '业务配置JSON（进货时间窗/寄售时间窗/工作日/返佣比例等）|search:off',
  `remark` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '备注|search:off',
  `tenant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '租户',
  `merchant_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '商户',
  `created_by` bigint unsigned NOT NULL DEFAULT '0' COMMENT '创建人ID',
  `dept_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '所属部门ID',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间，非 NULL 表示已删除',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_config_key_tenant` (`config_key`, `tenant_id`, `merchant_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='会员业务配置';

-- 初始化默认配置（INSERT IGNORE，避免重跑覆盖运营改过的配置）
INSERT IGNORE INTO `member_business_config`
  (`id`, `config_key`, `payload`, `remark`, `created_at`, `updated_at`)
VALUES
  (
    1,
    'global',
    JSON_OBJECT(
      'purchase', JSON_OBJECT(
        'startTime', '10:00',
        'endTime',   '10:30',
        'allowedWeekdays', JSON_ARRAY(1,2,3,4,5)
      ),
      'consign', JSON_OBJECT(
        'startTime', '14:30',
        'endTime',   NULL
      ),
      'selfRebateTiers', JSON_ARRAY(
        JSON_OBJECT('nthOrder', 2, 'rewardYuan', 88),
        JSON_OBJECT('nthOrder', 3, 'rewardYuan', 188),
        JSON_OBJECT('nthOrder', 4, 'rewardYuan', 288)
      ),
      'selfTurnoverRewardRate', 1.0,
      'directPromoteRate',      0.4
    ),
    '默认配置：进货 10:00-10:30 工作日；寄售 14:30 起无截止；自购 2/3/4 单 88/188/288；自购 1%；直推 0.4%',
    NOW(), NOW()
  );

-- ----------------------------------------------------------------
-- 2. ALTER TABLE：用 PREPARE/EXECUTE 实现幂等 ADD COLUMN（不依赖 DELIMITER）
-- ----------------------------------------------------------------

-- member_level.daily_purchase_limit
SET @sql := IF(
  (SELECT COUNT(*) FROM information_schema.COLUMNS
     WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'member_level' AND COLUMN_NAME = 'daily_purchase_limit') = 0,
  'ALTER TABLE `member_level` ADD COLUMN `daily_purchase_limit` int unsigned NOT NULL DEFAULT 1 COMMENT ''该等级每日限购单数|search:eq'' AFTER `need_team_turnover`',
  'SELECT 1'
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- member_user.daily_purchase_limit
SET @sql := IF(
  (SELECT COUNT(*) FROM information_schema.COLUMNS
     WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'member_user' AND COLUMN_NAME = 'daily_purchase_limit') = 0,
  'ALTER TABLE `member_user` ADD COLUMN `daily_purchase_limit` int unsigned NOT NULL DEFAULT 1 COMMENT ''本会员每日限购单数（按等级初始化，可单独调整）|search:eq'' AFTER `is_qualified`',
  'SELECT 1'
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- member_user.today_purchase_count
SET @sql := IF(
  (SELECT COUNT(*) FROM information_schema.COLUMNS
     WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'member_user' AND COLUMN_NAME = 'today_purchase_count') = 0,
  'ALTER TABLE `member_user` ADD COLUMN `today_purchase_count` int unsigned NOT NULL DEFAULT 0 COMMENT ''今日已购单数|search:off'' AFTER `daily_purchase_limit`',
  'SELECT 1'
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- member_user.last_purchase_date
SET @sql := IF(
  (SELECT COUNT(*) FROM information_schema.COLUMNS
     WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'member_user' AND COLUMN_NAME = 'last_purchase_date') = 0,
  'ALTER TABLE `member_user` ADD COLUMN `last_purchase_date` date DEFAULT NULL COMMENT ''最近购买日期（跨日重置 today_purchase_count）|search:off'' AFTER `today_purchase_count`',
  'SELECT 1'
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- member_user.total_purchase_count
SET @sql := IF(
  (SELECT COUNT(*) FROM information_schema.COLUMNS
     WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'member_user' AND COLUMN_NAME = 'total_purchase_count') = 0,
  'ALTER TABLE `member_user` ADD COLUMN `total_purchase_count` int unsigned NOT NULL DEFAULT 0 COMMENT ''历史总购单数（用于阶梯返佣判断）|search:off'' AFTER `last_purchase_date`',
  'SELECT 1'
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- ----------------------------------------------------------------
-- 3. 等级数据：UPSERT，不强覆盖运营改过的 name/level_no
-- ----------------------------------------------------------------
INSERT INTO `member_level`
  (`id`, `name`, `level_no`, `icon`, `duration_days`, `need_active_count`, `need_team_turnover`,
   `daily_purchase_limit`, `is_top`, `auto_deploy`, `remark`, `sort`, `status`, `created_at`, `updated_at`)
VALUES
  (1, '普通会员',     1, '', 0,   0,        0, 1, 0, 0, '默认等级',                    1, 1, NOW(), NOW()),
  (2, '高级合伙人', 2, '', 0,  10,  100000000, 2, 0, 0, '团队 10 有效用户、100 万业绩', 2, 1, NOW(), NOW()),
  (3, '核心合伙人', 3, '', 0, 50, 1000000000, 5, 1, 0, '团队 50 有效用户、1000 万业绩，最高等级',  3, 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE
  `daily_purchase_limit` = VALUES(`daily_purchase_limit`),
  `updated_at` = NOW();

-- 给已有会员补 daily_purchase_limit（仅初始为 0 的，避免覆盖运营手动调整过的值）
UPDATE `member_user` u
LEFT JOIN `member_level` l ON l.id = u.level_id AND l.deleted_at IS NULL
SET u.daily_purchase_limit = COALESCE(l.daily_purchase_limit, 1)
WHERE u.deleted_at IS NULL AND u.daily_purchase_limit = 0;
