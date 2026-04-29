package service

import (
	"context"
	"gbaseadmin/app/member/internal/model"
	"gbaseadmin/utility/snowflake"
	"github.com/gogf/gf/v2/net/ghttp"
)

type IShopOrder interface {
	Create(ctx context.Context, in *model.ShopOrderCreateInput) error
	Update(ctx context.Context, in *model.ShopOrderUpdateInput) error
	Delete(ctx context.Context, id snowflake.JsonInt64) error
	BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error
	Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.ShopOrderDetailOutput, err error)
	List(ctx context.Context, in *model.ShopOrderListInput) (list []*model.ShopOrderListOutput, total int, err error)
	Export(ctx context.Context, in *model.ShopOrderListInput) (list []*model.ShopOrderListOutput, err error)
	BatchUpdate(ctx context.Context, in *model.ShopOrderBatchUpdateInput) error
	Import(ctx context.Context, file *ghttp.UploadFile) (success int, fail int, err error)
}

var localShopOrder IShopOrder

func ShopOrder() IShopOrder {
	return localShopOrder
}

func RegisterShopOrder(i IShopOrder) {
	localShopOrder = i
}
