package service

import (
	"context"
	"gbaseadmin/app/member/internal/model"
	"gbaseadmin/utility/snowflake"
	"github.com/gogf/gf/v2/net/ghttp"
)

type IWarehouseTrade interface {
	Create(ctx context.Context, in *model.WarehouseTradeCreateInput) error
	Update(ctx context.Context, in *model.WarehouseTradeUpdateInput) error
	Delete(ctx context.Context, id snowflake.JsonInt64) error
	BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error
	Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.WarehouseTradeDetailOutput, err error)
	List(ctx context.Context, in *model.WarehouseTradeListInput) (list []*model.WarehouseTradeListOutput, total int, err error)
	Export(ctx context.Context, in *model.WarehouseTradeListInput) (list []*model.WarehouseTradeListOutput, err error)
	BatchUpdate(ctx context.Context, in *model.WarehouseTradeBatchUpdateInput) error
	Import(ctx context.Context, file *ghttp.UploadFile) (success int, fail int, err error)
}

var localWarehouseTrade IWarehouseTrade

func WarehouseTrade() IWarehouseTrade {
	return localWarehouseTrade
}

func RegisterWarehouseTrade(i IWarehouseTrade) {
	localWarehouseTrade = i
}
