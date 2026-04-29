package service

import (
	"context"
	"gbaseadmin/app/demo/internal/model"
	"gbaseadmin/utility/snowflake"
	"github.com/gogf/gf/v2/net/ghttp"
)

type IProduct interface {
	Create(ctx context.Context, in *model.ProductCreateInput) error
	Update(ctx context.Context, in *model.ProductUpdateInput) error
	Delete(ctx context.Context, id snowflake.JsonInt64) error
	BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error
	Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.ProductDetailOutput, err error)
	List(ctx context.Context, in *model.ProductListInput) (list []*model.ProductListOutput, total int, err error)
	Export(ctx context.Context, in *model.ProductListInput) (list []*model.ProductListOutput, err error)
	BatchUpdate(ctx context.Context, in *model.ProductBatchUpdateInput) error
	Import(ctx context.Context, file *ghttp.UploadFile) (success int, fail int, err error)
}

var localProduct IProduct

func Product() IProduct {
	return localProduct
}

func RegisterProduct(i IProduct) {
	localProduct = i
}
