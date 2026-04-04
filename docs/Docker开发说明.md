# Docker 开发说明

## 目录约定

- 开发环境：`docker/dev/docker-compose.yml`
- 国内镜像开发环境：`docker/dev/docker-compose.cn.yml`
- 生产环境：`docker/prod/docker-compose.yml`
- 开发环境变量源：`docker/dev/.env`
- 后端运行时环境文件：`admin-go/.env`

## 使用方式

推荐始终通过下面的脚本启动，而不是直接在 `admin-go/` 里执行 compose：

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

## env 同步规则

每次执行 `docker/dev/compose.ps1` 或 `docker/dev/compose.sh` 时，脚本都会：

1. 读取 `docker/dev/.env`
2. 覆盖复制到 `admin-go/.env`
3. 再调用对应的 `docker compose`

这样可以保证：

- Docker 使用的变量和 GoFrame 本地配置一致
- `admin-go/.env` 不需要手动维护两份
- 目录迁移后仍然只有一份开发环境源配置

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
