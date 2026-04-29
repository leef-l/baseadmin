# GBaseAdmin - Demo 应用

Codegen 生成器的全场景覆盖示例。所有代码均由 `admin-go/codegen` 生成，用于验证生成器在树形、外键、枚举、字典、富组件、密码、金额、时间范围等场景下的正确性。

Demo 只用于开发验证和体验测试，不参与生产部署。

## 代码入口

- 启动入口：`main.go`
- 路由绑定：`internal/cmd/cmd.go`
- 当前挂载控制器：10 个模块，全部为认证接口
- 默认开发端口：`10026`
- 路由前缀：`/api/demo`

## 数据库表

| 表名 | 说明 | 覆盖场景 |
|------|------|----------|
| `demo_category` | 分类管理 | 树形表（parent_id） |
| `demo_customer` | 客户管理 | 跨应用外键（ref:system_users）、枚举、图片上传 |
| `demo_product` | 产品管理 | 金额字段、图片上传、富文本、JSON 编辑器 |
| `demo_campaign` | 营销活动 | 枚举、时间字段、Textarea |
| `demo_order` | 订单管理 | 金额、多枚举、跨应用外键 |
| `demo_work_order` | 工单管理 | 多枚举、文件上传 |
| `demo_contract` | 合同管理 | 文件上传、金额、时间字段 |
| `demo_survey` | 问卷调查 | JSON 编辑器、枚举 |
| `demo_appointment` | 预约管理 | 时间字段、枚举、跨应用外键 |
| `demo_audit_log` | 审计日志 | JSON 编辑器、IP 字段、InputUrl |

10 张表由迁移 `000015_add_rich_demo_app` 创建。

## 与 system/upload 的区别

- 无 health 控制器（无 `/healthz`、`/readyz` 端点）
- 无 `manifest/config/config.yaml`，开发配置在 `hack/config.yaml`
- 无独立的中间件链（RequestID、DomainContext、AccessLog 等均未挂载）
- 不参与 `deploy.yml` 生产部署，Docker compose 中作为开发服务存在

## 对应前端页面

| 页面 | 说明 |
|------|------|
| `views/demo/appointment/index.vue` | 预约管理 |
| `views/demo/audit_log/index.vue` | 审计日志 |
| `views/demo/campaign/index.vue` | 营销活动 |
| `views/demo/category/index.vue` | 分类管理（树形） |
| `views/demo/contract/index.vue` | 合同管理 |
| `views/demo/customer/index.vue` | 客户管理 |
| `views/demo/order/index.vue` | 订单管理 |
| `views/demo/product/index.vue` | 产品管理 |
| `views/demo/survey/index.vue` | 问卷调查 |
| `views/demo/work_order/index.vue` | 工单管理 |

## 本地直跑

```bash
cd admin-go
go run ./app/demo/main.go
```

## 重新生成

如需重新生成 demo 应用（例如修改了 codegen 模板后验证输出）：

```bash
cd admin-go/codegen
../../scripts/run-go-task-with-limits.sh go run . --table demo_category,demo_customer,demo_product,demo_campaign,demo_order,demo_work_order,demo_contract,demo_survey,demo_appointment,demo_audit_log --force
```
