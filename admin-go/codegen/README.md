# GBaseAdmin 代码生成器

## 铁律

- 数据库结构和初始化真源在 `admin-go/database/migrations/`
- `codegen/sql/*.sql` 只用于模板验证、示例和本地辅助，不是部署真源
- 代码生成、DAO 生成、菜单写入、离线验证都基于 `go` / `gf` / 数据库结构完成
- 生成前端页面时，组件名必须受 `vue-vben-admin/apps/web-antd/src/adapter/component/index.ts` 约束
- 如果生成需求引入新组件，必须先在 adapter 完成适配，再修改模板和字段映射
- 生成后端逻辑时，必须遵守 GoFrame ORM 约定：写库优先使用 DO 对象，`created_at` / `updated_at` / `deleted_at` 依赖框架自动维护，软删除统一走 `Delete()`

## 表名规范

所有数据库表必须使用 `{应用名}_{模块名}` 格式命名：

| 表名 | 应用 | 模块 |
|------|------|------|
| `system_dept` | system | dept |
| `system_role` | system | role |
| `system_menu` | system | menu |
| `system_users` | system | users |
| `system_user_role` | system | user_role |
| `system_role_menu` | system | role_menu |
| `system_role_dept` | system | role_dept |
| `system_user_dept` | system | user_dept |
| `upload_dir` | upload | dir |
| `upload_file` | upload | file |
| `upload_config` | upload | config |
| `upload_dir_rule` | upload | dir_rule |
| `demo_demo` | demo | demo |

代码生成器根据第一个 `_` 拆分表名，前半部分为应用名，后半部分为模块名。

## 使用方法

```bash
cd admin-go/codegen

# 生成单表
go run . --table system_dept

# 生成单表后端/前端代码，不触发 gf gen dao
go run . --table system_dept --only backend

# 生成多表
go run . --table system_dept,system_role

# 只生成前端
go run . --table system_dept --only frontend

# 如需执行 gf gen dao，显式开启
go run . --table system_dept --only backend --with-dao

# 如需新建应用骨架，显式开启 gf init
go run . --table demo_demo --with-init

# 强制覆盖已有文件
go run . --table system_dept --force

# 预览（不写入文件）
go run . --table system_dept --dry-run

# 预览并输出结构化 manifest
go run . --table system_dept --dry-run --manifest-out ./tmp/codegen-manifest.json

# 只生成菜单数据（写入数据库）
go run . --table system_dept --only menu

# 生成代码同时写入菜单
go run . --table system_dept --menu
```

## 验证方式

推荐至少跑两层：

```bash
cd admin-go/codegen

# 1. 离线模板验证：不依赖真实生成目录
go run verify_codegen.go

# 2. 语法/单测验证
go test ./...

# 如模板有预期变更，更新 golden snapshot
go test ./... -run TestTemplateGoldenSnapshots -args -update-golden
```

如果本机有可用 MySQL 且 `admin-go/.env` 已配置完成，还可以继续跑生成级端到端验证：

```bash
cd admin-go/codegen

# 低资源服务器建议分阶段执行
go run ./cmd/verifye2e --stage render --keep-temp
go run ./cmd/verifye2e --stage dao --temp-root /tmp/baseadmin-codegen-e2e-xxxx
go run ./cmd/verifye2e --stage build --temp-root /tmp/baseadmin-codegen-e2e-xxxx

# 资源充足时也可以一次跑完
go run ./cmd/verifye2e --stage all
```

说明：

- `verify_codegen.go` 直接渲染模板并检查搜索、外键、字典、树形等关键片段是否生成
- `go test ./...` 现在包含 `testdata/golden/` 快照比对，模板输出漂移会直接在测试阶段暴露
- `cmd/verifye2e --stage render` 只做建表、模板渲染和工作区准备
- `cmd/verifye2e --stage dao` 单独执行 `gf gen dao`
- `cmd/verifye2e --stage build` 单独执行 `go build`
- `sql/e2e_verify.sql` 是 `verifye2e` 使用的专用验证表结构

