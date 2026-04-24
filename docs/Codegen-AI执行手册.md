# Codegen AI 执行手册

本文给 AI/代理执行 `admin-go/codegen/` 任务使用。目标是把 CRUD 生成说清楚到可以直接执行，不依赖口头经验。

## 1. 适用范围

当前仓库是 GBaseAdmin 基础框架，只允许围绕这些范围做 codegen：

```text
admin-go/app/system
admin-go/app/upload
vue-vben-admin/apps/web-antd/src/api/system
vue-vben-admin/apps/web-antd/src/api/upload
vue-vben-admin/apps/web-antd/src/views/system
vue-vben-admin/apps/web-antd/src/views/upload
admin-go/codegen
admin-go/database/migrations
```

`admin-go/codegen/codegen.yaml` 的 `allowed_apps` 默认只有：

```yaml
allowed_apps:
  - system
  - upload
```

如果需求属于其他业务项目，不要在本仓库新增业务应用。应在对应业务仓库处理。

## 2. AI 必读顺序

处理 codegen 前按顺序读取：

1. `CLAUDE.md`
2. `docs/代码库导读.md`
3. 本文
4. `admin-go/codegen/README.md`
5. 相关现有页面，例如 `vue-vben-admin/apps/web-antd/src/views/system/*` 或 `vue-vben-admin/apps/web-antd/src/views/upload/*`

如果这些文档冲突，以 `CLAUDE.md` 的铁律优先，再同步更新文档。

## 3. 资源保护命令

在服务器上执行任何命令前先检查负载：

```bash
./scripts/wait-for-cpu-idle.sh
```

在 `admin-go/codegen` 目录执行 Go 任务时必须使用受限入口：

```bash
../../scripts/run-go-task-with-limits.sh go run . --table upload_dir_rule --dry-run
../../scripts/run-go-task-with-limits.sh go test ./...
```

前端排查必须使用 Node 限流入口，禁止裸 `pnpm`、`npm`、`npx`：

```bash
./scripts/run-node-task-with-limits.sh pnpm -C vue-vben-admin -F @vben/web-antd typecheck
```

禁止在本机执行 Docker。需要 Docker 行为时只更新文档或配置，由 CI/部署环境执行。

## 4. Codegen 任务卡

AI 开始生成前，必须先把需求整理成任务卡。信息不完整时先从现有表、迁移和页面推断；仍无法确认再问用户。

```yaml
codegen_task:
  app: upload
  module: dir_rule
  table: upload_dir_rule
  purpose: 管理上传目录规则
  migration_required: true
  migration_file: admin-go/database/migrations/000009_xxx.up.sql
  generate:
    backend: true
    frontend: true
    menu: false
    with_dao: false
    with_init: false
  overwrite:
    force: false
  fields:
    - name: id
      type: BIGINT UNSIGNED
      comment: 主键
    - name: title
      type: VARCHAR(64)
      comment: 规则名称|search:like|keyword:on|priority:90
    - name: status
      type: TINYINT(1)
      comment: 状态:0=关闭,1=开启
  expected_outputs:
    backend:
      - admin-go/app/upload/api/upload/v1/dir_rule.go
      - admin-go/app/upload/internal/controller/dir_rule/dir_rule.go
      - admin-go/app/upload/internal/logic/dir_rule/dir_rule.go
      - admin-go/app/upload/internal/service/dir_rule.go
      - admin-go/app/upload/internal/model/dir_rule.go
      - admin-go/app/upload/internal/consts/dir_rule.go
    frontend:
      - vue-vben-admin/apps/web-antd/src/api/upload/dir_rule/index.ts
      - vue-vben-admin/apps/web-antd/src/api/upload/dir_rule/types.ts
      - vue-vben-admin/apps/web-antd/src/views/upload/dir_rule/index.vue
      - vue-vben-admin/apps/web-antd/src/views/upload/dir_rule/modules/form.vue
      - vue-vben-admin/apps/web-antd/src/views/upload/dir_rule/modules/detail-drawer.vue
  validation:
    - bash ./scripts/verify-vben-pages.sh
    - cd admin-go/codegen && ../../scripts/run-go-task-with-limits.sh go test ./...
```

