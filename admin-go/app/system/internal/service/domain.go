package service

import (
	"context"

	"gbaseadmin/app/system/internal/model"
	"gbaseadmin/utility/snowflake"
)

type IDomain interface {
	Create(ctx context.Context, in *model.DomainCreateInput) error
	Update(ctx context.Context, in *model.DomainUpdateInput) error
	Delete(ctx context.Context, id snowflake.JsonInt64) error
	BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error
	Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.DomainDetailOutput, err error)
	List(ctx context.Context, in *model.DomainListInput) (list []*model.DomainListOutput, total int, err error)
	ApplyNginx(ctx context.Context, id snowflake.JsonInt64) (out *model.DomainApplyNginxOutput, err error)
	ApplySSL(ctx context.Context, id snowflake.JsonInt64) (out *model.DomainApplySSLOutput, err error)
}

var localDomain IDomain

func Domain() IDomain {
	return localDomain
}

func RegisterDomain(i IDomain) {
	localDomain = i
}