## 失败退出规则

- `codegen` 主流程现在会收敛每张表、每个应用后置生成阶段的失败项
- 单个表的解析、前端生成、后端生成、菜单写入、应用目录创建、`internal/cmd` 目录创建、`hack/config.yaml`、`gf gen dao`、`main.go/cmd.go/middleware` 任一环节失败，最终都会非零退出
- 生成过程不会因为某一张表失败而立刻中断；会尽量继续跑完其余任务，并在最后输出失败摘要
- `dry-run` 同样适用这套规则；预览阶段只要有模板渲染失败，也会以非零退出结束

## 预览与落盘规则

- backend/frontend 生成现在先把所有目标文件渲染到内存，再统一写入磁盘
- 文件落盘阶段使用随机临时文件再原子替换目标文件，避免固定 `.codegen-tmp` 名称在并发或中断后互相覆盖
- `dry-run` 会基于真实 plan 输出更准确的状态：`新文件`、`有变化`、`无变化`、`跳过（已存在）`、`保护（enhance 文件不覆盖）`
- 传 `--manifest-out` 时，CLI 会把文件计划和菜单计划写入 JSON，便于 CI、审查工具或人工核对
- 当前 manifest 已覆盖模板来源、输出路径、动作类型、文件大小，以及菜单目录/页面/按钮计划

## 菜单批量写入规则

- `--menu` 和 `--only menu` 现在会先收集本次成功解析的表，再统一批量生成菜单
- 非 `dry-run` 模式下，整批菜单会复用同一个数据库连接，并在同一个事务里写入目录、页面和按钮
- 批量写库前会先做本地预校验，提前拦截空元数据、缺应用名/模块名、同批次菜单 path 冲突和权限冲突
- 只要任意一个模块的目录、菜单页或按钮写入失败，整批菜单事务会回滚，CLI 最终以非零退出
- `dry-run` 仍然只打印预览，不会写库

## Parser 分层约定

`ParseTable()` 现在按固定阶段执行：

1. 表身份解析：从 `{app}_{module}` 拆出应用和模块
2. 字段元数据构建：列信息转 `FieldMeta`
3. 关联字段解析：统一补齐外键和 `parent_id` 的显示字段、关联表信息
4. 模板派生元数据收口：由 `FinalizeTemplateMeta()` 统一计算 `Has*`、`SearchFields`、`KeywordSearchFields`

后续如果继续扩展规则，优先加到对应阶段函数里，不要把新规则继续堆回 `ParseTable()` 主流程。

### 规则注册表

- 字段隐藏、图片字段、可搜索文本、精确搜索、金额字段、关联显示字段优先级现在统一收敛在 `parser/rule_registry.go`
- `field_mapper.go`、`parser.go` 和关联字段推断共享同一套规则，不再各自维护近似但不完全一致的硬编码列表
- 后续若要扩展命名启发式，优先改注册表和对应测试，而不是在多个函数里重复加 `if/else`

### CLI 参数一览

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `--table` | string | （必填） | 表名，多表用逗号分隔 |
| `--only` | string | （空） | 只生成指定部分：`backend`、`frontend`、`menu` |
| `--force` | bool | false | 强制覆盖已有文件 |
| `--config` | string | `./codegen.yaml` | 配置文件路径 |
| `--dry-run` | bool | false | 预览模式，不写入文件 |
| `--menu` | bool | false | 生成代码同时写入菜单数据到数据库 |
| `--with-dao` | bool | false | 生成后显式执行 `gf gen dao` |
| `--with-init` | bool | false | 应用目录不存在时显式执行 `gf init` |
| `--manifest-out` | string | （空） | 将本次文件与菜单计划写入 JSON manifest |

## 自动创建应用

