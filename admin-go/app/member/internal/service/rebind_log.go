package service

import (
	"context"
	"gbaseadmin/app/member/internal/model"
	"gbaseadmin/utility/snowflake"
	"github.com/gogf/gf/v2/net/ghttp"
)

type IRebindLog interface {
	Create(ctx context.Context, in *model.RebindLogCreateInput) error
	Update(ctx context.Context, in *model.RebindLogUpdateInput) error
	Delete(ctx context.Context, id snowflake.JsonInt64) error
	BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error
	Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.RebindLogDetailOutput, err error)
	List(ctx context.Context, in *model.RebindLogListInput) (list []*model.RebindLogListOutput, total int, err error)
	Export(ctx context.Context, in *model.RebindLogListInput) (list []*model.RebindLogListOutput, err error)
	Import(ctx context.Context, file *ghttp.UploadFile) (success int, fail int, err error)
}

var localRebindLog IRebindLog

func RebindLog() IRebindLog {
	return localRebindLog
}

func RegisterRebindLog(i IRebindLog) {
	localRebindLog = i
}
