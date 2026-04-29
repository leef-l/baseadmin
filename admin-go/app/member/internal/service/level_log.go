package service

import (
	"context"
	"gbaseadmin/app/member/internal/model"
	"gbaseadmin/utility/snowflake"
	"github.com/gogf/gf/v2/net/ghttp"
)

type ILevelLog interface {
	Create(ctx context.Context, in *model.LevelLogCreateInput) error
	Update(ctx context.Context, in *model.LevelLogUpdateInput) error
	Delete(ctx context.Context, id snowflake.JsonInt64) error
	BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error
	Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.LevelLogDetailOutput, err error)
	List(ctx context.Context, in *model.LevelLogListInput) (list []*model.LevelLogListOutput, total int, err error)
	Export(ctx context.Context, in *model.LevelLogListInput) (list []*model.LevelLogListOutput, err error)
	Import(ctx context.Context, file *ghttp.UploadFile) (success int, fail int, err error)
}

var localLevelLog ILevelLog

func LevelLog() ILevelLog {
	return localLevelLog
}

func RegisterLevelLog(i ILevelLog) {
	localLevelLog = i
}