当表名前缀对应的应用目录不存在时，代码生成器默认只创建目录并写入生成文件；如需额外执行 `gf init app/{appName} -a` 创建应用骨架，请显式传 `--with-init`。

例如：`--table demo_demo --with-init` 会创建 `app/demo/` 应用骨架。

## 生成文件列表

### 后端（输出到 `app/{appName}/`）

| 模板 | 输出路径 | 说明 |
|------|---------|------|
| `api.tpl` | `api/{app}/v1/{module}.go` | API 请求/响应结构体 |
| `controller.tpl` | `internal/controller/{module}/{module}.go` | 控制器 |
| `logic.tpl` | `internal/logic/{module}/{module}.go` | 业务逻辑 |
| `service.tpl` | `internal/service/{module}.go` | 服务接口 |
| `model.tpl` | `internal/model/{module}.go` | DTO 模型 |
| `consts.tpl` | `internal/consts/{module}.go` | 枚举常量 |

### 前端（输出到 `vue-vben-admin/apps/web-antd/src/`）

| 模板 | 输出路径 | 说明 |
|------|---------|------|
| `types.tpl` | `api/{app}/{module}/types.ts` | TypeScript 类型定义 |
| `api.tpl` | `api/{app}/{module}/index.ts` | API 请求函数 |
| `list.tpl` | `views/{app}/{module}/index.vue` | 列表页面 |
| `form.tpl` | `views/{app}/{module}/modules/form.vue` | 表单弹窗 |
| `detail-drawer.tpl` | `views/{app}/{module}/modules/detail-drawer.vue` | 详情抽屉 |

## 配置文件

`codegen.yaml`：

```yaml
database:
  host: ${MYSQL_HOST}
  port: ${MYSQL_PORT}
  user: ${MYSQL_USER}
  password: ${MYSQL_PASSWORD}
  dbname: ${MYSQL_DATABASE}

backend:
  output: ../app/             # 后端输出根目录

frontend:
  output: ../../vue-vben-admin/apps/web-antd/src/

skip_fields:                  # 这些字段在表单中隐藏，不生成前端组件
  - created_at
  - updated_at
  - deleted_at
  - created_by
  - dept_id

menu_apps:                    # 菜单应用目录配置（新增应用在此添加即可）
  system:
    title: 系统管理
    icon: SettingOutlined
  upload:
    title: 上传管理
    icon: CloudUploadOutlined

menu_modules:                 # 模块级菜单配置（可选）
  upload/file:
    icon: FileTextOutlined
  upload/dir_rule:
    icon: PartitionOutlined
```

### 环境变量支持

推荐统一使用 `admin-go/.env`。`codegen` 会自动读取该文件，再展开 `codegen.yaml` 里的 `${...}` 变量：

```yaml
database:
  host: ${MYSQL_HOST}
  port: ${MYSQL_PORT}
  user: ${MYSQL_USER}
  password: ${MYSQL_PASSWORD}
  dbname: ${MYSQL_DATABASE}
```

### 菜单应用配置

`menu_apps` 定义了菜单生成器为每个应用创建目录时使用的标题和图标。新增应用只需在此添加一条配置，无需修改生成器源码。未配置的应用默认使用 `{应用名}管理` 和 `AppstoreOutlined` 图标。

### 模块菜单图标

页面菜单默认会根据模块名和表注释自动选择一个较合适的 Ant Design 图标，例如 `users -> UserOutlined`、`dept -> ApartmentOutlined`、`menu -> MenuOutlined`、`config -> ControlOutlined`、`dir -> FolderOpenOutlined`、`dir_rule -> PartitionOutlined`、`file -> FileTextOutlined`。

如果自动推断不符合业务语义，可以在 `menu_modules` 里按 `应用名/模块名` 单独覆盖：

```yaml
menu_modules:
  upload/file:
    sort: 40
    icon: FileTextOutlined
    is_show: 1
```

## 数据库表设计规范

