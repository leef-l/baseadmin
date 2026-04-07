# GBaseAdmin

这是一个精简后的后台基础框架仓库，只保留：

- `admin-go/app/system`
- `admin-go/app/upload`
- `vue-vben-admin/apps/web-antd/src/`
- `admin-go/codegen/`

详细说明统一放在 `docs/` 目录：

- [文档索引](docs/README.md)
- [基础框架说明](docs/基础框架说明.md)
- [Docker 开发说明](docs/Docker开发说明.md)

Docker 开发入口已迁移到根目录：

- `docker/dev/docker-compose.yml`
- `docker/dev/docker-compose.cn.yml`
- `docker/prod/docker-compose.yml`
- `docker/dev/compose.ps1`
- `docker/dev/compose.sh`

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
- MySQL：`127.0.0.1:40001`
- Redis：`127.0.0.1:40002`

### 初始化数据

- `admin-go/codegen/sql/init.sql` 只保留 `system` / `upload` 当前相关表和最小种子数据
- 默认初始化账号为 `admin` / `admin123`

### 默认开发端口

- MySQL：`40001`
- Redis：`40002`
- system：`40003`
- upload：`40004`
- frontend：`40005`（容器内 `5666`）
- adminer：`40006`

### 常用排查命令

```powershell
docker compose --env-file .\docker\dev\.env -f .\docker\dev\docker-compose.yml ps
docker compose --env-file .\docker\dev\.env -f .\docker\dev\docker-compose.yml logs --tail=100
```

更多说明见：

- [Docker 开发说明](docs/Docker开发说明.md)
- [基础框架说明](docs/基础框架说明.md)
