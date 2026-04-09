INSERT INTO `system_menu` (
  `id`, `parent_id`, `title`, `type`, `path`, `component`, `permission`, `icon`,
  `sort`, `is_show`, `is_cache`, `link_url`, `status`, `created_by`, `dept_id`,
  `created_at`, `updated_at`, `deleted_at`
) VALUES
  (1000000000000000060, 0, '仪表盘', 1, '/dashboard', NULL, '', 'DashboardOutlined', 0, 1, 0, NULL, 1, 0, 1000000000000000001, '2026-03-30 21:20:22', '2026-03-30 21:20:22', NULL),
  (1000000000000000061, 1000000000000000060, '分析页', 2, '/analytics', 'dashboard/analytics/index', '', 'AreaChartOutlined', 1, 1, 1, NULL, 1, 0, 1000000000000000001, '2026-03-30 21:20:22', '2026-03-30 21:20:22', NULL),
  (1000000000000000062, 1000000000000000060, '工作台', 2, '/workspace', 'dashboard/workspace/index', '', 'DesktopOutlined', 2, 1, 0, NULL, 1, 0, 1000000000000000001, '2026-03-30 21:20:22', '2026-03-30 21:20:22', NULL);

INSERT INTO `system_role_menu` (`role_id`, `menu_id`) VALUES
  (1000000000000000002, 1000000000000000060),
  (1000000000000000002, 1000000000000000061),
  (1000000000000000002, 1000000000000000062);