- 主键统一使用 `id BIGINT UNSIGNED`（Snowflake ID）
- 软删除使用 `deleted_at DATETIME`
- 公共字段：`created_at`、`updated_at`、`deleted_at`、`created_by`、`dept_id`
- 树形结构使用 `parent_id BIGINT UNSIGNED`
- 状态字段使用 `status TINYINT(1)`，注释格式：`状态:0=关闭,1=开启`
- 枚举字段注释格式：`字段说明:值1=标签1,值2=标签2`
- 外键字段命名：`{关联模块}_id`（如 `dept_id`、`role_id`）
- 多选外键字段命名：`{关联模块}_ids`（如 `role_ids`）

## 字段注释与枚举格式

数据库字段的 `COMMENT` 决定了前端表单标签和组件类型。格式：

```
{标签}:{值1}={显示名1},{值2}={显示名2},...
```

### Tooltip 提示（括号语法）

当字段注释的标签部分包含中文括号 `（）` 或英文括号 `()` 时，括号内的内容会自动提取为 Tooltip 提示文字，括号前的文字作为精简标签。

```
{精简标签}（{提示文字}）
```

效果：表单 label 和表格列头显示精简标签 + 问号图标，鼠标悬停显示提示文字。

示例：

| 注释 | 精简标签 | Tooltip 提示 | 前端效果 |
|------|---------|-------------|---------|
| `部门名称` | 部门名称 | （无） | 普通文字标签 |
| `排序（升序）` | 排序 | 升序 | 排序 ❓ |
| `支付金额（分）` | 支付金额 | 分 | 支付金额 ❓ |
| `面值（分，满减时为抵扣额）` | 面值 | 分，满减时为抵扣额 | 面值 ❓ |
| `状态:0=关闭,1=开启` | 状态 | （无） | 普通文字标签 + 枚举 |

> 括号语法和枚举语法可以组合使用：`排序（升序）:0=默认,1=热门` → 精简标签="排序"，Tooltip="升序"，枚举=[{0,"默认"},{1,"热门"}]

### 搜索指令

字段注释还支持通过 `|` 追加搜索指令，覆盖默认的智能识别规则：

```sql
`summary` varchar(255) COMMENT '摘要|search:off|keyword:only|priority:95'
`order_no` varchar(64) COMMENT '订单号|search:eq'
`category_id` bigint COMMENT '分类|search:tree'
```

可用指令：

- `search:off`：不生成独立搜索控件
- `search:on`：启用智能识别
- `search:eq`：强制精确搜索
- `search:like`：强制模糊搜索
- `search:range`：强制区间搜索
- `search:select`：强制下拉搜索
- `search:tree`：强制树形下拉搜索
- `keyword:on`：加入全局关键词搜索
- `keyword:off`：不加入全局关键词搜索
- `keyword:only`：只加入全局关键词搜索，不单独生成控件
- `priority:90`：指定搜索优先级，值越大越靠前

### 枚举格式

示例：

| 注释 | 解析结果 |
|------|---------|
| `部门名称` | 标签="部门名称"，无枚举 |
| `状态:0=关闭,1=开启` | 标签="状态"，枚举=[{0,"关闭"},{1,"开启"}] |
| `类型:1=普通,2=VIP,3=管理员` | 标签="类型"，枚举=[{1,"普通"},{2,"VIP"},{3,"管理员"}] |

枚举字段会自动在后端生成 Go 常量（`internal/consts/{module}.go`），前端生成对应的 options 数组。

### 外键显式声明

默认情况下，`*_id` 会优先按“当前应用同名前缀表”推断关联表，例如 `demo_article.category_id -> demo_category`。

如果字段需要指向跨应用表，或字段名不足以唯一推断关联表，可以在字段注释后追加 `|ref:` 指令：

```sql
`user_id` BIGINT UNSIGNED NOT NULL COMMENT '作者|ref:system_users.username'
```

