package service

import (
	"context"
	"gbaseadmin/app/member/internal/model"
	"gbaseadmin/utility/snowflake"
	"github.com/gogf/gf/v2/net/ghttp"
)

type ILevel interface {
	Create(ctx context.Context, in *model.LevelCreateInput) error
	Update(ctx context.Context, in *model.LevelUpdateInput) error
	Delete(ctx context.Context, id snowflake.JsonInt64) error
	BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error
	Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.LevelDetailOutput, err error)
	List(ctx context.Context, in *model.LevelListInput) (list []*model.LevelListOutput, total int, err error)
	Export(ctx context.Context, in *model.LevelListInput) (list []*model.LevelListOutput, err error)
	BatchUpdate(ctx context.Context, in *model.LevelBatchUpdateInput) error
	Import(ctx context.Context, file *ghttp.UploadFile) (success int, fail int, err error)
}

var localLevel ILevel

func Level() ILevel {
	return localLevel
}

func RegisterLevel(i ILevel) {
	localLevel = i
}
