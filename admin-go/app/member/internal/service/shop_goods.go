package service

import (
	"context"
	"gbaseadmin/app/member/internal/model"
	"gbaseadmin/utility/snowflake"
	"github.com/gogf/gf/v2/net/ghttp"
)

type IShopGoods interface {
	Create(ctx context.Context, in *model.ShopGoodsCreateInput) error
	Update(ctx context.Context, in *model.ShopGoodsUpdateInput) error
	Delete(ctx context.Context, id snowflake.JsonInt64) error
	BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error
	Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.ShopGoodsDetailOutput, err error)
	List(ctx context.Context, in *model.ShopGoodsListInput) (list []*model.ShopGoodsListOutput, total int, err error)
	Export(ctx context.Context, in *model.ShopGoodsListInput) (list []*model.ShopGoodsListOutput, err error)
	BatchUpdate(ctx context.Context, in *model.ShopGoodsBatchUpdateInput) error
	Import(ctx context.Context, file *ghttp.UploadFile) (success int, fail int, err error)
}

var localShopGoods IShopGoods

func ShopGoods() IShopGoods {
	return localShopGoods
}

func RegisterShopGoods(i IShopGoods) {
	localShopGoods = i
}