- `system_users` 表示关联表名
- `username` 表示显示字段；省略时会按默认优先级自动查找显示字段

生成器会据此：
- 后端按 `system_users` 填充关联展示字段
- 前端自动导入对应应用的列表 API
- 搜索表单自动生成外键下拉

## 前端组件自动映射

代码生成器根据字段名、数据库类型、枚举数量自动选择前端组件，无需手动配置。

### 按字段名匹配

| 字段名模式 | 映射组件 | 说明 |
|-----------|---------|------|
| `parent_id` | `TreeSelectSingle` | 树形单选 |
| `parent_ids` | `TreeSelectMulti` | 树形多选 |
| `*_ids` | `SelectMulti` | 多选下拉框 |
| `*_id`（排除 `id`、`dept_id`） | `Select` | 单选下拉框（外键关联） |
| `password`、`*_password`、`*_pwd` | `Password` | 密码输入框 |
| `*_url`、`*_link` | `InputUrl` | URL 输入框 |
| `*_at` | `DateTimePicker` | 日期时间选择器 |
| `sort`、`order`、`*_num`、`*_price`、`*_amount`、`*_income`、`*_balance` | `InputNumber` | 数字输入框 |
| `icon` | `IconPicker` | 图标选择器 |
| `avatar`、`cover`、`logo`、`banner`、`thumbnail`、`poster` | `ImageUpload` | 图片上传（精确名匹配） |
| `*_image`、`*_img`、`*_photo`、`*_pic`、`*_cover`、`*_banner`、`*_logo`、`*_thumbnail`、`*_poster` | `ImageUpload` | 图片上传（后缀匹配） |
| `*_file`、`*_attachment` | `FileUpload` | 文件上传 |
| `*_content`、`*_body`、`*_html` | `RichText` | 富文本编辑器 |
| `*_json`、`*_config`、`*_settings` | `JsonEditor` | JSON 编辑器 |

### 富组件自动映射

`parser/field_mapper.go` 中的 `MapComponent` 函数根据字段命名规则自动映射富组件：

| 字段名规则 | 组件 | 说明 |
|-----------|------|------|
| `avatar`、`cover`、`logo`、`banner`、`thumbnail`、`poster`、`*_image`、`*_img`、`*_photo`、`*_pic`、`*_cover` 等 | `ImageUpload` | 图片上传（文件管理器 `mode=image`） |
| `*_file`、`*_attachment` | `FileUpload` | 文件上传（文件管理器 `mode=all`） |
| `*_content`、`*_body`、`*_html` | `RichText` | TinyMCE 富文本编辑器 |
| `*_json`、`*_config`、`*_settings` | `JsonEditor` | JSON 编辑器（tree + code 双模式） |

这些组件在 `adapter/component/index.ts` 中注册，上传类组件调用 `/upload/uploader/upload` 接口。

> **注意**：使用 `ImageUpload` 或 `FileUpload` 组件的表单页面需要 `upload` 应用处于运行状态，否则文件管理 API 不可用。

### 按枚举数量匹配

| 条件 | 映射组件 |
|------|---------|
| `status`/`is_*` + 2个枚举值 | `Switch` | 开关切换 |
| `status`/`is_*` + 3个以上枚举值 | `Radio` | 单选按钮组 |
| `type`/`level`/`grade` | `Select` | 下拉选择 |

### 按数据库类型匹配

| 数据库类型 | 映射组件 |
|-----------|---------|
| `TEXT`、`LONGTEXT`、`MEDIUMTEXT`、`TINYTEXT` | `Textarea` |
| 其他 | `Input` |

### 全部可用组件

下列组件名是 `codegen` 的单一契约真源，定义在 `parser/component_registry.go`，`field_mapper`、`verify_codegen.go`、模板契约测试共享同一份列表。

