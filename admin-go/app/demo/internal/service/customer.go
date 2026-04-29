package service

import (
	"context"
	"gbaseadmin/app/demo/internal/model"
	"gbaseadmin/utility/snowflake"
	"github.com/gogf/gf/v2/net/ghttp"
)

type ICustomer interface {
	Create(ctx context.Context, in *model.CustomerCreateInput) error
	Update(ctx context.Context, in *model.CustomerUpdateInput) error
	Delete(ctx context.Context, id snowflake.JsonInt64) error
	BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error
	Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.CustomerDetailOutput, err error)
	List(ctx context.Context, in *model.CustomerListInput) (list []*model.CustomerListOutput, total int, err error)
	Export(ctx context.Context, in *model.CustomerListInput) (list []*model.CustomerListOutput, err error)
	BatchUpdate(ctx context.Context, in *model.CustomerBatchUpdateInput) error
	Import(ctx context.Context, file *ghttp.UploadFile) (success int, fail int, err error)
}

var localCustomer ICustomer

func Customer() ICustomer {
	return localCustomer
}

func RegisterCustomer(i ICustomer) {
	localCustomer = i
}
