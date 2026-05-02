-- ================================================================
-- 000020: 后台菜单：会员业务配置 / 合同模板 / 合同记录
--
-- 注意：本迁移作用于 system 服务所连数据库（baseadmin 站点库），
-- 与其它 000xxx 迁移（funddisk 站点库）数据库不同，需在 baseadmin 库下执行。
-- 其它 000018/000019 迁移作用于 funddisk 业务库，互不影响。
-- ================================================================

SET NAMES utf8mb4;

INSERT IGNORE INTO `system_menu`
  (`id`, `parent_id`, `title`, `type`, `path`, `component`, `permission`, `icon`, `sort`, `is_show`, `is_cache`, `status`, `created_at`, `updated_at`)
VALUES
  (1000000000000002140, 1000000000000002000, '业务配置', 2, '/member/biz-config', 'member/biz_config/index', 'member:biz-config:edit', 'lucide:settings',  5, 1, 0, 1, NOW(), NOW()),
  (1000000000000002150, 1000000000000002000, '合同模板', 2, '/member/contract-template', 'member/contract_template/index', 'member:contract-template:list', 'lucide:file-cog', 95, 1, 0, 1, NOW(), NOW()),
  (1000000000000002160, 1000000000000002000, '合同记录', 2, '/member/contract', 'member/contract/index', 'member:contract:list', 'lucide:file-signature', 96, 1, 0, 1, NOW(), NOW());

INSERT IGNORE INTO `system_menu`
  (`id`, `parent_id`, `title`, `type`, `permission`, `sort`, `is_show`, `status`, `created_at`, `updated_at`)
VALUES
  (1000000000000002141, 1000000000000002140, '保存配置', 3, 'member:biz-config:save', 1, 1, 1, NOW(), NOW()),
  (1000000000000002151, 1000000000000002150, '新增模板', 3, 'member:contract-template:create', 1, 1, 1, NOW(), NOW()),
  (1000000000000002152, 1000000000000002150, '修改模板', 3, 'member:contract-template:update', 2, 1, 1, NOW(), NOW()),
  (1000000000000002153, 1000000000000002150, '删除模板', 3, 'member:contract-template:delete', 3, 1, 1, NOW(), NOW()),
  (1000000000000002154, 1000000000000002150, '查看模板', 3, 'member:contract-template:detail', 4, 1, 1, NOW(), NOW()),
  (1000000000000002161, 1000000000000002160, '合同列表', 3, 'member:contract:list', 1, 1, 1, NOW(), NOW()),
  (1000000000000002162, 1000000000000002160, '下载合同', 3, 'member:contract:download', 2, 1, 1, NOW(), NOW());
