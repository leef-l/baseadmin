package service

import (
	"context"

	"gbaseadmin/app/system/internal/model"
	"gbaseadmin/utility/snowflake"
)

type ITenant interface {
	Create(ctx context.Context, in *model.TenantCreateInput) error
	Update(ctx context.Context, in *model.TenantUpdateInput) error
	Delete(ctx context.Context, id snowflake.JsonInt64) error
	BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error
	Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.TenantDetailOutput, err error)
	List(ctx context.Context, in *model.TenantListInput) (list []*model.TenantListOutput, total int, err error)
}

var localTenant ITenant

func Tenant() ITenant {
	return localTenant
}

func RegisterTenant(i ITenant) {
	localTenant = i
}
