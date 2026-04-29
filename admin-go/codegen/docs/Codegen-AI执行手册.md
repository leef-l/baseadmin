# Codegen AI 执行手册

本文给 AI/代理维护 codegen 时使用。业务字段备注规则见 [字段备注与生成规则.md](./字段备注与生成规则.md)。

## 先读顺序

1. `admin-go/codegen/README.md`
2. `admin-go/codegen/docs/字段备注与生成规则.md`
3. 本文
4. 需要改行为时，再读 `parser/`、`templates/`、`generator/` 对应文件

## 运行约束

当前服务器资源有限，所有 codegen 相关命令都必须走受限入口：

```bash
cd admin-go/codegen
../../scripts/wait-for-cpu-idle.sh
../../scripts/run-go-task-with-limits.sh go run verify_codegen.go
../../scripts/run-go-task-with-limits.sh go run . --table system_dept --dry-run
```

禁止事项：

- 不要裸跑 `pnpm`、`npm`、`npx`。
- 不要裸跑 Docker。
- 当前协作要求下不要跑 `go test`。如果以后 CI 或本地策略允许，也必须通过 `../../scripts/run-go-task-with-limits.sh` 执行。
- 不要手工改 `dao/`、`model/do/`、`model/entity/` 这类由 `gf gen dao` 生成的文件。

需要前端命令时，只能走根目录脚本：

```bash
./scripts/run-node-task-with-limits.sh <命令>
```

## 真源边界

| 内容 | 真源 |
|------|------|
| 数据库结构 | `admin-go/database/migrations/` |
| codegen 示例/离线验证 SQL | `admin-go/codegen/sql/*.sql` |
| 组件白名单 | `admin-go/codegen/parser/component_registry.go` |
| 字段命名启发式 | `admin-go/codegen/parser/rule_registry.go` |
| 字段备注解析 | `admin-go/codegen/parser/comment_parser.go` |
| 搜索派生规则 | `admin-go/codegen/parser/parser.go` |
| 后端 CRUD 生成 | `admin-go/codegen/templates/backend/*.tpl` |
| 前端页面生成 | `admin-go/codegen/templates/frontend/*.tpl` |
| 菜单写入 | `admin-go/codegen/generator/menu/` |

`admin-go/codegen/sql/*.sql` 不是线上迁移真源，不能用它替代 migrations。

## 排查路径

字段没有按预期生成时，按这个顺序排查：

1. 看 migration 里的字段名、类型、`NULL/DEFAULT`、`COMMENT`。
2. 看 `comment_parser.go` 是否能解析对应备注语法。
3. 看 `buildFieldMeta()` 是否正确设置 `FieldMeta`。
4. 看 `MapComponent()` 是否把字段映射到预期组件。
5. 看 `applySearchMeta()` 和 `finalizeSearchMeta()` 是否改变搜索行为。
6. 看前后端模板是否消费了对应 `FieldMeta` / `TableMeta` 字段。
7. 用受限 Go 入口执行 `go run verify_codegen.go` 或 `go run . --dry-run --manifest-out ...` 验证输出。

表级功能没有按预期生成时，优先看 `FinalizeTemplateMeta()`：

- `HasParentID` 控制树形接口和树形页面。
- `HasTenantScope` 控制租户权限守卫；codegen CRUD 表必须同时具备 `tenant_id` 和 `merchant_id`。
- `HasCreatedBy` / `HasDeptID` 控制部门数据权限；codegen CRUD 表必须同时具备 `created_by` 和 `dept_id`。
- `tenant_id`、`merchant_id`、`created_by`、`dept_id` 缺任意一个都会在解析阶段直接失败。
- `HasDict` 控制字典 API 导入。
- `HasForeignKey` 控制关联 API 导入和回显。
- `HasKeywordSearch` 控制全局关键词搜索。

## 字段备注重点

字段备注必须写到 migration 里。最小合格标准：

