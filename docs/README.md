# 文档目录

`docs/` 只保留当前仍长期维护的项目说明。一次性协作提示、临时排障记录不再放在这里；长跑巡检的滚动日志集中在 `流程日志/` 子目录。

## 首选阅读

- ★ [代码库导读](代码库导读.md) — 新人从这里开始：真源表 + 架构地图 + 日常维护红线

## 项目总览

- [System 应用说明](../admin-go/app/system/README.MD)
- [Upload 应用说明](../admin-go/app/upload/README.MD)
- [代码生成器说明](../admin-go/codegen/README.md)
- [Codegen AI 执行手册](Codegen-AI执行手册.md)

## 开发与部署

- [Docker 开发说明](Docker开发说明.md)
- [生产运维说明](生产运维说明.md)

## 协作约定

- [AI 协作铁律](../CLAUDE.md)
- [Git 提交流程铁律](Git提交流程铁律.md)

## AI 执行入口

- Codegen 任务先读 [Codegen AI 执行手册](Codegen-AI执行手册.md)，再读 [代码生成器说明](../admin-go/codegen/README.md)
- 执行 codegen、Go 测试、前端校验时必须遵守 [AI 协作铁律](../CLAUDE.md) 的负载守卫和 Node 限流要求

## 流程日志（长跑巡检专用，非日常阅读）

- [持续优化执行说明](流程日志/持续优化执行说明.md)
- [优化审计记录](流程日志/优化审计记录.md)
