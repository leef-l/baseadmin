package service

import (
	"context"
	"gbaseadmin/app/demo/internal/model"
	"gbaseadmin/utility/snowflake"
	"github.com/gogf/gf/v2/net/ghttp"
)

type IWorkOrder interface {
	Create(ctx context.Context, in *model.WorkOrderCreateInput) error
	Update(ctx context.Context, in *model.WorkOrderUpdateInput) error
	Delete(ctx context.Context, id snowflake.JsonInt64) error
	BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error
	Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.WorkOrderDetailOutput, err error)
	List(ctx context.Context, in *model.WorkOrderListInput) (list []*model.WorkOrderListOutput, total int, err error)
	Export(ctx context.Context, in *model.WorkOrderListInput) (list []*model.WorkOrderListOutput, err error)
	BatchUpdate(ctx context.Context, in *model.WorkOrderBatchUpdateInput) error
	Import(ctx context.Context, file *ghttp.UploadFile) (success int, fail int, err error)
}

var localWorkOrder IWorkOrder

func WorkOrder() IWorkOrder {
	return localWorkOrder
}

func RegisterWorkOrder(i IWorkOrder) {
	localWorkOrder = i
}
