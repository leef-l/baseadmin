# GBaseAdmin AI 协作说明

## 必须先阅读

1. [EasyMVP对接OpenHands与Aider引擎设计实现文档](docs/EasyMVP对接OpenHands与Aider引擎设计实现文档.md)
2. [仓库总说明](README.MD)
3. [基础框架说明](docs/基础框架说明.md)
4. [Docker 开发说明](docs/Docker开发说明.md)
5. [system 服务路由入口](admin-go/app/system/internal/cmd/cmd.go)
6. [upload 服务路由入口](admin-go/app/upload/internal/cmd/cmd.go)

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

## 铁律

1. 生成器优先：CRUD 重复问题先看 `admin-go/codegen/`
2. 不要手改生成层：`dao/do/entity`
3. 菜单与权限联动：后端接口、前端入口、`system_menu` 要一起看
4. 文档统一放 `docs/`，根目录只保留简短说明
5. 路径引用尽量写成 Markdown 相对链接，保证能直接点击打开

## 常用命令

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
