# SaaS 多租户、商户与域名说明

本文说明当前项目里的平台、租户、商户、用户、角色、部门、域名之间的关系，以及后台如何创建和管理这些对象。

## 核心模型

当前项目采用“一张用户表 + 租户/商户运营主体表 + 角色权限”的 SaaS 模型。

| 对象 | 表 | 含义 | 是否登录主体 |
|------|----|------|--------------|
| 平台 | `tenant_id=0`、`merchant_id=0` | 系统运营方，超级管理员所在层级 | 否 |
| 租户 | `system_tenant` | SaaS 客户主体，例如一家企业、一个品牌或一个站点运营方 | 否 |
| 商户 | `system_merchant` | 租户下面的业务主体，例如门店、子品牌、加盟商 | 否 |
| 用户 | `system_users` | 真正登录后台的人 | 是 |
| 角色 | `system_role` | 权限集合，决定菜单、按钮和数据权限 | 否 |
| 部门 | `system_dept` | 组织树和部门数据权限载体 | 否 |
| 域名 | `system_domain` | 自定义域名和访问域名绑定 | 否 |

重点：

- 租户和商户不是登录账号，它们是运营主体。
- 登录后台的永远是 `system_users` 里的用户。
- 用户通过 `tenant_id` / `merchant_id` 归属到平台、租户或商户。
- 用户通过 `system_user_role` 绑定角色，角色通过 `system_role_menu` 获得菜单和按钮权限。
- 部门数据权限仍然存在；租户/商户权限负责 SaaS 主体隔离，部门权限负责组织范围隔离。

## 三层关系

```text
平台
├── 平台管理员用户
├── 租户 A
│   ├── 租户管理员用户
│   ├── 租户角色 / 部门
│   ├── 商户 A1
│   │   ├── 商户管理员用户
│   │   └── 商户角色 / 部门
│   └── 商户 A2
└── 租户 B
```

归属字段含义：

| `tenant_id` | `merchant_id` | 含义 |
|-------------|---------------|------|
| `0` | `0` | 平台级数据 |
| `租户ID` | `0` | 租户级数据 |
| `租户ID` | `商户ID` | 商户级数据 |

平台超级管理员没有租户限制；租户用户只能访问自己的租户；商户用户只能访问自己的租户和商户。

## 用户、角色、部门如何配合

当前没有为租户或商户创建独立登录表，原因是后台登录、权限、菜单、密码、JWT、数据权限都已经围绕 `system_users` 工作。

一个租户或商户需要登录后台时，会有对应用户：

- 租户管理员：`system_users.tenant_id = 租户ID`，`merchant_id = 0`
- 商户管理员：`system_users.tenant_id = 租户ID`，`merchant_id = 商户ID`
- 普通租户/商户员工：同样在 `system_users`，按实际归属绑定角色

角色也是按归属隔离：

- 平台角色：`tenant_id=0`、`merchant_id=0`
- 租户角色：`tenant_id=租户ID`、`merchant_id=0`
- 商户角色：`tenant_id=租户ID`、`merchant_id=商户ID`

部门用于组织树和部门数据权限，不替代租户/商户：

- 部门解决“这个用户在组织内能看哪些部门的数据”
- 租户/商户解决“这个用户属于哪个 SaaS 运营主体”
- 两者同时存在时，查询会先按租户/商户收口，再按部门数据权限收口

## 创建租户

入口：

- 页面：`/system/tenant`
- API：`POST /api/system/tenant/create`
- 权限：`system:tenant:create`
- 后端：`admin-go/app/system/internal/logic/tenant/tenant.go`

租户只能由平台侧账号创建。租户创建时可以勾选“创建管理员”，后端会在同一个事务里自动创建：

1. `system_tenant` 租户记录
2. 租户根部门
3. 租户管理员角色
4. 租户管理员用户
5. 用户角色绑定
6. 租户管理员角色默认菜单权限

自动管理员默认规则：

- 未填写管理员用户名时，会用 `{租户编码}_admin`
- 管理员密码走统一密码策略和 bcrypt 加密
- 租户管理员默认可拿到部门、角色、用户、域名、商户等租户运营所需菜单
- 域名的 Nginx 应用和 SSL 申请这类高风险按钮不会默认下发给租户管理员

如果创建租户时没有创建管理员，也可以后续在“用户管理”里手工创建用户，并选择对应租户、角色和部门。

## 创建商户

入口：

- 页面：`/system/merchant`
- API：`POST /api/system/merchant/create`
- 权限：`system:merchant:create`
- 后端：`admin-go/app/system/internal/logic/merchant/merchant.go`

商户必须归属于租户。创建商户时同样可以勾选“创建管理员”，后端会在同一个事务里自动创建：

1. `system_merchant` 商户记录
2. 商户部门
3. 商户管理员角色
4. 商户管理员用户
5. 用户角色绑定
6. 商户管理员角色默认菜单权限

商户管理员默认菜单范围比租户管理员更小：

- 可管理自己范围内的部门、角色、用户、域名
- 不默认拥有商户管理权限
- 不默认拥有域名 Nginx 应用和 SSL 申请按钮

租户账号创建商户时，后端会强制把商户归属写成当前账号所属租户，不能通过请求参数越权创建到别的租户。

## 数据权限公共方法

租户/商户数据权限的公共方法位于：

- `admin-go/app/system/internal/logic/shared/tenant_scope.go`
- `admin-go/app/system/internal/middleware/context_scope.go`
- `admin-go/codegen/templates/backend/middleware_context.tpl`

核心方法：

