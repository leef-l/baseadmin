-- 上传管理与系统管理补齐 SaaS 归属字段

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

SET @schema_name := DATABASE();

SET @sql := (SELECT IF(COUNT(*) = 0,
  'ALTER TABLE `upload_config` ADD COLUMN `tenant_id` bigint unsigned NOT NULL DEFAULT ''0'' COMMENT ''租户'' AFTER `dept_id`',
  'DO 0')
FROM information_schema.columns
WHERE table_schema = @schema_name AND table_name = 'upload_config' AND column_name = 'tenant_id');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(COUNT(*) = 0,
  'ALTER TABLE `upload_config` ADD COLUMN `merchant_id` bigint unsigned NOT NULL DEFAULT ''0'' COMMENT ''商户'' AFTER `tenant_id`',
  'DO 0')
FROM information_schema.columns
WHERE table_schema = @schema_name AND table_name = 'upload_config' AND column_name = 'merchant_id');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(COUNT(*) = 0,
  'ALTER TABLE `upload_config` ADD KEY `idx_tenant_merchant` (`tenant_id`, `merchant_id`)',
  'DO 0')
FROM information_schema.statistics
WHERE table_schema = @schema_name AND table_name = 'upload_config' AND index_name = 'idx_tenant_merchant');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(COUNT(*) = 0,
  'ALTER TABLE `upload_dir` ADD COLUMN `tenant_id` bigint unsigned NOT NULL DEFAULT ''0'' COMMENT ''租户'' AFTER `dept_id`',
  'DO 0')
FROM information_schema.columns
WHERE table_schema = @schema_name AND table_name = 'upload_dir' AND column_name = 'tenant_id');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(COUNT(*) = 0,
  'ALTER TABLE `upload_dir` ADD COLUMN `merchant_id` bigint unsigned NOT NULL DEFAULT ''0'' COMMENT ''商户'' AFTER `tenant_id`',
  'DO 0')
FROM information_schema.columns
WHERE table_schema = @schema_name AND table_name = 'upload_dir' AND column_name = 'merchant_id');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(COUNT(*) = 0,
  'ALTER TABLE `upload_dir` ADD KEY `idx_tenant_merchant` (`tenant_id`, `merchant_id`)',
  'DO 0')
FROM information_schema.statistics
WHERE table_schema = @schema_name AND table_name = 'upload_dir' AND index_name = 'idx_tenant_merchant');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(COUNT(*) = 0,
  'ALTER TABLE `upload_dir_rule` ADD COLUMN `tenant_id` bigint unsigned NOT NULL DEFAULT ''0'' COMMENT ''租户'' AFTER `dept_id`',
  'DO 0')
FROM information_schema.columns
WHERE table_schema = @schema_name AND table_name = 'upload_dir_rule' AND column_name = 'tenant_id');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(COUNT(*) = 0,
  'ALTER TABLE `upload_dir_rule` ADD COLUMN `merchant_id` bigint unsigned NOT NULL DEFAULT ''0'' COMMENT ''商户'' AFTER `tenant_id`',
  'DO 0')
FROM information_schema.columns
WHERE table_schema = @schema_name AND table_name = 'upload_dir_rule' AND column_name = 'merchant_id');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(COUNT(*) = 0,
  'ALTER TABLE `upload_dir_rule` ADD KEY `idx_tenant_merchant` (`tenant_id`, `merchant_id`)',
  'DO 0')
FROM information_schema.statistics
WHERE table_schema = @schema_name AND table_name = 'upload_dir_rule' AND index_name = 'idx_tenant_merchant');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(COUNT(*) = 0,
  'ALTER TABLE `upload_file` ADD COLUMN `tenant_id` bigint unsigned NOT NULL DEFAULT ''0'' COMMENT ''租户'' AFTER `dept_id`',
  'DO 0')
FROM information_schema.columns
WHERE table_schema = @schema_name AND table_name = 'upload_file' AND column_name = 'tenant_id');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(COUNT(*) = 0,
  'ALTER TABLE `upload_file` ADD COLUMN `merchant_id` bigint unsigned NOT NULL DEFAULT ''0'' COMMENT ''商户'' AFTER `tenant_id`',
  'DO 0')
FROM information_schema.columns
WHERE table_schema = @schema_name AND table_name = 'upload_file' AND column_name = 'merchant_id');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(COUNT(*) = 0,
  'ALTER TABLE `upload_file` ADD KEY `idx_tenant_merchant` (`tenant_id`, `merchant_id`)',
  'DO 0')
