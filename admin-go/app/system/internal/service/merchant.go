package service

import (
	"context"

	"gbaseadmin/app/system/internal/model"
	"gbaseadmin/utility/snowflake"
)

type IMerchant interface {
	Create(ctx context.Context, in *model.MerchantCreateInput) error
	Update(ctx context.Context, in *model.MerchantUpdateInput) error
	Delete(ctx context.Context, id snowflake.JsonInt64) error
	BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error
	Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.MerchantDetailOutput, err error)
	List(ctx context.Context, in *model.MerchantListInput) (list []*model.MerchantListOutput, total int, err error)
}

var localMerchant IMerchant

func Merchant() IMerchant {
	return localMerchant
}

func RegisterMerchant(i IMerchant) {
	localMerchant = i
}