| 组件名 | 说明 |
|--------|------|
| `Input` | 文本输入框 |
| `InputNumber` | 数字输入框 |
| `Textarea` | 多行文本框 |
| `Switch` | 开关（两态切换） |
| `Radio` | 单选按钮组 |
| `Select` | 下拉选择（单选） |
| `SelectMulti` | 下拉选择（多选） |
| `TreeSelectSingle` | 树形下拉（单选） |
| `TreeSelectMulti` | 树形下拉（多选） |
| `ImageUpload` | 图片上传 |
| `FileUpload` | 文件上传 |
| `RichText` | 富文本编辑器 |
| `JsonEditor` | JSON 编辑器 |
| `Password` | 密码输入框 |
| `InputUrl` | URL 输入框 |
| `DateTimePicker` | 日期时间选择器 |
| `IconPicker` | 图标选择器 |

另外，`go test ./...` 会校验 `templates/frontend/form.tpl` 对这份组件清单的分支覆盖，避免“新增组件常量但模板漏接”的漂移。

## 类型映射

### Go 类型映射

| 数据库类型 | Go 类型 |
|-----------|---------|
| `BIGINT UNSIGNED` | `JsonInt64`（防止前端精度丢失） |
| `BIGINT` | `JsonInt64` |
| `INT`、`MEDIUMINT`、`SMALLINT`、`TINYINT` | `int` |
| `FLOAT`、`DOUBLE`、`DECIMAL` | `float64` |
| `DATETIME`、`TIMESTAMP`、`DATE`、`TIME` | `*gtime.Time` |
| 其他 | `string` |

### TypeScript 类型映射

| 数据库类型 | TS 类型 |
|-----------|---------|
| `BIGINT` | `string`（Snowflake ID 防精度丢失） |
| `INT`、`TINYINT`、`FLOAT`、`DOUBLE`、`DECIMAL` | `number` |
| 其他 | `string` |

## 隐藏字段

以下字段自动从表单中排除（不出现在新增/编辑表单中）：

- `id` — 主键，自动生成
- `created_at`、`updated_at`、`deleted_at` — 时间戳，自动维护
- `created_by` — 创建人，自动填充
- `dept_id` — 部门 ID，自动填充

此外，`codegen.yaml` 中的 `skip_fields` 配置项可自定义额外需要隐藏的字段。

## 智能特性检测

| 特性 | 触发条件 | 生成效果 |
|------|---------|---------|
| 树形结构 | 表中存在 `parent_id` 字段 | 后端生成树形查询接口，前端生成树形表格 |
| 密码加密 | 字段名为 `password`/`*_password`/`*_pwd` | 后端自动 bcrypt 加密 |
| 外键关联 | 字段名为 `*_id`（排除 `id`、`dept_id`） | 自动批量查询关联表（`WHERE id IN (...)`），填充显示字段（title/name/username/nickname/real_name/label/phone/mobile） |
| 多选外键 | 字段名为 `*_ids` | 前端多选组件，后端数组处理 |
| Snowflake ID | 所有 `BIGINT` 主键/外键 | 使用 `JsonInt64` 防止 JS 精度丢失 |
| 软删除 | 存在 `deleted_at` 字段 | 查询自动过滤已删除记录 |
| 枚举常量 | 字段注释包含枚举定义 | 后端生成 Go 常量，前端生成 options |
| Tooltip 提示 | 字段注释标签含 `（）` 或 `()` | 前端表单 label 和列头自动渲染 Tooltip 问号图标 |
| 模糊搜索 | 字段名为 `title`/`name`/`phone`/`email` 等 | 后端 `WhereLike` 模糊查询，前端搜索栏自动添加 Input |
| 精确搜索 | 字段名后缀 `_no`/`_code`/`_sn` | 编号类字段用精确匹配 `Where` 而非 `WhereLike` |
| 金额格式化 | 字段名含 `price`/`amount`/`balance`/`fee`/`cost` | 列表自动"分→元"格式化显示（`/ 100`） |
| 批量删除 | 所有表 | 前端勾选框 + 批量删除按钮，后端 `WhereIn` 批量软删除 |
| CSV 导出 | 所有表 | 后端 CSV 流式输出，前端导出按钮（Blob 下载） |
| 详情抽屉 | 所有表 | 只读详情展示，枚举 Tag、图片预览、富文本渲染 |
| 时间范围筛选 | 所有表 | 前端 RangePicker + 后端 `created_at` 区间查询 |
| 列表排序 | 所有表 | 前端列头排序 + 后端动态 `OrderBy`/`OrderDir` |

