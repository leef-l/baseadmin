package service

import (
	"context"
	"gbaseadmin/app/demo/internal/model"
	"gbaseadmin/utility/snowflake"
	"github.com/gogf/gf/v2/net/ghttp"
)

type ISurvey interface {
	Create(ctx context.Context, in *model.SurveyCreateInput) error
	Update(ctx context.Context, in *model.SurveyUpdateInput) error
	Delete(ctx context.Context, id snowflake.JsonInt64) error
	BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error
	Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.SurveyDetailOutput, err error)
	List(ctx context.Context, in *model.SurveyListInput) (list []*model.SurveyListOutput, total int, err error)
	Export(ctx context.Context, in *model.SurveyListInput) (list []*model.SurveyListOutput, err error)
	BatchUpdate(ctx context.Context, in *model.SurveyBatchUpdateInput) error
	Import(ctx context.Context, file *ghttp.UploadFile) (success int, fail int, err error)
}

var localSurvey ISurvey

func Survey() ISurvey {
	return localSurvey
}

func RegisterSurvey(i ISurvey) {
	localSurvey = i
}
