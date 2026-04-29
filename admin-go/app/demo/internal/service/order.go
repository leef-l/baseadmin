package service

import (
	"context"
	"gbaseadmin/app/demo/internal/model"
	"gbaseadmin/utility/snowflake"
	"github.com/gogf/gf/v2/net/ghttp"
)

type IOrder interface {
	Create(ctx context.Context, in *model.OrderCreateInput) error
	Update(ctx context.Context, in *model.OrderUpdateInput) error
	Delete(ctx context.Context, id snowflake.JsonInt64) error
	BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error
	Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.OrderDetailOutput, err error)
	List(ctx context.Context, in *model.OrderListInput) (list []*model.OrderListOutput, total int, err error)
	Export(ctx context.Context, in *model.OrderListInput) (list []*model.OrderListOutput, err error)
	BatchUpdate(ctx context.Context, in *model.OrderBatchUpdateInput) error
	Import(ctx context.Context, file *ghttp.UploadFile) (success int, fail int, err error)
}

var localOrder IOrder

func Order() IOrder {
	return localOrder
}

func RegisterOrder(i IOrder) {
	localOrder = i
}
