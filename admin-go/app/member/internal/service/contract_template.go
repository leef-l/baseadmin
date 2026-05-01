package service

import (
	"context"
	"gbaseadmin/app/member/internal/model"
	"gbaseadmin/utility/snowflake"
	"github.com/gogf/gf/v2/net/ghttp"
)

type IContractTemplate interface {
	Create(ctx context.Context, in *model.ContractTemplateCreateInput) error
	Update(ctx context.Context, in *model.ContractTemplateUpdateInput) error
	Delete(ctx context.Context, id snowflake.JsonInt64) error
	BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error
	Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.ContractTemplateDetailOutput, err error)
	List(ctx context.Context, in *model.ContractTemplateListInput) (list []*model.ContractTemplateListOutput, total int, err error)
	Export(ctx context.Context, in *model.ContractTemplateListInput) (list []*model.ContractTemplateListOutput, err error)
	BatchUpdate(ctx context.Context, in *model.ContractTemplateBatchUpdateInput) error
	Import(ctx context.Context, file *ghttp.UploadFile) (success int, fail int, err error)
}

var localContractTemplate IContractTemplate

func ContractTemplate() IContractTemplate {
	return localContractTemplate
}

func RegisterContractTemplate(i IContractTemplate) {
	localContractTemplate = i
}
