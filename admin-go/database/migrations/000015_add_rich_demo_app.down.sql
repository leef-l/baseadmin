-- 回滚丰富 Demo 体验应用。

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

DELETE FROM `system_role_menu`
WHERE `menu_id` IN (
  SELECT `id` FROM `system_menu`
  WHERE `path` LIKE '/demo%' OR `permission` LIKE 'demo:%'
);

DELETE FROM `system_menu`
WHERE `path` LIKE '/demo%' OR `permission` LIKE 'demo:%';

DROP TABLE IF EXISTS `demo_audit_log`;
DROP TABLE IF EXISTS `demo_appointment`;
DROP TABLE IF EXISTS `demo_survey`;
DROP TABLE IF EXISTS `demo_contract`;
DROP TABLE IF EXISTS `demo_work_order`;
DROP TABLE IF EXISTS `demo_order`;
DROP TABLE IF EXISTS `demo_campaign`;
DROP TABLE IF EXISTS `demo_product`;
DROP TABLE IF EXISTS `demo_customer`;
DROP TABLE IF EXISTS `demo_category`;

DELETE FROM `system_merchant`
WHERE `id` = 324500000000000101 AND `code` = 'demo-merchant';

DELETE FROM `system_tenant`
WHERE `id` = 324500000000000001 AND `code` = 'demo-tenant';

SET FOREIGN_KEY_CHECKS = 1;