## 5. 表结构契约

### 表名

表名必须是：

```text
{app}_{module}
```

示例：

```text
system_dept
system_role
system_users
upload_file
upload_dir_rule
```

生成器按第一个 `_` 拆分：

```text
upload_dir_rule -> app=upload, module=dir_rule
```

多段模块名保留为包名：

```text
user_review -> package user_review
```

### 常规字段

推荐基础结构：

```sql
`id` BIGINT UNSIGNED NOT NULL COMMENT '主键',
`created_at` DATETIME NULL COMMENT '创建时间',
`updated_at` DATETIME NULL COMMENT '更新时间',
`deleted_at` DATETIME NULL COMMENT '删除时间',
`created_by` BIGINT UNSIGNED NULL COMMENT '创建人',
`dept_id` BIGINT UNSIGNED NULL COMMENT '部门'
```

约定：

- `id`、`created_at`、`updated_at`、`deleted_at`、`created_by`、`dept_id` 默认不进入表单。
- 有 `deleted_at` 时，生成逻辑按软删除处理。
- 有 `created_by` 和 `dept_id` 时，生成逻辑要考虑数据权限。
- 有 `parent_id` 时，生成树形列表和树形查询。
- 有 `sort` 时，列表默认按排序字段处理。
- `BIGINT` 主键和外键在前端按 `string` 处理，避免 JS 精度丢失。

### 字段注释语法

字段 `COMMENT` 是 codegen 的重要输入。AI 写迁移时必须认真设计注释。

基本格式：

```text
标签
标签:值1=显示名1,值2=显示名2
标签（提示文字）
标签:dict:字典类型
标签|指令1|指令2
```

示例：

```sql
`title` varchar(64) COMMENT '标题|search:like|keyword:on|priority:90'
`status` tinyint(1) COMMENT '状态:0=关闭,1=开启'
`type` tinyint(1) COMMENT '类型:1=普通,2=高级,3=系统'
`amount` int COMMENT '金额（分）|search:range'
`user_id` bigint unsigned COMMENT '用户|ref:system_users.username|search:select'
`gender` tinyint(1) COMMENT '性别:dict:gender'
```

可用指令：

| 指令 | 作用 |
| --- | --- |
| `ref:system_users.username` | 指定外键表和显示字段 |
| `search:off` | 不生成搜索控件 |
| `search:on` | 启用智能搜索 |
| `search:eq` | 精确搜索 |
| `search:like` | 模糊搜索 |
| `search:range` | 区间搜索 |
| `search:select` | 下拉搜索 |
| `search:tree` | 树形搜索 |
| `keyword:on` | 加入全局关键词搜索 |
| `keyword:off` | 不加入全局关键词搜索 |
| `keyword:only` | 只加入全局关键词搜索，不单独生成搜索控件 |
| `priority:90` | 搜索控件排序优先级 |

## 6. 组件映射契约

AI 不要在模板里发明组件名。可用组件只来自：

```text
admin-go/codegen/parser/component_registry.go
vue-vben-admin/apps/web-antd/src/adapter/component/index.ts
```

字段命名到组件的主要映射：

