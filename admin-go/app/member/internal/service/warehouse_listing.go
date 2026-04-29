package service

import (
	"context"
	"gbaseadmin/app/member/internal/model"
	"gbaseadmin/utility/snowflake"
	"github.com/gogf/gf/v2/net/ghttp"
)

type IWarehouseListing interface {
	Create(ctx context.Context, in *model.WarehouseListingCreateInput) error
	Update(ctx context.Context, in *model.WarehouseListingUpdateInput) error
	Delete(ctx context.Context, id snowflake.JsonInt64) error
	BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error
	Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.WarehouseListingDetailOutput, err error)
	List(ctx context.Context, in *model.WarehouseListingListInput) (list []*model.WarehouseListingListOutput, total int, err error)
	Export(ctx context.Context, in *model.WarehouseListingListInput) (list []*model.WarehouseListingListOutput, err error)
	BatchUpdate(ctx context.Context, in *model.WarehouseListingBatchUpdateInput) error
	Import(ctx context.Context, file *ghttp.UploadFile) (success int, fail int, err error)
}

var localWarehouseListing IWarehouseListing

func WarehouseListing() IWarehouseListing {
	return localWarehouseListing
}

func RegisterWarehouseListing(i IWarehouseListing) {
	localWarehouseListing = i
}
