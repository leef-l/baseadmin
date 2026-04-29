package service

import (
	"context"
	"gbaseadmin/app/demo/internal/model"
	"gbaseadmin/utility/snowflake"
	"github.com/gogf/gf/v2/net/ghttp"
)

type IContract interface {
	Create(ctx context.Context, in *model.ContractCreateInput) error
	Update(ctx context.Context, in *model.ContractUpdateInput) error
	Delete(ctx context.Context, id snowflake.JsonInt64) error
	BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error
	Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.ContractDetailOutput, err error)
	List(ctx context.Context, in *model.ContractListInput) (list []*model.ContractListOutput, total int, err error)
	Export(ctx context.Context, in *model.ContractListInput) (list []*model.ContractListOutput, err error)
	BatchUpdate(ctx context.Context, in *model.ContractBatchUpdateInput) error
	Import(ctx context.Context, file *ghttp.UploadFile) (success int, fail int, err error)
}

var localContract IContract

func Contract() IContract {
	return localContract
}

func RegisterContract(i IContract) {
	localContract = i
}
