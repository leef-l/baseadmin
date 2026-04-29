package service

import (
	"context"
	"gbaseadmin/app/member/internal/model"
	"gbaseadmin/utility/snowflake"
)

type IShopCategory interface {
	Create(ctx context.Context, in *model.ShopCategoryCreateInput) error
	Update(ctx context.Context, in *model.ShopCategoryUpdateInput) error
	Delete(ctx context.Context, id snowflake.JsonInt64) error
	BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error
	Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.ShopCategoryDetailOutput, err error)
	List(ctx context.Context, in *model.ShopCategoryListInput) (list []*model.ShopCategoryListOutput, total int, err error)
	Export(ctx context.Context, in *model.ShopCategoryListInput) (list []*model.ShopCategoryListOutput, err error)
	Tree(ctx context.Context, in *model.ShopCategoryTreeInput) (tree []*model.ShopCategoryTreeOutput, err error)
	BatchUpdate(ctx context.Context, in *model.ShopCategoryBatchUpdateInput) error
}

var localShopCategory IShopCategory

func ShopCategory() IShopCategory {
	return localShopCategory
}

func RegisterShopCategory(i IShopCategory) {
	localShopCategory = i
}
