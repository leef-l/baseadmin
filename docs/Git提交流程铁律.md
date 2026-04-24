# Git 提交流程铁律

## 铁律

1. 整体功能完成后，必须立即提交并推送到 GitHub，禁止停留在“本地已完成但未推送”
2. 默认使用仓库根目录脚本 `./scripts/feature-publish.sh "type(scope): summary"` 完成提交
3. 仓库验证统一只认 GitHub Actions；要验证，就提交并推送，然后查看 `CI GBaseAdmin`
4. 如果推送的是 `main`，除了 `CI GBaseAdmin`，还必须关注 `Deploy GBaseAdmin`
5. 提交前必须确认当前改动属于同一整体功能；不允许把无关改动混进同一个功能提交
6. 如果功能包含数据库变更，必须先补 `admin-go/database/migrations/` 迁移文件，再执行提交脚本
7. 推送目标默认是 `origin` 当前分支；除非有明确要求，不要临时改推送目标

## 标准命令

```bash
./scripts/feature-publish.sh "feat(system): complete feature delivery"
```

默认脚本会自动执行：

1. 校验工作区里没有未暂存 / 未跟踪改动
2. 基于已暂存内容执行提交
3. `git push origin 当前分支`
4. 触发 GitHub Actions 校验

如果确认当前工作区所有改动都属于同一批交付，也可以显式使用：

```bash
./scripts/feature-publish.sh --all "feat(system): complete feature delivery"
```

这时脚本会自动执行：

1. `git add -A`
2. `git commit -m "..."`
3. `git push origin 当前分支`
4. 触发 GitHub Actions 校验

## 使用边界

- 默认模式下，必须先用精确 `git add <path>` 准备好本次提交；脚本不会替你吞掉其他脏改动
- 如果工作区里混有别人的改动、临时实验或其他未完成功能，不要使用 `--all`
- 只有在确认所有改动都属于同一批交付时，才允许使用 `--all`
- 如果 GitHub Actions 失败，不要把“本地看起来没问题”当作完成状态；修完后重新提交并推送

## 适用说明

- 这条规则用于“一个整体功能已经完成”的节点，不要求每个小修改都单独推送
- 如果当前工作区存在不属于本次功能的改动，先拆分，再提交
- 如果本次功能涉及数据库，迁移文件和业务代码必须同批提交
- 根目录 `lefthook.yml` 只保留提醒，不承担仓库验证

## GitHub Actions 现状

- `CI GBaseAdmin`：按路径执行 `verify-baseadmin-scope`、后端测试、codegen 测试、前端页面规则校验、前端 typecheck、`web-antd` 单测
- `CI GBaseAdmin`：同时输出 Actions Summary，并上传后端/前端的 JUnit、JSON、文本日志 artifact，测试记录统一看这里
- `Deploy GBaseAdmin`：`main` 分支在部署前再次执行上线门禁，并顺序部署 `system`、`upload`、`frontend`
- `Deploy GBaseAdmin`：使用专用 smoke 账号做真实登录、目录规则、上传、公网 URL 冒烟；需要仓库 secrets `SMOKE_USERNAME`、`SMOKE_PASSWORD`

说明：

- 日常交付不再要求本地手工跑完整校验；默认路径就是 `commit + push`
- 如需本地排查，仍可手工执行 `scripts/verify-baseadmin-scope.sh`、`scripts/verify-vben-pages.sh` 等轻量脚本；后端测试必须走 `scripts/run-go-task-with-limits.sh go test ./...`，前端校验必须走 `scripts/run-node-task-with-limits.sh pnpm -C vue-vben-admin -F @vben/web-antd typecheck`
