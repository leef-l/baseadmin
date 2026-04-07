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
- 只要改 Docker 相关配置、端口、镜像、环境变量、挂载、启动命令，必须同步检查并更新这三处：`docker/dev/docker-compose.yml`、`docker/dev/docker-compose.cn.yml`、`docker/prod/docker-compose.yml`

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
8. 整体功能完成后必须立即提交到 GitHub：禁止把“已完成但未提交”当作稳定状态。默认使用仓库脚本 `./scripts/feature-publish.sh "type(scope): summary"` 一次性执行 `git add -A`、`git commit`、`git push origin 当前分支`
9. 数据库铁律与提交流程同时生效：只要整体功能包含数据库变更，必须先补齐 `golang-migrate` 迁移文件，再执行功能提交流程；禁止跳过迁移直接提交业务代码
10. Docker 三文件联动：任何 Docker 改动都不能只改一份 compose，必须同步更新开发版、国内镜像开发版、生产版；如果某一份故意不同，必须在变更处写清原因
11. 管理端前端必须服从本仓库 `vue-vben-admin` 组件规则：`admin-go/codegen/` 生成的表单、列表、详情页面，只能使用 `vue-vben-admin/apps/web-antd/src/adapter/component/index.ts` 已注册、已适配的组件名与交互模式。禁止为了“先跑起来”直接在模板里拼裸 DOM 组件方案、第三方临时组件方案或绕过 adapter 的写法
12. 缺组件先适配再生成：如果业务需要的组件在 `adapter/component/index.ts` 不存在，必须先补适配组件，再修改 `codegen` 映射和模板；禁止直接在生成模板里写一个仓库未注册的组件名
13. vben 页面风格统一：生成页必须优先对齐 `vue-vben-admin/apps/web-antd/src/views/system/*`、`vue-vben-admin/apps/web-antd/src/views/upload/*` 的现有写法，包括 `useVbenModal` / `useVbenForm` / `useVbenVxeGrid` 的使用顺序、打开关闭时序、导入导出方式和表单 schema 组织方式
14. 页面风格必须可校验：业务页风格检查统一使用 `./scripts/verify-vben-pages.sh`。修改 `vue-vben-admin/apps/web-antd/src/views/` 后，至少保证这份脚本规则仍成立；如果规则需要扩展，必须先更新脚本再改页面

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