## 生成的 CRUD 功能清单

每张表自动生成以下完整功能：

### 后端接口

| 接口 | 方法 | 路径 | 说明 |
|------|------|------|------|
| 创建 | POST | `/{module}/create` | 新增记录 |
| 更新 | PUT | `/{module}/update` | 修改记录 |
| 删除 | DELETE | `/{module}/delete` | 软删除单条 |
| 批量删除 | DELETE | `/{module}/batch-delete` | 软删除多条 |
| 详情 | GET | `/{module}/detail` | 获取单条详情（含关联字段） |
| 列表 | GET | `/{module}/list` | 分页列表（支持搜索、筛选、排序、时间范围） |
| 导出 | GET | `/{module}/export` | CSV 导出（最多 10000 条，支持筛选条件） |
| 树形 | GET | `/{module}/tree` | 仅 `parent_id` 表，返回树形结构 |

### 前端页面

| 组件 | 功能 |
|------|------|
| 列表页 | VxeTable 表格 + 搜索表单 + 时间范围筛选 + 列头排序 + 枚举 Tag + 金额格式化 |
| 表单弹窗 | 新建/编辑表单，自动组件映射，密码字段编辑时可选填 |
| 详情抽屉 | 只读详情展示（Descriptions），枚举 Tag、图片预览、富文本/JSON 渲染 |
| 工具栏 | 新建按钮 + 批量删除按钮（非树形）+ 导出按钮 |

### 枚举常量语义化

枚举字段自动生成语义化 Go 常量，覆盖 30+ 常见中文标签映射：

```go
// 改前（数字兜底）
const UserStatusV0 = 0 // 禁用
const UserStatusV1 = 1 // 启用

// 改后（语义化）
const UserStatusDisabled = 0 // 禁用
const UserStatusEnabled = 1  // 启用
```

支持的映射包括：启用/禁用、开启/关闭、是/否、男/女、显示/隐藏、待处理/进行中/已完成、待审核/已通过/已拒绝、待支付/已支付/已退款、草稿/已发布/已下架、成功/失败、充值/消费/提现等。未匹配的标签使用 `V{数值}` 兜底。

## 菜单生成

使用 `--menu` 或 `--only menu` 可将菜单数据写入 `system_menu` 表。

每个模块生成以下菜单记录：

| 菜单 | 类型 | 说明 |
|------|------|------|
| `{模块名}管理` | 目录 | 一级菜单，挂载到对应应用目录下 |
| `{模块名}列表` | 页面 | 列表页面，路由指向生成的 `index.vue` |
| 按钮权限 | 按钮 | 包含 新增、修改、删除、查看、导出 五个操作按钮 |

菜单写入前会检查是否已存在，避免重复插入。

## 命名转换规则

| snake_case | CamelCase（Go 导出） | DAO 风格 | camelCase（JSON/TS） |
|-----------|---------------------|----------|---------------------|
| `dept_name` | `DeptName` | `DeptName` | `deptName` |
| `parent_id` | `ParentID` | `ParentId` | `parentID` |
| `link_url` | `LinkURL` | `LinkUrl` | `linkURL` |
| `id` | `ID` | `Id` | `id` |

常见缩写（`ID`、`URL`、`IP`、`API`、`HTTP` 等）在 CamelCase 中保持全大写。
