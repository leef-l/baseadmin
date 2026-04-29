package service

import (
	"context"
	"gbaseadmin/app/member/internal/model"
	"gbaseadmin/utility/snowflake"
	"github.com/gogf/gf/v2/net/ghttp"
)

type IWarehouseGoods interface {
	Create(ctx context.Context, in *model.WarehouseGoodsCreateInput) error
	Update(ctx context.Context, in *model.WarehouseGoodsUpdateInput) error
	Delete(ctx context.Context, id snowflake.JsonInt64) error
	BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error
	Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.WarehouseGoodsDetailOutput, err error)
	List(ctx context.Context, in *model.WarehouseGoodsListInput) (list []*model.WarehouseGoodsListOutput, total int, err error)
	Export(ctx context.Context, in *model.WarehouseGoodsListInput) (list []*model.WarehouseGoodsListOutput, err error)
	BatchUpdate(ctx context.Context, in *model.WarehouseGoodsBatchUpdateInput) error
	Import(ctx context.Context, file *ghttp.UploadFile) (success int, fail int, err error)
}

var localWarehouseGoods IWarehouseGoods

func WarehouseGoods() IWarehouseGoods {
	return localWarehouseGoods
}

func RegisterWarehouseGoods(i IWarehouseGoods) {
	localWarehouseGoods = i
}
