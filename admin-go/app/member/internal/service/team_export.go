package service

import (
	"context"
	"gbaseadmin/app/member/internal/model"
	"gbaseadmin/utility/snowflake"
	"github.com/gogf/gf/v2/net/ghttp"
)

type ITeamExport interface {
	Create(ctx context.Context, in *model.TeamExportCreateInput) error
	Update(ctx context.Context, in *model.TeamExportUpdateInput) error
	Delete(ctx context.Context, id snowflake.JsonInt64) error
	BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error
	Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.TeamExportDetailOutput, err error)
	List(ctx context.Context, in *model.TeamExportListInput) (list []*model.TeamExportListOutput, total int, err error)
	Export(ctx context.Context, in *model.TeamExportListInput) (list []*model.TeamExportListOutput, err error)
	BatchUpdate(ctx context.Context, in *model.TeamExportBatchUpdateInput) error
	Import(ctx context.Context, file *ghttp.UploadFile) (success int, fail int, err error)
}

var localTeamExport ITeamExport

func TeamExport() ITeamExport {
	return localTeamExport
}

func RegisterTeamExport(i ITeamExport) {
	localTeamExport = i
}