| 方法 | 作用 |
|------|------|
| `ResolveTenantAccessScope` | 从 JWT 上下文解析当前账号的租户/商户范围 |
| `ApplyTenantScopeToModel` | 给列表、树形、导出查询追加租户/商户过滤 |
| `ApplyTenantScopeToWrite` | 租户/商户账号写入时覆盖请求里的归属字段 |
| `EnsureTenantMerchantAccessible` | 校验目标租户/商户存在、启用且可访问 |
| `EnsureTenantScopedRowAccessible` | 单条数据写操作前校验行归属 |
| `EnsureTenantScopedRowsAccessible` | 批量写操作前校验多行归属 |

codegen 生成的 CRUD 模板已经按表字段自动接入：

| 场景 | 生成调用 |
|------|----------|
| 新增 | `ApplyTenantScopeToWrite` + `EnsureTenantMerchantAccessible` |
| 更新 | 写入归属覆盖 + 目标归属校验 + 行级归属校验 |
| 删除 | 行级归属校验 |
| 批量删除 | 批量行级归属校验 |
| 详情 | 行级归属校验 |
| 列表 | `ApplyTenantScopeToModel` |
| 树形 | `ApplyTenantScopeToModel` |
| 导出 | `ApplyTenantScopeToModel` |
| 批量更新 | 批量行级归属校验 |
| 导入 | 写入归属覆盖 + 目标归属校验 |

触发规则：

- codegen CRUD 表必须同时有 `tenant_id` 和 `merchant_id`：启用租户/商户权限
- codegen CRUD 表必须同时有 `created_by` 和 `dept_id`：继续启用部门数据权限
- `tenant_id`、`merchant_id`、`created_by`、`dept_id` 缺任意一个都会在生成阶段失败
- 后台列表页里，租户/商户归属列和筛选控件只给平台超级管理员显示；租户、商户和普通员工账号不展示这些字段。

## 域名绑定

入口：

- 页面：`/system/domain`
- API：
  - `POST /api/system/domain/create`
  - `POST /api/system/domain/apply-nginx`
  - `POST /api/system/domain/apply-ssl`
- 权限：
  - `system:domain:create/update/delete/list/batch-delete`
  - `system:domain:apply`
  - `system:domain:ssl`

域名表：`system_domain`

关键字段：

| 字段 | 含义 |
|------|------|
| `domain` | 绑定域名 |
| `owner_type` | `1=租户`、`2=商户` |
| `tenant_id` | 归属租户 |
| `merchant_id` | 归属商户，`0` 表示租户级域名 |
| `app_code` | 当前主要使用 `admin` |
| `verify_status` | `0=未校验`、`1=已校验` |
| `nginx_status` | `0=未应用`、`1=已应用` |
| `ssl_status` | `0=未配置`、`1=已配置` |
| `status` | `0=关闭`、`1=开启` |

当前流程：

1. 在域名 DNS 里把域名或泛解析指向当前服务器
2. 后台创建域名绑定
3. 点击“应用 Nginx 配置”
4. 后端写入宝塔 Nginx 配置并执行 `nginx -t`、`nginx -s reload`
5. 点击“申请 SSL 证书”
6. 后端调用宝塔 Python 环境的 ACME 能力申请证书
7. 证书申请成功后再次生成带 SSL 的 Nginx 配置并 reload

实际文件位置：

| 内容 | 路径 |
|------|------|
| Nginx vhost | `/www/server/panel/vhost/nginx/baseadmin-managed-{domain}.conf` |
| 证书目录 | `/www/server/panel/vhost/cert/{domain}/` |
| 证书文件 | `fullchain.pem`、`privkey.pem` |
| ACME challenge | `/www/wwwroot/baseadmin.easytestdev.online/admin/.well-known/acme-challenge/` |

域名访问时，`DomainContext` 中间件会读取 `Host` / `X-Forwarded-Host`：

1. 查找 `system_domain`
2. 要求 `app_code=admin`、`verify_status=1`、`status=1`
3. 匹配后把域名归属写入请求上下文
4. `Auth` 中间件校验当前 JWT 的租户/商户归属是否允许访问该域名

因此，一个租户域名不能被其他租户账号拿着 JWT 直接访问。

## 菜单与权限

迁移会自动写入这些菜单：

- `/system/tenant`：租户管理
- `/system/merchant`：商户管理
- `/system/domain`：域名管理

默认授权给基线超级管理员角色。租户/商户自动管理员会按菜单 profile 获得一组受限菜单，后续可在角色管理里调整。

如果页面 404，优先排查：

1. `system_menu` 是否有对应页面菜单
2. `system_role_menu` 是否给当前用户角色授权
3. `GET /api/system/auth/info` 是否返回按钮权限码
4. `GET /api/system/auth/menus` 是否返回页面路由
5. 前端 `accessCodes` 是否已刷新，旧登录态可能需要重新登录

## 常见问题

### 为什么不只用一个用户表加角色？

当前确实仍然只有一个用户表作为登录主体。但租户和商户是运营主体，不是用户身份本身。单靠角色无法表达：

- 租户/商户生命周期、联系人、编码、到期时间
- 商户必须归属于哪个租户
- 自定义域名归属到租户还是商户
- 数据行必须按租户/商户强隔离
- 创建租户/商户时自动初始化管理员、部门和角色

所以当前是“一张用户表登录 + 租户/商户主体表建模 + 角色授权 + 部门数据权限”的组合。

### 租户或商户怎么进后台？

租户或商户本身不能登录，必须有对应用户。创建租户/商户时勾选“创建管理员”可以自动生成后台用户；也可以后续在用户管理中手工创建并绑定租户/商户、角色和部门。

### 部门数据权限还保留吗？

保留。部门权限解决组织范围，租户/商户权限解决 SaaS 主体范围，两者是不同维度。不要用部门表硬模拟租户/商户，否则自定义域名、商户归属、自动初始化和跨租户隔离都会变复杂。
