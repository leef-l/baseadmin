# Database Migrations

铁律：

- 以后所有数据库结构变更、初始化数据变更、默认菜单/默认账号/默认上传配置变更，必须新增 `golang-migrate` 迁移文件。
- 不再把 `docker/mysql/init.sql`、`admin-go/codegen/sql/init.sql` 当作数据库真源。
- Docker 容器启动前必须先执行 `migrate up`，再启动应用。
- 默认由 `system` 容器负责自动执行迁移；`upload` 容器默认不重复抢跑迁移。

常用命令：

```bash
cd admin-go
go run ./cmd/migrate up
go run ./cmd/migrate version
go run ./cmd/migrate create add_system_logs
```

默认迁移目录：`database/migrations`