```sql
`name` varchar(100) NOT NULL COMMENT '名称|search:like|keyword:on'
`status` tinyint unsigned NOT NULL DEFAULT 1 COMMENT '状态:0=关闭,1=开启|search:select'
`amount` bigint unsigned NOT NULL DEFAULT 0 COMMENT '金额（分）'
`user_id` bigint unsigned NOT NULL COMMENT '用户|ref:system_users.username|search:select'
```

备注会影响：

- 前端表单 label、列表列头、详情字段和导出表头。
- Tooltip 文案。
- 静态枚举常量和前端 options。
- 动态字典加载。
- 搜索控件、搜索操作符和关键词搜索。
- 外键关联表、显示字段、前端 API 导入和后端批量回填。

不要只写“状态”“类型”“金额”这类模糊备注。状态要写值域，金额要写单位，外键要写业务对象，跨应用外键要写 `ref:`。

## SaaS / 数据权限契约

租户隔离不是靠前端隐藏字段实现，后端必须生成守卫。

触发条件和硬约束：

- 所有 codegen CRUD 表必须同时有 `tenant_id` 和 `merchant_id`：生成租户/商户权限。
- 所有 codegen CRUD 表必须同时有 `created_by` 和 `dept_id`：生成部门数据权限。
- 上面四个字段缺任意一个都不允许继续生成。

生成模板必须保持这些调用：

- Create：`ApplyTenantScopeToWrite`、`EnsureTenantMerchantAccessible`
- Update：写入归属覆盖、目标归属校验、行级归属校验
- Delete / BatchDelete：行级归属校验
- Detail：行级归属校验
- List / Tree / Export：`ApplyTenantScopeToModel`
- BatchUpdate：批量行级归属校验
- Import：写入归属覆盖和归属合法性校验
- 前端 List：`tenant_id` / `merchant_id` 的列表列和搜索筛选必须包在平台超级管理员判断里。

部门数据权限：

- Create / Import 注入 `created_by`、`dept_id`。
- List / Tree / Export 调用 `ApplyDataScope`。
- `tenant_id`、`merchant_id`、`created_by`、`dept_id` 都是生成前置条件，不再是可选触发字段。

如果改 CRUD 模板，必须逐项确认这些入口没有漏掉。

## 修改模板的注意点

后端模板：

- 写库优先使用 `do.{DaoName}`。
- 不手动维护 `created_at`、`updated_at`、`deleted_at`。
- 软删除统一走 `Delete()`。
- 租户、商户、部门、创建人字段由中间件或上下文注入，不要信任普通用户输入。
- 关联字段填充要批量查询，避免 N+1。

前端模板：

- 组件必须在 `parser/component_registry.go` 和 Vben adapter 中存在。
- 字典字段要处理缺少字典模块的配置分支。
- 搜索区字段过多时优先依赖 `keyword` 和优先级，不要把所有文本字段都生成独立控件。
- ID / bigint 在 TypeScript 侧按 string 处理，避免 JS 精度丢失。

## 推荐验证

不改行为，只改文档：

```bash
git diff -- admin-go/codegen/README.md admin-go/codegen/docs
```

改 parser 或模板：

```bash
cd admin-go/codegen
../../scripts/wait-for-cpu-idle.sh
../../scripts/run-go-task-with-limits.sh go run verify_codegen.go
../../scripts/run-go-task-with-limits.sh go run . --table demo_article --dry-run --manifest-out /tmp/baseadmin-codegen-manifest.json
```

有 MySQL 且需要端到端验证：

```bash
cd admin-go/codegen
../../scripts/run-go-task-with-limits.sh go run ./cmd/verifye2e --stage render --keep-temp
../../scripts/run-go-task-with-limits.sh go run ./cmd/verifye2e --stage dao --temp-root /tmp/baseadmin-codegen-e2e-xxxx
../../scripts/run-go-task-with-limits.sh go run ./cmd/verifye2e --stage build --temp-root /tmp/baseadmin-codegen-e2e-xxxx
```

如果项目策略允许单测，也必须走受限 Go 入口；当前协作要求下不要执行 `go test`。
