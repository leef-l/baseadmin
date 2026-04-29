package service

import (
	"context"
	"gbaseadmin/app/member/internal/model"
	"gbaseadmin/utility/snowflake"
	"github.com/gogf/gf/v2/net/ghttp"
)

type IWalletLog interface {
	Create(ctx context.Context, in *model.WalletLogCreateInput) error
	Update(ctx context.Context, in *model.WalletLogUpdateInput) error
	Delete(ctx context.Context, id snowflake.JsonInt64) error
	BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error
	Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.WalletLogDetailOutput, err error)
	List(ctx context.Context, in *model.WalletLogListInput) (list []*model.WalletLogListOutput, total int, err error)
	Export(ctx context.Context, in *model.WalletLogListInput) (list []*model.WalletLogListOutput, err error)
	Import(ctx context.Context, file *ghttp.UploadFile) (success int, fail int, err error)
}

var localWalletLog IWalletLog

func WalletLog() IWalletLog {
	return localWalletLog
}

func RegisterWalletLog(i IWalletLog) {
	localWalletLog = i
}
