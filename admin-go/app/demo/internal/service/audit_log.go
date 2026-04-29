package service

import (
	"context"
	"gbaseadmin/app/demo/internal/model"
	"gbaseadmin/utility/snowflake"
	"github.com/gogf/gf/v2/net/ghttp"
)

type IAuditLog interface {
	Create(ctx context.Context, in *model.AuditLogCreateInput) error
	Update(ctx context.Context, in *model.AuditLogUpdateInput) error
	Delete(ctx context.Context, id snowflake.JsonInt64) error
	BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error
	Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.AuditLogDetailOutput, err error)
	List(ctx context.Context, in *model.AuditLogListInput) (list []*model.AuditLogListOutput, total int, err error)
	Export(ctx context.Context, in *model.AuditLogListInput) (list []*model.AuditLogListOutput, err error)
	Import(ctx context.Context, file *ghttp.UploadFile) (success int, fail int, err error)
}

var localAuditLog IAuditLog

func AuditLog() IAuditLog {
	return localAuditLog
}

func RegisterAuditLog(i IAuditLog) {
	localAuditLog = i
}
