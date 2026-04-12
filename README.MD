# GBaseAdmin

这是一个精简后的后台基础框架仓库，只保留：

- `admin-go/app/system`
- `admin-go/app/upload`
- `vue-vben-admin/apps/web-antd/src/`
- `admin-go/codegen/`

当前代码里实际保留的业务模块：

- `system`：`auth`、`dept`、`role`、`menu`、`users`，以及公开探活 `healthz` / `readyz`
- `upload`：`config`、`dir`、`dir_rule`、`file`、`uploader`，以及公开探活 `healthz` / `readyz`
- `web-antd` 页面：`views/system/{dept,menu,role,users}`、`views/upload/{config,dir,dir_rule,file}`
- 上传复用组件：`components/file-manager`、`components/upload`、`components/tinymce`、`components/json-editor`

详细说明统一放在 `docs/` 目录，新人从代码库导读开始读：

- ★ [代码库导读](docs/代码库导读.md) — 首选入口，真源表 + 架构 + 维护红线
- [文档索引](docs/README.md)
- [System 应用说明](admin-go/app/system/README.MD)
- [Upload 应用说明](admin-go/app/upload/README.MD)
- [代码生成器说明](admin-go/codegen/README.md)
- [Docker 开发说明](docs/Docker开发说明.md)
- [生产运维说明](docs/生产运维说明.md)
- [Git 提交流程铁律](docs/Git提交流程铁律.md)

Docker 开发入口已迁移到根目录：

- `docker/dev/docker-compose.yml`
- `docker/dev/docker-compose.cn.yml`
- `docker/prod/docker-compose.yml`
- `docker/dev/compose.ps1`
- `docker/dev/compose.sh`

Docker 铁律：

- 任何 Docker 改动都必须同步检查并更新 `docker/dev/docker-compose.yml`、`docker/dev/docker-compose.cn.yml`、`docker/prod/docker-compose.yml`

## Dev 启动方式

最短启动路径（从仓库根目录执行）：

```powershell
# Windows
.\docker\dev\compose.ps1 up -d --build
```

```bash
# Linux / macOS
./docker/dev/compose.sh up -d --build
```

默认只拉起后端、数据库和 adminer；前端需要显式启用 `frontend` profile。默认开发账号 `admin / admin123`，常用端口 MySQL `10020`、Redis `10021`、system `10022`、upload `10023`、frontend `10024`、adminer `10025`。

env 同步、国内镜像、profile、代理、初始化迁移等完整细节见 [Docker 开发说明](docs/Docker开发说明.md)。本地直跑（非 Docker）、常用排查命令同样集中在该文档。

线上服务管理、GitHub Actions 部署巡检、`systemd` / Nginx / 健康检查命令见 [生产运维说明](docs/生产运维说明.md)。

## Git 提交铁律

- 整体功能完成后必须立即提交并推送到 GitHub，禁止把“本地已完成但未推送”当作结束状态
- 统一使用仓库脚本：

```bash
./scripts/feature-publish.sh "feat(scope): summary"
```

- 默认只发布已经 `git add` 过的内容；工作区里如果还有未暂存或未跟踪文件，脚本会拒绝执行
- 只有确认当前所有改动都属于同一批交付时，才使用 `./scripts/feature-publish.sh --all "feat(scope): summary"`
- 如果这次功能包含数据库改动，必须先补 `admin-go/database/migrations/`，再执行提交脚本
- 仓库启用了 `lefthook`：`pre-commit` 会跑页面规范校验，`pre-push` 会跑后端测试、codegen 测试、前端 typecheck 和 `web-antd` 单测
- `scripts/run-go-task-with-limits.sh` 和 `scripts/run-node-task-with-limits.sh` 会在启动前先检查宿主机 CPU 使用率；默认高于 `80%` 时等待，降到 `50%` 后继续，可通过 `RESOURCE_LOAD_MAX_PERCENT` / `RESOURCE_LOAD_RESUME_PERCENT` 调整

## 常用校验命令

从仓库根目录执行：

```bash
# 仓库范围与页面规范
bash ./scripts/verify-baseadmin-scope.sh
bash ./scripts/verify-vben-pages.sh

# 后端与 codegen 测试（带负载守卫）
cd admin-go && ../scripts/run-go-task-with-limits.sh go test ./...
cd admin-go/codegen && ../../scripts/run-go-task-with-limits.sh go test ./...

# 前端校验（带负载守卫）
bash ./scripts/run-node-task-with-limits.sh pnpm -C vue-vben-admin -F @vben/web-antd typecheck
bash ./scripts/run-node-task-with-limits.sh pnpm -C vue-vben-admin test:web-antd
```

## Vben 组件铁律

- 管理端页面和 `admin-go/codegen` 生成模板，必须遵守本仓库 `vue-vben-admin` 组件规则
- 允许使用的表单组件和交互组件，以 `vue-vben-admin/apps/web-antd/src/adapter/component/index.ts` 中已注册、已适配的组件为准
- 如果缺组件，必须先在 adapter 完成适配，再修改 codegen 映射和模板；禁止直接在模板里生成仓库未注册的组件名
- 生成页必须优先对齐 `vue-vben-admin/apps/web-antd/src/views/system/*`、`vue-vben-admin/apps/web-antd/src/views/upload/*` 的现有写法
- 页面风格校验脚本：`./scripts/verify-vben-pages.sh`

更多说明见：

- [代码库导读](docs/代码库导读.md)
- [Git 提交流程铁律](docs/Git提交流程铁律.md)
- [Docker 开发说明](docs/Docker开发说明.md)