| 字段规则 | 组件 |
| --- | --- |
| `parent_id` | `TreeSelectSingle` |
| `parent_ids` | `TreeSelectMulti` |
| `*_ids` | `SelectMulti` |
| `*_id`，排除 `id`、`dept_id` | `Select` |
| `status` 或 `is_*`，两值枚举 | `Switch` |
| `status` 或 `is_*`，多值枚举 | `Radio` |
| `type`、`level`、`grade` | `Select` |
| `password`、`*_password`、`*_pwd`、`*_secret`、`*_token` | `Password` |
| `*_url`、`*_link` | `InputUrl` |
| `*_at` | `DateTimePicker` |
| `sort`、`order`、`*_num`、`*_price`、`*_amount`、`*_income`、`*_balance` | `InputNumber` |
| `icon` | `IconPicker` |
| `avatar`、`cover`、`logo`、`banner`、`thumbnail`、`poster`、`*_image`、`*_img`、`*_photo`、`*_pic` | `ImageUpload` |
| `*_file`、`*_attachment` | `FileUpload` |
| `*_content`、`*_body`、`*_html` | `RichText` |
| `*_json`、`*_config`、`*_settings` | `JsonEditor` |
| `TEXT`、`LONGTEXT`、`MEDIUMTEXT`、`TINYTEXT` | `Textarea` |
| 其他 | `Input` |

如果业务需要新组件，顺序必须是：

```text
adapter/component/index.ts
-> admin-go/codegen/parser/component_registry.go
-> admin-go/codegen/parser/field_mapper.go
-> admin-go/codegen/templates/frontend/form.tpl
-> admin-go/codegen/verify_codegen.go
-> admin-go/codegen/testdata/golden
-> scripts/verify-vben-pages.sh（如页面规则需要扩展）
```

## 7. 标准执行流程

### 7.1 生成前检查

从仓库根目录执行：

```bash
./scripts/wait-for-cpu-idle.sh
git status --short
```

如果工作区有无关改动，不要覆盖。只修改本任务需要的文件。

### 7.2 数据库先行

如果任务涉及新增表或改字段，先写迁移：

```text
admin-go/database/migrations/0000xx_xxx.up.sql
admin-go/database/migrations/0000xx_xxx.down.sql
```

不要把 `admin-go/codegen/sql/*.sql` 当部署真源。`codegen/sql` 只用于验证样例。

### 7.3 预览

在 `admin-go/codegen` 目录执行：

```bash
../../scripts/wait-for-cpu-idle.sh
../../scripts/run-go-task-with-limits.sh go run . \
  --table upload_dir_rule \
  --dry-run \
  --manifest-out /tmp/baseadmin-codegen-manifest.json
```

AI 必须检查 manifest 中的输出路径，确认没有写到 `system/upload` 范围之外。

### 7.4 落盘生成

只生成后端：

```bash
../../scripts/wait-for-cpu-idle.sh
../../scripts/run-go-task-with-limits.sh go run . --table upload_dir_rule --only backend
```

只生成前端：

```bash
../../scripts/wait-for-cpu-idle.sh
../../scripts/run-go-task-with-limits.sh go run . --table upload_dir_rule --only frontend
```

后端、前端一起生成：

```bash
../../scripts/wait-for-cpu-idle.sh
../../scripts/run-go-task-with-limits.sh go run . --table upload_dir_rule
```

菜单写库只有明确需要时才执行：

```bash
../../scripts/wait-for-cpu-idle.sh
../../scripts/run-go-task-with-limits.sh go run . --table upload_dir_rule --only menu
```

`--force` 只能在确认覆盖的是生成文件时使用。不要用它覆盖手写增强文件。

### 7.5 DAO 处理

`--with-dao` 默认关闭。需要 DAO 时才显式执行：

```bash
../../scripts/wait-for-cpu-idle.sh
../../scripts/run-go-task-with-limits.sh go run . --table upload_dir_rule --only backend --with-dao
```

禁止手改：

```text
admin-go/app/*/dao/
admin-go/app/*/internal/model/do/
admin-go/app/*/internal/model/entity/
```

这些目录只能由 `gf gen dao` 生成。

## 8. 验证要求

### 文档或小改动

至少执行：

```bash
./scripts/wait-for-cpu-idle.sh
git diff --check
```

### 修改 codegen

至少执行：

```bash
./scripts/wait-for-cpu-idle.sh
cd admin-go/codegen
../../scripts/run-go-task-with-limits.sh go run verify_codegen.go
../../scripts/run-go-task-with-limits.sh go test ./...
```

