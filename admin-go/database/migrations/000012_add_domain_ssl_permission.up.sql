SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `system_menu` (`id`, `parent_id`, `title`, `type`, `path`, `component`, `permission`, `icon`, `sort`, `is_show`, `is_cache`, `link_url`, `status`, `created_by`, `dept_id`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (1000000000000000106, 1000000000000000100, '申请SSL证书', 3, NULL, NULL, 'system:domain:ssl', '', 6, 0, 0, NULL, 1, 0, 1000000000000000001, '2026-04-28 00:00:00', '2026-04-28 00:00:00', NULL)
ON DUPLICATE KEY UPDATE
  `title` = VALUES(`title`),
  `permission` = VALUES(`permission`),
  `sort` = VALUES(`sort`),
  `status` = VALUES(`status`),
  `deleted_at` = NULL;

INSERT INTO `system_role_menu` (`role_id`, `menu_id`)
VALUES (1000000000000000002, 1000000000000000106)
ON DUPLICATE KEY UPDATE `menu_id` = VALUES(`menu_id`);

SET FOREIGN_KEY_CHECKS = 1;
