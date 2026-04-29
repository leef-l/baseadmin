# 文档目录

`docs/` 只保留当前仍长期维护的项目说明。一次性协作提示、临时排障记录和过期巡检日志不再放在这里。

## 首选阅读

- ★ [代码库导读](代码库导读.md) — 新人从这里开始：真源表 + 架构地图 + 日常维护红线

## 项目总览

- [System 应用说明](../admin-go/app/system/README.MD)
- [Upload 应用说明](../admin-go/app/upload/README.MD)
- [Demo 应用说明](../admin-go/app/demo/README.MD) — codegen 全场景示例
- [SaaS 多租户与域名说明](SaaS多租户与域名说明.md)
- [守护进程管理说明](守护进程管理说明.md)
- [代码生成器说明](../admin-go/codegen/README.md)
- [Codegen AI 执行手册](../admin-go/codegen/docs/Codegen-AI执行手册.md)

## 安全与认证

- Token 黑名单、登出流程、登录限流、时序攻击防护、域名严格模式 → 见 [代码库导读](代码库导读.md) 认证主链路部分
- Redis Token 黑名单运维 → 见 [生产运维说明](生产运维说明.md)

## 开发与部署

- [Docker 开发说明](Docker开发说明.md)
- [生产运维说明](生产运维说明.md)

## 协作约定

- [AI 协作铁律](../CLAUDE.md)
- [Git 提交流程铁律](Git提交流程铁律.md)

## AI 执行入口

- Codegen 任务先读 [代码生成器说明](../admin-go/codegen/README.md)、[字段备注与生成规则](../admin-go/codegen/docs/字段备注与生成规则.md)，再读 [Codegen AI 执行手册](../admin-go/codegen/docs/Codegen-AI执行手册.md)
- 执行 codegen、Go 命令、前端校验时必须遵守 [AI 协作铁律](../CLAUDE.md) 的负载守卫和 Node 限流要求；当前协作要求下不要执行 `go test`
