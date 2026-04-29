-- 回滚会员分销系统

SET FOREIGN_KEY_CHECKS = 0;

DELETE FROM `system_role_menu` WHERE `menu_id` IN (
  SELECT `id` FROM `system_menu` WHERE (`path` LIKE '/member%' OR `permission` LIKE 'member:%') AND `deleted_at` IS NULL
);

DELETE FROM `system_menu` WHERE `path` LIKE '/member%' OR `permission` LIKE 'member:%';

DROP TABLE IF EXISTS `member_team_export`;
DROP TABLE IF EXISTS `member_warehouse_trade`;
DROP TABLE IF EXISTS `member_warehouse_listing`;
DROP TABLE IF EXISTS `member_warehouse_goods`;
DROP TABLE IF EXISTS `member_shop_order`;
DROP TABLE IF EXISTS `member_shop_goods`;
DROP TABLE IF EXISTS `member_shop_category`;
DROP TABLE IF EXISTS `member_rebind_log`;
DROP TABLE IF EXISTS `member_wallet_log`;
DROP TABLE IF EXISTS `member_wallet`;
DROP TABLE IF EXISTS `member_level_log`;
DROP TABLE IF EXISTS `member_user`;
DROP TABLE IF EXISTS `member_level`;

SET FOREIGN_KEY_CHECKS = 1;
