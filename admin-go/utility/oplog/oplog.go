package oplog

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

var insertOperationLog = func(ctx context.Context, data g.Map) {
	_, _ = g.DB().Ctx(ctx).Insert(ctx, "system_operation_log", data)
}

// Record 记录操作日志
// module: 模块名（如 order）
// action: 操作（create/update/delete/batch-delete/import）
// targetID: 操作目标 ID
// detail: 操作详情（可选）
func Record(ctx context.Context, module, action, targetID, detail string) {
	data := g.Map{
		"module":     module,
		"action":     action,
		"target_id":  targetID,
		"detail":     detail,
		"created_at": gtime.Now(),
	}
	go insertOperationLog(context.Background(), data)
}
