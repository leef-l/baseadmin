-- 回滚上传管理与系统管理 SaaS 归属字段

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

ALTER TABLE `system_menu`
  DROP KEY `idx_tenant_merchant`,
  DROP COLUMN `merchant_id`,
  DROP COLUMN `tenant_id`;

ALTER TABLE `system_daemon`
  DROP KEY `idx_tenant_merchant`,
  DROP COLUMN `merchant_id`,
  DROP COLUMN `tenant_id`;

ALTER TABLE `system_merchant`
  DROP KEY `idx_tenant_merchant`,
  DROP COLUMN `merchant_id`;

ALTER TABLE `system_tenant`
  DROP KEY `idx_tenant_merchant`,
  DROP COLUMN `merchant_id`,
  DROP COLUMN `tenant_id`;

ALTER TABLE `upload_file`
  DROP KEY `idx_tenant_merchant`,
  DROP COLUMN `merchant_id`,
  DROP COLUMN `tenant_id`;

ALTER TABLE `upload_dir_rule`
  DROP KEY `idx_tenant_merchant`,
  DROP COLUMN `merchant_id`,
  DROP COLUMN `tenant_id`;

ALTER TABLE `upload_dir`
  DROP KEY `idx_tenant_merchant`,
  DROP COLUMN `merchant_id`,
  DROP COLUMN `tenant_id`;

ALTER TABLE `upload_config`
  DROP KEY `idx_tenant_merchant`,
  DROP COLUMN `merchant_id`,
  DROP COLUMN `tenant_id`;

SET FOREIGN_KEY_CHECKS = 1;
