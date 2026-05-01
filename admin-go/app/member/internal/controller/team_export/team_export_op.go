package team_export

import (
	"context"

	v1 "gbaseadmin/app/member/api/member/v1"
	"gbaseadmin/app/member/internal/logic/team_export"
	"gbaseadmin/app/member/internal/middleware"
)

// Run 后台触发导出。
func (c *cTeamExport) Run(ctx context.Context, req *v1.TeamExportRunReq) (res *v1.TeamExportRunRes, err error) {
	operator := int64(middleware.GetUserID(ctx))
	out, err := team_export.Export(ctx, &team_export.ExportInput{
		UserID:     int64(req.UserID),
		ExportType: 1,
		OperatorID: operator,
		Remark:     req.Remark,
	})
	if err != nil {
		return nil, err
	}
	return &v1.TeamExportRunRes{
		ExportID:    out.ExportID,
		FileURL:     out.FileURL,
		FileSize:    out.FileSize,
		MemberCount: out.MemberCount,
	}, nil
}

// Deploy 站点裂变（占位）。
func (c *cTeamExport) Deploy(ctx context.Context, req *v1.TeamExportDeployReq) (res *v1.TeamExportDeployRes, err error) {
	operator := int64(middleware.GetUserID(ctx))
	out, err := team_export.Deploy(ctx, &team_export.DeployInput{
		ExportID:   int64(req.ExportID),
		Domain:     req.Domain,
		OperatorID: operator,
	})
	if err != nil {
		return nil, err
	}
	return &v1.TeamExportDeployRes{
		DeployStatus: out.DeployStatus,
		DeployDomain: out.DeployDomain,
	}, nil
}
