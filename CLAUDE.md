# GBaseAdmin AI 协作说明

## 必须先阅读

1. [仓库总说明](README.MD)
2. [代码库导读](docs/代码库导读.md)（首选入口）
3. [Codegen AI 执行手册](docs/Codegen-AI执行手册.md)（涉及 `admin-go/codegen/`、CRUD 生成、模板或 Vben 生成页时必读）
4. [代码生成器说明](admin-go/codegen/README.md)
5. [Docker 开发说明](docs/Docker开发说明.md)
6. [system 服务路由入口](admin-go/app/system/internal/cmd/cmd.go)
7. [upload 服务路由入口](admin-go/app/upload/internal/cmd/cmd.go)

当前仓库已经收缩为后台基础框架，AI/代理默认只围绕以下范围工作：

- 后端：`admin-go/app/system`
- 后端：`admin-go/app/upload`
- 管理端：`vue-vben-admin/apps/web-antd/src/` 下的 `system` / `upload` 以及后台公共壳
- 生成器：`admin-go/codegen/`
- Docker 入口：`docker/dev/`、`docker/prod/`

## 目录边界

```text
docker/
docs/
admin-go/
vue-vben-admin/apps/web-antd/src/
```

历史 `.claude/agents/*.md` 分工文档已经移除；协作约束只保留当前这份 `CLAUDE.md`，避免再引用旧路径和旧架构说明。

谨慎处理区：

```text
vue-vben-admin/packages/
vue-vben-admin/internal/
vue-vben-admin/playground/
admin-go/app/*/dao/
admin-go/app/*/internal/model/do/
admin-go/app/*/internal/model/entity/
```

## Docker 约定

- 开发 compose 在 `docker/dev/docker-compose.yml`
- 国内镜像开发 compose 在 `docker/dev/docker-compose.cn.yml`
- 生产 compose 在 `docker/prod/docker-compose.yml`
- 运行 Docker 时优先用 `docker/dev/compose.ps1` 或 `docker/dev/compose.sh`
- 两个脚本都会先把 `docker/dev/.env` 覆盖到 `admin-go/.env`
- Docker 相关资源统一放在 `docker/build/`、`docker/mysql/`、`docker/nginx/`
- 只要改 Docker 相关配置、端口、镜像、环境变量、挂载、启动命令，必须同步检查并更新这三处：`docker/dev/docker-compose.yml`、`docker/dev/docker-compose.cn.yml`、`docker/prod/docker-compose.yml`

## 铁律

铁律按主题分为 5 组。**A 组管代码生成器与前端适配，B 组管提交与数据库，C 组管 Docker，D 组管资源保护,E 组管通用约定。** 合计 10 条，任何一条都不可降级。

### A. 代码生成器与前端适配

1. **生成器是 CRUD 唯一入口**。任何重复的 CRUD 改动先看 `admin-go/codegen/`；禁止直接手改 `admin-go/app/*/dao/`、`admin-go/app/*/internal/model/do/`、`admin-go/app/*/internal/model/entity/`，**这三处只能由生成器产出**。
2. **管理端组件必须走 adapter 闭环**。codegen 生成的表单/列表/详情，只能使用 `vue-vben-admin/apps/web-antd/src/adapter/component/index.ts` 已注册、已适配的组件名与交互模式。若业务需要新组件：**先写 adapter → 再改 `admin-go/codegen/` 映射与模板 → 再生成页面**。禁止在模板里拼裸 DOM、引入仓库未注册的第三方组件、或绕过 adapter 的任何写法。
3. **vben 页面风格统一**。生成页必须优先对齐 `vue-vben-admin/apps/web-antd/src/views/system/*`、`vue-vben-admin/apps/web-antd/src/views/upload/*` 的现有写法，包括 `useVbenModal` / `useVbenForm` / `useVbenVxeGrid` 的使用顺序、打开关闭时序、导入导出方式和表单 schema 组织方式。
4. **页面风格必须可校验**。业务页风格检查统一使用 `./scripts/verify-vben-pages.sh`。修改 `vue-vben-admin/apps/web-antd/src/views/` 后必须保证这份脚本仍然通过；如果规则需要扩展，**先更新脚本，再改页面**。
5. **codegen 测试必须全面**。每次修改 `admin-go/codegen/` 后，必须同步维护离线验证脚本和覆盖样例，不允许把"需要本机跑 `npm` / `pnpm` / Docker 才能验证"当作唯一验收路径。测试表必须覆盖以下全部场景，缺一不可：
   - 树形表（parent_id）+ 非树形表
   - 单选外键（*_id 指向同应用表）+ 跨应用外键（*_id 指向其他应用表）
   - 树形外键（关联表有 parent_id）+ 普通外键（关联表无 parent_id）
   - 所有组件类型：Input、InputNumber、Textarea、Switch、Radio、Select、TreeSelect、ImageUpload、FileUpload、RichText、JsonEditor、Password、InputUrl、DateTimePicker、IconPicker
   - 枚举字段（2 值 Switch + 3 值 Radio + 多值 Select）
   - 金额字段（分→元）
   - 搜索字段（模糊 + 精确）
   - Tooltip 字段（括号提示）
   - 密码字段
   - 多段模块名（如 user_review → 包名 user_review）
   - 字段 comment 为空的回退
   - 字典字段（dict:xxx）
   - 数据权限字段（created_by + dept_id）
   - sort 排序字段
   - 自定义时间字段（*_at）
   - 验证规则（email/phone/url/max-length）

