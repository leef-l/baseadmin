package service

import (
	"context"
	"gbaseadmin/app/system/internal/model"
	"gbaseadmin/utility/snowflake"
	"github.com/gogf/gf/v2/net/ghttp"
)

type ITenantPlan interface {
	Create(ctx context.Context, in *model.TenantPlanCreateInput) error
	Update(ctx context.Context, in *model.TenantPlanUpdateInput) error
	Delete(ctx context.Context, id snowflake.JsonInt64) error
	BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error
	Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.TenantPlanDetailOutput, err error)
	List(ctx context.Context, in *model.TenantPlanListInput) (list []*model.TenantPlanListOutput, total int, err error)
	Export(ctx context.Context, in *model.TenantPlanListInput) (list []*model.TenantPlanListOutput, err error)
	BatchUpdate(ctx context.Context, in *model.TenantPlanBatchUpdateInput) error
	Import(ctx context.Context, file *ghttp.UploadFile) (success int, fail int, err error)
}

var localTenantPlan ITenantPlan

func TenantPlan() ITenantPlan {
	return localTenantPlan
}

func RegisterTenantPlan(i ITenantPlan) {
	localTenantPlan = i
}
