# Git 提交流程铁律

## 铁律

1. 整体功能完成后，必须立即提交并推送到 GitHub，禁止停留在“本地已完成但未推送”
2. 默认使用仓库根目录脚本 `./scripts/feature-publish.sh "type(scope): summary"` 完成提交
3. 提交前必须确认当前改动属于同一整体功能；不允许把无关改动混进同一个功能提交
4. 如果功能包含数据库变更，必须先补 `admin-go/database/migrations/` 迁移文件，再执行提交脚本
5. 推送目标默认是 `origin` 当前分支；除非有明确要求，不要临时改推送目标

## 标准命令

```bash
./scripts/feature-publish.sh "feat(system): complete feature delivery"
```

脚本会自动执行：

1. `git add -A`
2. `git commit -m "..."`
3. `git push origin 当前分支`

## 适用说明

- 这条规则用于“一个整体功能已经完成”的节点，不要求每个小修改都单独推送
- 如果当前工作区存在不属于本次功能的改动，先拆分，再提交
- 如果本次功能涉及数据库，迁移文件和业务代码必须同批提交
