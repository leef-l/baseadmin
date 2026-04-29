# GBaseAdmin - System 应用

系统管理应用，提供基础的权限管理功能。

## 代码入口

- 启动入口：`main.go`
- 路由绑定：`internal/cmd/cmd.go`
- 当前挂载控制器：`auth`、`daemon`、`dept`、`domain`、`tenant`、`merchant`、`role`、`menu`、`users`、`health`、`hello`
- 默认开发端口：`10022`

## 数据库表

| 表名 | 说明 |
|------|------|
| `system_dept` | 部门表 |
| `system_role` | 角色表 |
| `system_menu` | 菜单表 |
| `system_users` | 用户表 |
| `system_user_role` | 用户-角色关联表 |
| `system_role_menu` | 角色-菜单关联表 |
| `system_role_dept` | 角色-部门关联表 |
| `system_user_dept` | 用户-部门关联表 |
| `system_tenant` | 租户表 |
| `system_merchant` | 商户表 |
| `system_domain` | 域名绑定表 |
| `system_daemon` | 守护进程配置表 |

## API 路由

对外接口分两层：

- 公开路由：`GET /healthz`、`GET /readyz`、`GET /hello`
- 业务接口前缀：`/api/system/`

`web-antd` 内部请求路径会以 `/api` 作为 baseURL，再访问 `/system/...`。

| 模块 | 路径 | 鉴权 | 说明 |
|------|------|------|------|
| 健康检查 | `GET /healthz` | 否 | 返回 `{ status: "ok", service: "system" }` |
| 就绪检查 | `GET /readyz` | 否 | 返回数据库与关键依赖检查结果 |
| 示例 | `GET /hello` | 否 | GoFrame 默认 hello 路由，当前仍挂载 |
| 认证 | `POST /api/system/auth/login` | 否 | 用户名密码登录 |
| 认证 | `POST /api/system/auth/ticket-login` | 否 | 应用间票据登录 |
| 认证 | `POST /api/system/auth/ticket` | 是 | 为当前登录用户生成目标应用票据 |
| 认证 | `GET /api/system/auth/info` | 是 | 当前用户信息，返回 `roles` 和 `perms` |
| 认证 | `GET /api/system/auth/menus` | 是 | 当前用户菜单树，供 `web-antd` 动态路由使用 |
| 认证 | `POST /api/system/auth/change-password` | 是 | 修改当前用户密码 |
| 认证 | `POST /api/system/auth/logout` | 是 | 退出登录，Token 加入黑名单 |
| 租户 | `/api/system/tenant/list,create,update,delete,batch-delete,detail` | 是 | 租户管理，可创建租户管理员 |
| 商户 | `/api/system/merchant/list,create,update,delete,batch-delete,detail` | 是 | 商户管理，可创建商户管理员 |
| 域名 | `/api/system/domain/list,create,update,delete,batch-delete,detail` | 是 | 租户/商户自定义域名管理 |
| 域名 | `/api/system/domain/apply-nginx,apply-ssl` | 是 | 写入宝塔 Nginx 配置、申请 SSL 证书 |
| 守护进程 | `/api/system/daemon/list,create,update,delete,batch-delete,detail` | 是 | 宝塔 Supervisor 守护进程管理 |
| 守护进程 | `/api/system/daemon/restart,batch-restart,stop,batch-stop,log` | 是 | 重启、暂停、查看日志 |
| 部门 | `/api/system/dept/list,create,update,delete,batch-delete,detail,tree` | 是 | 部门管理 |
| 角色 | `/api/system/role/list,create,update,delete,batch-delete,detail,tree` | 是 | 角色管理 |
| 角色 | `/api/system/role/grant-menu,menu-ids,grant-dept,dept-ids` | 是 | 菜单授权和数据权限 |
| 菜单 | `/api/system/menu/list,create,update,delete,batch-delete,detail,tree` | 是 | 菜单管理 |
| 用户 | `/api/system/users/list,create,update,delete,batch-delete,detail,reset-password` | 是 | 用户管理和密码重置 |

## 权限模型

