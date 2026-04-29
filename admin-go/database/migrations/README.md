# Database Migrations

铁律：

- 以后所有数据库结构变更、初始化数据变更、默认菜单/默认账号/默认上传配置变更，必须新增 `golang-migrate` 迁移文件。
- 如果迁移里新增默认菜单种子，还必须同步补齐默认角色授权种子，至少保证基线超级管理员角色可见。
- 不再把 `docker/mysql/init.sql`、`admin-go/codegen/sql/init.sql` 当作数据库真源。
- Docker 容器启动前必须先执行 `migrate up`，再启动应用。
- 默认由 `system` 容器负责自动执行迁移；`upload` 容器默认不重复抢跑迁移。

常用命令：

```bash
cd admin-go
go run ./cmd/migrate up
go run ./cmd/migrate version
go run ./cmd/migrate create add_system_logs
```

默认迁移目录：`database/migrations`

## 当前迁移主线

| 版本 | 文件前缀 | 说明 |
|------|----------|------|
| `000001` | `baseline_system_upload` | system / upload 基线表、默认账号、默认角色、默认菜单、默认上传配置 |
| `000002` | `fill_missing_menu_icons` | 补齐菜单图标 |
| `000003` | `remove_dashboard_seed` | 移除 dashboard 种子，保持当前精简后台范围 |
| `000004` | `add_batch_delete_menu_seed` | 补批量删除按钮权限 |
| `000005` | `add_upload_dir_rule_file_type` | 上传目录规则增加文件类型匹配 |
| `000006` | `add_upload_dir_rule_storage_types` | 上传目录规则增加存储类型匹配 |
| `000007` | `add_upload_dir_rule_keep_name_and_expand_matchers` | 上传目录规则补保留原文件名和匹配条件 |
| `000008` | `add_upload_dir_keep_name` | 上传目录补保留原文件名配置 |
| `000009` | `repair_upload_dir_rule_storage_types` | 修复上传目录规则存储类型数据 |
| `000010` | `add_saas_tenant_merchant` | 新增租户、商户表，为用户/角色/部门补租户和商户字段，补租户/商户菜单 |
| `000011` | `add_saas_domain_binding` | 新增域名绑定表，补域名管理和应用 Nginx 权限 |
| `000012` | `add_domain_ssl_permission` | 补域名 SSL 申请按钮权限 |
| `000013` | `add_system_daemon` | 新增守护进程表，补守护进程菜单和按钮权限 |
| `000014` | `add_scope_fields_to_upload_and_system` | 为上传管理和系统管理补齐租户、商户归属字段 |

## SaaS 迁移约定

- 所有交给 codegen 生成 CRUD 的业务表必须同时有 `tenant_id`、`merchant_id`、`created_by` 和 `dept_id`。
- `tenant_id`、`merchant_id` 用于 SaaS 运营主体隔离；`created_by`、`dept_id` 用于部门数据权限注入和过滤。
- `tenant_id=0`、`merchant_id=0` 表示平台级数据。
- `tenant_id>0`、`merchant_id=0` 表示租户级数据。
- `tenant_id>0`、`merchant_id>0` 表示商户级数据。
- 如果迁移新增菜单，必须同步插入 `system_role_menu`，至少保证基线超级管理员角色可见。

## 字段备注要求

新增字段时必须认真写 `COMMENT`，因为 codegen 会直接读取字段备注生成前端标签、枚举、Tooltip、搜索和外键关联。

示例：

```sql
`status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态:0=关闭,1=开启',
`amount` bigint unsigned NOT NULL DEFAULT '0' COMMENT '金额（分）',
`user_id` bigint unsigned NOT NULL COMMENT '用户|ref:system_users.username|search:select'
```

详细规则见 `admin-go/codegen/docs/字段备注与生成规则.md`。
