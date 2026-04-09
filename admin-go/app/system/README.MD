# GBaseAdmin - System 应用

系统管理应用，提供基础的权限管理功能。

## 代码入口

- 启动入口：`main.go`
- 路由绑定：`internal/cmd/cmd.go`
- 当前挂载控制器：`auth`、`dept`、`role`、`menu`、`users`、`health`、`hello`
- 默认开发端口：`44003`

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
| 认证 | `POST /api/system/auth/login` | 否 | 登录 |
| 认证 | `GET /api/system/auth/info` | 是 | 当前用户信息，返回 `roles` 和 `perms` |
| 认证 | `GET /api/system/auth/menus` | 是 | 当前用户菜单树，供 `web-antd` 动态路由使用 |
| 认证 | `POST /api/system/auth/change-password` | 是 | 修改当前用户密码 |
| 部门 | `/api/system/dept/list,create,update,delete,batch-delete,detail,tree` | 是 | 部门管理 |
| 角色 | `/api/system/role/list,create,update,delete,batch-delete,detail,tree` | 是 | 角色管理 |
| 角色 | `/api/system/role/grant-menu,menu-ids,grant-dept,dept-ids` | 是 | 菜单授权和数据权限 |
| 菜单 | `/api/system/menu/list,create,update,delete,batch-delete,detail,tree` | 是 | 菜单管理 |
| 用户 | `/api/system/users/list,create,update,delete,batch-delete,detail,reset-password` | 是 | 用户管理和密码重置 |

## 权限模型

- 角色设置 `is_admin=1` 为超级管理员，跳过权限判断，可查看所有菜单和数据
- 普通角色通过 `system_role_menu` 关联表控制菜单权限
- 数据权限通过 `system_role_dept` 关联表控制（`data_scope=5` 时生效）
- `middleware/auth.go` 会在服务端校验按钮权限，不再只依赖前端隐藏按钮
- 数据权限已经进入运行时查询和写操作校验，当前至少覆盖 `users` / `dept` 的列表、详情及主要写接口

## 菜单与前端契约

- `auth/info` 同时驱动前端用户信息和权限码加载
- `auth/menus` 返回的是 `web-antd` 后端权限模式所需的菜单树
- `type=3` 按钮权限不会进入前端路由树
- `type=4` 外链会在前端映射为 `meta.link`
- `type=5` 内链会在前端映射为 `meta.iframeSrc`
- 菜单 `linkURL` / `linkUrl` 两种字段名，前端都会兼容读取
- 外链 / 内嵌页在后端未填写 `component` 时，前端会兜底为 `IFrameView`
- 菜单里的 `path`、`component`、`icon`、`permission` 必须与 `vue-vben-admin/apps/web-antd/src/views` 和路由约定保持一致

## 对应前端页面

| 页面 | 说明 |
|------|------|
| `views/system/dept/index.vue` | 部门树和部门表单 |
| `views/system/menu/index.vue` | 菜单树、菜单类型、外链/内嵌页配置 |
| `views/system/role/index.vue` | 角色 CRUD、菜单授权、数据权限授权 |
| `views/system/users/index.vue` | 用户 CRUD、部门筛选、角色分配、重置密码 |

## 本地直跑

```bash
cd admin-go
go run ./app/system/main.go
```
