# Docker 开发说明

本说明描述仓库保留的 Docker 方案，适用于具备 Docker 条件的环境。低配机器建议保持 backend-only，不要默认拉起前端容器。

## 目录约定

- 开发环境：`docker/dev/docker-compose.yml`
- 国内镜像开发环境：`docker/dev/docker-compose.cn.yml`
- 生产环境：`docker/prod/docker-compose.yml`
- 开发环境变量源：`docker/dev/.env`
- 后端运行时环境文件：`admin-go/.env`

## 铁律

- 任何 Docker 改动都必须同时检查并更新这三处 compose：
  - `docker/dev/docker-compose.yml`
  - `docker/dev/docker-compose.cn.yml`
  - `docker/prod/docker-compose.yml`
- 如果其中某一份必须保持差异，必须在对应改动处写清楚原因，不能默认遗漏

## 使用方式

推荐始终通过下面的脚本启动，而不是直接在 `admin-go/` 里执行 compose：

```powershell
.\docker\dev\compose.ps1 up -d --build
.\docker\dev\compose.ps1 down -v
.\docker\dev\compose.ps1 -China up -d --build
.\docker\dev\compose.ps1 --profile frontend up -d --build frontend
```

```bash
./docker/dev/compose.sh up -d --build
./docker/dev/compose.sh down -v
./docker/dev/compose.sh -China up -d --build
./docker/dev/compose.sh --profile frontend up -d --build frontend
```

默认启动只包含后端、数据库和 adminer；`frontend` 改为显式 profile，避免低配机器在容器启动时自动执行 `pnpm install`。
`compose.sh` / `compose.ps1` 在检测到 `frontend` 启动请求时，会先校验宿主机可用内存，默认要求至少 `2048MB`；不足时会直接拒绝启动。确需强制启动时可临时设置 `ALLOW_LOW_MEMORY_FRONTEND=1`，也可通过 `FRONTEND_MIN_HOST_MEM_MB` 调整阈值。

## 前端代理约定

- `web-antd` 开发环境 `VITE_GLOB_API_URL=/api`
- `docker/dev/.env` 默认注入 `VITE_PROXY_SYSTEM_TARGET=http://system:10022`
- `docker/dev/.env` 默认注入 `VITE_PROXY_UPLOAD_TARGET=http://upload:10023`
- `docker/dev/.env` 默认注入 `VITE_PROXY_DEMO_TARGET=http://demo:10026`
- 如果不是走 Docker 前端容器，而是宿主机直跑 `web-antd`，`apps/web-antd/.env.development` 默认把 `/api/system`、`/api/upload`、`/api/demo` 分别转发到对应后端端口
- 前端容器内端口为 `10024`，compose 暴露端口为 `10024`

## env 同步规则

每次执行 `docker/dev/compose.ps1` 或 `docker/dev/compose.sh` 时，脚本都会：

1. 读取 `docker/dev/.env`
2. 覆盖复制到 `admin-go/.env`
3. 再调用对应的 `docker compose`

这样可以保证：

- Docker 使用的变量和 GoFrame 本地配置一致
- `admin-go/.env` 不需要手动维护两份
- 目录迁移后仍然只有一份开发环境源配置

## 默认连接与账号

- 登录账号：`admin`
- 登录密码：`admin123`
- MySQL：`127.0.0.1:10020`
- Redis：`127.0.0.1:10021`
- system：`127.0.0.1:10022`
- upload：`127.0.0.1:10023`
- frontend：`127.0.0.1:10024`（需显式启用 `frontend` profile）
- adminer：`127.0.0.1:10025`
- demo：`127.0.0.1:10026`（codegen 全场景示例应用）
- `system` 健康检查：`GET /healthz`
- `system` 就绪检查：`GET /readyz`
- `upload` 健康检查：`GET /healthz`
- `upload` 就绪检查：`GET /readyz`

## 初始化来源

- MySQL 容器不再挂载初始化 SQL
- `system` / `upload` / `demo` 容器启动前会自动执行 `golang-migrate up`
- `DB_MIGRATE_AUTO=0/false/no/off` 可关闭启动迁移；`auto` 模式默认仅 `system` 执行迁移
- 数据库真源迁移目录为 [admin-go/database/migrations](../admin-go/database/migrations)
- 默认上传配置使用本地存储，目录为 `resource/upload`

## Docker 资源目录

Docker 相关资源现在统一放在：

- Dockerfile
- MySQL 初始化 SQL
- Nginx 配置

对应目录分别是：

- `docker/build/`
- `docker/mysql/`
- `docker/nginx/`

compose 入口也统一在根目录 `docker/` 下。