### B. 提交与数据库

6. **完成即提交，含 DB 变更必先写 migrate**。整体功能完成后必须立即提交并推送，禁止把"本地已完成但未推送"当作结束状态。**若功能包含数据库变更，必须先补齐 `admin-go/database/migrations/` 下的 `golang-migrate` 迁移文件，再执行业务提交**。默认使用仓库脚本 `./scripts/feature-publish.sh "type(scope): summary"` 基于已暂存内容提交并推送；只有明确确认全部改动都属于同一批交付时，才使用 `--all`。

### C. Docker 联动

7. **Docker 三文件联动**。任何 Docker 改动都不能只改一份 compose，必须同步检查并更新 `docker/dev/docker-compose.yml`、`docker/dev/docker-compose.cn.yml`、`docker/prod/docker-compose.yml` 三处；如果某一份必须保持差异，必须在对应改动处写清原因。

### D. 资源保护

8. **负载守卫 + 批量限流**。执行任何命令前先检查服务器负载；**只要 CPU 使用率超过 `80%` 必须立即停止，等待回落到 `50%` 以下才允许继续**。任何批量、循环、构建、迁移、脚本回填、代码生成、压测或高频重复执行任务，**一律分批推进并持续监控负载**，即便只是 `500` 次级别的执行也不得一次性打满服务器。Go / Node 重任务统一走 `scripts/run-go-task-with-limits.sh` 或 `scripts/run-node-task-with-limits.sh`。Node 限流入口默认强制 `1h` 超时、`MemoryMax=1536M`、`CPUQuota=60%`，裸 `npm` / `pnpm` / `npx` / `pnpx` / `yarn` / `corepack` 必须被 `scripts/baseadmin-shell-guard.sh` 阻断。

### E. 通用约定

9. **文档与路径**。文档统一放 `docs/`，根目录只保留简短说明；长跑巡检的滚动日志放 `docs/流程日志/`。所有引用尽量用 Markdown 相对链接，保证可以直接点击打开。
10. **变更边界**。菜单三联动：后端接口、前端入口、`system_menu` 必须一起看。本机执行边界：允许 `gf`、`go` 和数据库访问；**禁止在本机裸执行 `npm`、`pnpm` 或 Docker**，所有文档和流程设计必须遵守这条约束。确需执行前端排查时，只能使用 `scripts/run-node-task-with-limits.sh <command> [args...]`，且不得覆盖 `1h` / `1536M` 默认上限；Docker 在本机仍禁止执行。

## 常用命令

本仓库保留 Docker 入口，但当前低配机器不要在本机执行这些命令。

```powershell
.\docker\dev\compose.ps1 up -d --build
.\docker\dev\compose.ps1 down -v
.\docker\dev\compose.ps1 -China up -d --build
```

```bash
./docker/dev/compose.sh up -d --build
./docker/dev/compose.sh down -v
./docker/dev/compose.sh -China up -d --build
```
