# GBaseAdmin AI 协作说明

## 必须先阅读

1. [仓库总说明](README.MD)
2. [基础框架说明](docs/基础框架说明.md)
3. [Docker 开发说明](docs/Docker开发说明.md)
4. [system 服务路由入口](admin-go/app/system/internal/cmd/cmd.go)
5. [upload 服务路由入口](admin-go/app/upload/internal/cmd/cmd.go)

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
6. 本机资源受限：允许在本机执行 `gf`、`go` 和数据库访问；禁止在本机直接执行 `npm`、`pnpm`，也禁止在本机使用 Docker。文档和流程设计必须遵守这条约束
7. codegen 测试必须全面：每次修改 `admin-go/codegen/` 后，必须同步维护离线验证脚本和覆盖样例，不允许把“需要本机跑 `npm` / `pnpm` / Docker 才能验证”当作唯一验收路径。测试表必须覆盖以下全部场景，缺一不可：
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