FROM information_schema.statistics
WHERE table_schema = @schema_name AND table_name = 'upload_file' AND index_name = 'idx_tenant_merchant');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(COUNT(*) = 0,
  'ALTER TABLE `system_tenant` ADD COLUMN `tenant_id` bigint unsigned NOT NULL DEFAULT ''0'' COMMENT ''租户'' AFTER `dept_id`',
  'DO 0')
FROM information_schema.columns
WHERE table_schema = @schema_name AND table_name = 'system_tenant' AND column_name = 'tenant_id');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(COUNT(*) = 0,
  'ALTER TABLE `system_tenant` ADD COLUMN `merchant_id` bigint unsigned NOT NULL DEFAULT ''0'' COMMENT ''商户'' AFTER `tenant_id`',
  'DO 0')
FROM information_schema.columns
WHERE table_schema = @schema_name AND table_name = 'system_tenant' AND column_name = 'merchant_id');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(COUNT(*) = 0,
  'ALTER TABLE `system_tenant` ADD KEY `idx_tenant_merchant` (`tenant_id`, `merchant_id`)',
  'DO 0')
FROM information_schema.statistics
WHERE table_schema = @schema_name AND table_name = 'system_tenant' AND index_name = 'idx_tenant_merchant');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(COUNT(*) = 0,
  'ALTER TABLE `system_merchant` ADD COLUMN `merchant_id` bigint unsigned NOT NULL DEFAULT ''0'' COMMENT ''商户'' AFTER `tenant_id`',
  'DO 0')
FROM information_schema.columns
WHERE table_schema = @schema_name AND table_name = 'system_merchant' AND column_name = 'merchant_id');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(COUNT(*) = 0,
  'ALTER TABLE `system_merchant` ADD KEY `idx_tenant_merchant` (`tenant_id`, `merchant_id`)',
  'DO 0')
FROM information_schema.statistics
WHERE table_schema = @schema_name AND table_name = 'system_merchant' AND index_name = 'idx_tenant_merchant');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

UPDATE `system_merchant`
SET `merchant_id` = `id`
WHERE `merchant_id` = 0;

SET @sql := (SELECT IF(COUNT(*) = 0,
  'ALTER TABLE `system_daemon` ADD COLUMN `tenant_id` bigint unsigned NOT NULL DEFAULT ''0'' COMMENT ''租户'' AFTER `dept_id`',
  'DO 0')
FROM information_schema.columns
WHERE table_schema = @schema_name AND table_name = 'system_daemon' AND column_name = 'tenant_id');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(COUNT(*) = 0,
  'ALTER TABLE `system_daemon` ADD COLUMN `merchant_id` bigint unsigned NOT NULL DEFAULT ''0'' COMMENT ''商户'' AFTER `tenant_id`',
  'DO 0')
FROM information_schema.columns
WHERE table_schema = @schema_name AND table_name = 'system_daemon' AND column_name = 'merchant_id');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(COUNT(*) = 0,
  'ALTER TABLE `system_daemon` ADD KEY `idx_tenant_merchant` (`tenant_id`, `merchant_id`)',
  'DO 0')
FROM information_schema.statistics
WHERE table_schema = @schema_name AND table_name = 'system_daemon' AND index_name = 'idx_tenant_merchant');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(COUNT(*) = 0,
  'ALTER TABLE `system_menu` ADD COLUMN `tenant_id` bigint unsigned NOT NULL DEFAULT ''0'' COMMENT ''租户'' AFTER `dept_id`',
  'DO 0')
FROM information_schema.columns
WHERE table_schema = @schema_name AND table_name = 'system_menu' AND column_name = 'tenant_id');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(COUNT(*) = 0,
  'ALTER TABLE `system_menu` ADD COLUMN `merchant_id` bigint unsigned NOT NULL DEFAULT ''0'' COMMENT ''商户'' AFTER `tenant_id`',
  'DO 0')
FROM information_schema.columns
WHERE table_schema = @schema_name AND table_name = 'system_menu' AND column_name = 'merchant_id');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(COUNT(*) = 0,
  'ALTER TABLE `system_menu` ADD KEY `idx_tenant_merchant` (`tenant_id`, `merchant_id`)',
  'DO 0')
FROM information_schema.statistics
WHERE table_schema = @schema_name AND table_name = 'system_menu' AND index_name = 'idx_tenant_merchant');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET FOREIGN_KEY_CHECKS = 1;
