package service

import (
	"context"
	"gbaseadmin/app/system/internal/model"
	"gbaseadmin/utility/snowflake"
	"github.com/gogf/gf/v2/net/ghttp"
)

type ICron interface {
	Create(ctx context.Context, in *model.CronCreateInput) error
	Update(ctx context.Context, in *model.CronUpdateInput) error
	Delete(ctx context.Context, id snowflake.JsonInt64) error
	BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error
	Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.CronDetailOutput, err error)
	List(ctx context.Context, in *model.CronListInput) (list []*model.CronListOutput, total int, err error)
	Export(ctx context.Context, in *model.CronListInput) (list []*model.CronListOutput, err error)
	BatchUpdate(ctx context.Context, in *model.CronBatchUpdateInput) error
	Import(ctx context.Context, file *ghttp.UploadFile) (success int, fail int, err error)
	Trigger(ctx context.Context, name string, handler string) (success bool, message string)
	LogList(ctx context.Context, cronID int64, pageNum, pageSize int) (list []*model.CronLogItem, total int, err error)
}

var localCron ICron

func Cron() ICron {
	return localCron
}

func RegisterCron(i ICron) {
	localCron = i
}