- 角色设置 `is_admin=1` 为超级管理员，跳过权限判断，可查看所有菜单和数据
- 普通角色通过 `system_role_menu` 关联表控制菜单权限
- 数据权限通过 `system_role_dept` 关联表控制（`data_scope=5` 时生效）
- `middleware/auth.go` 会在服务端校验按钮权限，不再只依赖前端隐藏按钮；同时检查 Token 黑名单，已注销 Token 会被拒绝
- 数据权限已经进入运行时查询和写操作校验，当前覆盖用户、部门、角色、租户/商户、域名、守护进程等系统管理入口；平台级专用能力仍会额外校验平台超级管理员
- 用户分配角色时，服务端会再次校验“当前操作者可分配的启用角色”；非超级管理员只能分配自己当前持有的启用角色
- 角色数据权限授权会按当前操作者的数据范围再次校验自定义部门列表，不能只靠前端树控件限制
- JWT 内携带 `tenant_id` / `merchant_id`，中间件和业务逻辑用它判断当前账号的 SaaS 数据范围
- 系统管理里仍需保留 `tenant_id`、`merchant_id`、`created_by`、`dept_id` 四个字段；平台级数据统一写 `0/0`
- `tenant_id=0`、`merchant_id=0` 表示平台级；`tenant_id>0`、`merchant_id=0` 表示租户级；两者都大于 0 表示商户级
- 租户和商户不是登录主体；真正登录后台的仍然是 `system_users`
- 租户/商户创建时可自动创建管理员用户、部门、角色和菜单授权
- `DomainContext` 会按访问域名匹配 `system_domain`，并阻止不属于该域名归属的账号访问

## SaaS 模型

当前是“一张用户表登录 + 租户/商户主体表建模 + 角色授权 + 部门数据权限”的模型：

| 层级 | 数据归属 | 登录账号 |
|------|----------|----------|
| 平台 | `tenant_id=0, merchant_id=0` | 平台管理员用户 |
| 租户 | `tenant_id=租户ID, merchant_id=0` | 租户管理员或租户员工 |
| 商户 | `tenant_id=租户ID, merchant_id=商户ID` | 商户管理员或商户员工 |

创建租户或商户时，如果选择创建管理员，后端会在事务里创建部门、角色、用户和角色菜单授权。后续也可以通过用户管理手工创建用户，再绑定租户/商户、部门和角色。

详细关系见仓库级文档 [SaaS 多租户与域名说明](../../../docs/SaaS多租户与域名说明.md)。

## 域名与守护进程

域名模块：

- 表：`system_domain`
- 页面：`views/system/domain/index.vue`
- Nginx 配置目录：`/www/server/panel/vhost/nginx`
- 证书目录：`/www/server/panel/vhost/cert/{domain}`
- SSL 申请通过宝塔 Python 环境调用 ACME 能力

守护进程模块：

- 表：`system_daemon`
- 页面：`views/system/daemon/index.vue`
- Supervisor profile：`/www/server/panel/plugin/supervisor/profile/{program}.ini`
- 宝塔 Supervisor 配置：`/www/server/panel/plugin/supervisor/config.json`
- 日志目录：`/www/server/panel/plugin/supervisor/log`

守护进程属于高危能力，只应授权给超级管理员或明确可信的运维角色。详细说明见 [守护进程管理说明](../../../docs/守护进程管理说明.md)。

## 菜单与前端契约

- `auth/info` 同时驱动前端用户信息和权限码加载
- 如果多个应用共用登录态，来源应用可以先在自己的登录态下调用 `POST /api/system/auth/ticket` 生成短效签名票据，再跳转目标应用的 `/auth/login?ticket=...`
- 目标应用收到 `ticket` 后，前端会自动调用 `POST /api/system/auth/ticket-login` 换取当前后台 JWT
- 参与互通的应用要使用同一份 `sso.sharedKey`，并各自配置不同的 `sso.appId`
- `auth/menus` 返回的是 `web-antd` 后端权限模式所需的菜单树
- `type=3` 按钮权限不会进入前端路由树
- `type=4` 外链会在前端映射为 `meta.link`
- `type=5` 内链会在前端映射为 `meta.iframeSrc`
- 菜单 `linkURL` / `linkUrl` 两种字段名，前端都会兼容读取
- 外链 / 内嵌页在后端未填写 `component` 时，前端会兜底为 `IFrameView`
- 菜单里的 `path`、`component`、`icon`、`permission` 必须与 `vue-vben-admin/apps/web-antd/src/views` 和路由约定保持一致
- 后端会按菜单类型强制校验 `path / component / permission / linkURL`，并拦截重复路由路径；前端表单也会同步做类型化校验

## 对应前端页面

| 页面 | 说明 |
|------|------|
| `views/system/dept/index.vue` | 部门树和部门表单 |
| `views/system/domain/index.vue` | 自定义域名、Nginx 应用、SSL 申请 |
| `views/system/daemon/index.vue` | 守护进程创建、编辑、重启、暂停、日志 |
| `views/system/menu/index.vue` | 菜单树、菜单类型、外链/内嵌页配置 |
| `views/system/merchant/index.vue` | 商户 CRUD 和商户管理员初始化 |
| `views/system/role/index.vue` | 角色 CRUD、菜单授权、数据权限授权 |
| `views/system/tenant/index.vue` | 租户 CRUD 和租户管理员初始化 |
| `views/system/users/index.vue` | 用户 CRUD、部门筛选、角色分配、重置密码 |

## 本地直跑

```bash
cd admin-go
go run ./app/system/main.go
```
