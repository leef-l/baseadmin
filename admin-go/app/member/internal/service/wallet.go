package service

import (
	"context"
	"gbaseadmin/app/member/internal/model"
	"gbaseadmin/utility/snowflake"
	"github.com/gogf/gf/v2/net/ghttp"
)

type IWallet interface {
	Create(ctx context.Context, in *model.WalletCreateInput) error
	Update(ctx context.Context, in *model.WalletUpdateInput) error
	Delete(ctx context.Context, id snowflake.JsonInt64) error
	BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error
	Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.WalletDetailOutput, err error)
	List(ctx context.Context, in *model.WalletListInput) (list []*model.WalletListOutput, total int, err error)
	Export(ctx context.Context, in *model.WalletListInput) (list []*model.WalletListOutput, err error)
	BatchUpdate(ctx context.Context, in *model.WalletBatchUpdateInput) error
	Import(ctx context.Context, file *ghttp.UploadFile) (success int, fail int, err error)
}

var localWallet IWallet

func Wallet() IWallet {
	return localWallet
}

func RegisterWallet(i IWallet) {
	localWallet = i
}
