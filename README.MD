# GBaseAdmin

这是一个精简后的后台基础框架仓库，只保留：

- `admin-go/app/system`
- `admin-go/app/upload`
- `vue-vben-admin/apps/web-antd/src/`
- `admin-go/codegen/`

详细说明统一放在 `docs/` 目录：

- [文档索引](docs/README.md)
- [Git 提交流程铁律](docs/Git提交流程铁律.md)
- [基础框架说明](docs/基础框架说明.md)
- [Docker 开发说明](docs/Docker开发说明.md)

Docker 开发入口已迁移到根目录：

- `docker/dev/docker-compose.yml`
- `docker/dev/docker-compose.cn.yml`
- `docker/prod/docker-compose.yml`
- `docker/dev/compose.ps1`
- `docker/dev/compose.sh`

Docker 铁律：

- 任何 Docker 改动都必须同步检查并更新 `docker/dev/docker-compose.yml`、`docker/dev/docker-compose.cn.yml`、`docker/prod/docker-compose.yml`

## Dev 启动方式

### 启动前提

- 从仓库根目录执行命令
- 具备 Docker / Docker Compose 环境
- 默认配置文件是 `docker/dev/.env`
- 每次执行脚本都会同步覆盖 `admin-go/.env`

### Windows

开发环境：

```powershell
.\docker\dev\compose.ps1 up -d --build
```

国内镜像：

```powershell
.\docker\dev\compose.ps1 -China up -d --build
```

停止并删除容器/卷：

```powershell
.\docker\dev\compose.ps1 down -v
```

### Linux / macOS

开发环境：

```bash
./docker/dev/compose.sh up -d --build
```

国内镜像：

```bash
./docker/dev/compose.sh -China up -d --build
```

停止并删除容器/卷：

```bash
./docker/dev/compose.sh down -v
```

### 脚本会自动做什么

- 读取 `docker/dev/.env`
- 同步覆盖到 `admin-go/.env`
- 自动选择 `docker/dev/docker-compose.yml` 或 `docker/dev/docker-compose.cn.yml`
- 执行 `docker compose up -d --build`

### 默认开发账号和连接

- 登录账号：`admin`
- 登录密码：`admin123`
- MySQL：`127.0.0.1:44001`
- Redis：`127.0.0.1:44002`

### 初始化数据

- 初始化现在统一走 `admin-go/database/migrations/`
- Docker 与部署启动前都会先执行 `golang-migrate up`
- `DB_MIGRATE_AUTO=0/false/no/off` 可显式关闭容器启动迁移；`auto` 下默认仅 `system` 自动执行
- 默认初始化账号为 `admin` / `admin123`

### 默认开发端口

- MySQL：`44001`
- Redis：`44002`
- system：`44003`
- upload：`44004`
- frontend：`44005`（容器内 `5666`）
- adminer：`44006`

### 常用排查命令

```powershell
docker compose --env-file .\docker\dev\.env -f .\docker\dev\docker-compose.yml ps
docker compose --env-file .\docker\dev\.env -f .\docker\dev\docker-compose.yml logs --tail=100
```

## Git 提交铁律

- 整体功能完成后必须立即提交并推送到 GitHub，禁止把“本地已完成但未推送”当作结束状态
- 统一使用仓库脚本：

```bash
./scripts/feature-publish.sh "feat(scope): summary"
```

- 如果这次功能包含数据库改动，必须先补 `admin-go/database/migrations/`，再执行提交脚本

## Vben 组件铁律

- 管理端页面和 `admin-go/codegen` 生成模板，必须遵守本仓库 `vue-vben-admin` 组件规则
- 允许使用的表单组件和交互组件，以 `vue-vben-admin/apps/web-antd/src/adapter/component/index.ts` 中已注册、已适配的组件为准
- 如果缺组件，必须先在 adapter 完成适配，再修改 codegen 映射和模板；禁止直接在模板里生成仓库未注册的组件名
- 生成页必须优先对齐 `vue-vben-admin/apps/web-antd/src/views/system/*`、`vue-vben-admin/apps/web-antd/src/views/upload/*` 的现有写法
- 页面风格校验脚本：`./scripts/verify-vben-pages.sh`

更多说明见：

- [Git 提交流程铁律](docs/Git提交流程铁律.md)
- [Docker 开发说明](docs/Docker开发说明.md)
- [基础框架说明](docs/基础框架说明.md)