如果模板输出有预期变化：

```bash
../../scripts/run-go-task-with-limits.sh go test ./... -run TestTemplateGoldenSnapshots -args -update-golden
../../scripts/run-go-task-with-limits.sh go test ./...
```

如果涉及生成页面：

```bash
./scripts/wait-for-cpu-idle.sh
bash ./scripts/verify-vben-pages.sh
```

如果涉及前端模板、adapter 或类型：

```bash
./scripts/wait-for-cpu-idle.sh
./scripts/run-node-task-with-limits.sh pnpm -C vue-vben-admin -F @vben/web-antd typecheck
```

### codegen 覆盖矩阵

修改 `admin-go/codegen/` 后，离线验证和样例必须覆盖：

- 树形表和非树形表。
- 单选外键和跨应用外键。
- 树形外键和普通外键。
- 所有注册组件：`Input`、`InputNumber`、`Textarea`、`Switch`、`Radio`、`Select`、`TreeSelectSingle`、`TreeSelectMulti`、`SelectMulti`、`ImageUpload`、`FileUpload`、`RichText`、`JsonEditor`、`Password`、`InputUrl`、`DateTimePicker`、`IconPicker`。
- 枚举字段：两值 Switch、三值 Radio、多值 Select。
- 金额字段，单位为分，列表按元展示。
- 搜索字段：模糊、精确、区间、下拉、树形。
- Tooltip 字段。
- 密码字段。
- 多段模块名。
- 字段 comment 为空的回退。
- `dict:xxx` 字典字段。
- 数据权限字段：`created_by` 和 `dept_id`。
- `sort` 排序字段。
- 自定义时间字段：`*_at`。
- 验证规则：email、phone、url、max-length。

## 9. 常见决策

### 什么时候必须用 codegen

- 新增标准 CRUD。
- 扩展标准列表、详情、表单、导出、批量删除。
- 需要后端 API、前端页面和菜单保持一致。
- 需要从数据库注释统一生成 label、枚举、搜索、组件。

### 什么时候不要用 codegen

- 登录、授权、上传、文件选择器、复杂工作流。
- 需要强业务编排的页面。
- 已有手写页面只做小范围交互修复。
- 不符合 `{app}_{module}` 表名规则的临时脚本。

### 外键如何设计

优先让字段名和关联表同名：

```sql
`dir_id` BIGINT UNSIGNED COMMENT '目录'
```

跨应用或显示字段不明确时必须显式写 `ref`：

```sql
`user_id` BIGINT UNSIGNED COMMENT '用户|ref:system_users.username'
```

如果关联表有 `parent_id`，前端会倾向树形选择。

### 搜索如何设计

常用选择：

```sql
`title` varchar(64) COMMENT '标题|search:like|keyword:on|priority:90'
`code` varchar(64) COMMENT '编码|search:eq|priority:80'
`status` tinyint(1) COMMENT '状态:0=关闭,1=开启|search:select'
`created_at` datetime COMMENT '创建时间|search:range'
```

不要让所有字段都进搜索栏。搜索控件过多会让生成页面难用。

## 10. 失败处理

- codegen 非零退出后，不要反复重试。
- 先看失败阶段：解析、作用域、模板渲染、落盘、菜单写库、`gf gen dao`。
- 数据库连不上时，不要改模板绕过；先报告环境问题。
- 模板缺组件时，不要在模板里硬写裸组件；先补 adapter 和注册表。
- 生成路径越界时，先检查表名和 `allowed_apps`。
- 如果已经生成了部分文件，先用 `git status --short` 和 `git diff` 明确改动，再继续。

## 11. AI 交付说明模板

完成后回复用户时只说高信号内容：

```text
已更新 codegen 文档，新增了 AI 执行手册，并把 README/CLAUDE/文档索引接上。
验证：git diff --check 通过。
提交：docs(codegen): add AI execution guide
```

如果未执行某项验证，要明确说明原因。
