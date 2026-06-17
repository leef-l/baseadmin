package oplog

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/utility/snowflake"
)

var insertOperationLog = func(ctx context.Context, data g.Map) {
	_, _ = g.DB().Ctx(ctx).Insert(ctx, "system_operation_log", data)
}

// Record 记录操作日志
// module: 模块名（如 order）
// action: 操作（create/update/delete/batch-delete/import）
// targetID: 操作目标 ID
// detail: 操作详情（可选）
// 自动从请求上下文提取 tenant_id / user_id / username
func Record(ctx context.Context, module, action, targetID, detail string) {
	module = strings.TrimSpace(module)
	action = strings.TrimSpace(action)
	if module == "" || action == "" {
		return
	}

	var (
		tenantID int64
		userID   int64
		username string
	)
	if req := ghttp.RequestFromCtx(ctx); req != nil {
		tenantID = req.GetCtxVar("jwt_tenant_id").Int64()
		userID = req.GetCtxVar("jwt_user_id").Int64()
		username = req.GetCtxVar("jwt_username").String()
	}

	data := g.Map{
		"id":         snowflake.Generate(),
		"module":     module,
		"action":     action,
		"target_id":  strings.TrimSpace(targetID),
		"detail":     strings.TrimSpace(detail),
		"tenant_id":  tenantID,
		"user_id":    userID,
		"username":   username,
		"created_at": gtime.Now(),
	}
	if ctx == nil {
		ctx = context.Background()
	}
	go insertOperationLog(context.WithoutCancel(ctx), data)
}
