package team_export

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	v1 "gbaseadmin/app/member/api/member/v1"
	"gbaseadmin/app/member/internal/logic/team_export"
	"gbaseadmin/app/member/internal/middleware"
)

// Run 后台触发导出（异步：立即返回 export_id，状态后台更新）。
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

// Status 查询导出任务状态（前端轮询用）。
func (c *cTeamExport) Status(ctx context.Context, req *v1.TeamExportStatusReq) (res *v1.TeamExportStatusRes, err error) {
	status, fileURL, fileSize, members, errReason, err := team_export.GetExport(ctx, int64(req.ExportID))
	if err != nil {
		return nil, err
	}
	return &v1.TeamExportStatusRes{
		ExportID:    fmt.Sprintf("%d", int64(req.ExportID)),
		Status:      status,
		StatusText:  exportStatusText(status),
		FileURL:     fileURL,
		FileSize:    fileSize,
		MemberCount: members,
		ErrReason:   errReason,
	}, nil
}

// Download 下载导出文件（流式）。
// 前端用 window.open('/api/member/team_export/download?exportId=...')。
func Download(r *ghttp.Request) {
	ctx := r.Context()
	exportID, _ := strconv.ParseInt(r.GetQuery("exportId").String(), 10, 64)
	if exportID <= 0 {
		r.Response.WriteStatus(400, "exportId 不能为空")
		return
	}
	status, fileURL, _, _, _, err := team_export.GetExport(ctx, exportID)
	if err != nil {
		g.Log().Warningf(ctx, "team_export Download err=%v", err)
		r.Response.WriteStatus(404, "导出不存在")
		return
	}
	if status != 2 || fileURL == "" {
		r.Response.WriteStatus(409, "导出尚未就绪")
		return
	}
	root := strings.TrimSpace(g.Cfg().MustGet(ctx, "member.teamExportRoot").String())
	if root == "" {
		root = filepath.Join("resource", "team-export")
	}
	if abs, absErr := filepath.Abs(root); absErr == nil {
		root = abs
	}
	rel := strings.TrimPrefix(fileURL, "/team-export/")
	full := filepath.Join(root, rel)
	f, err := os.Open(full)
	if err != nil {
		g.Log().Warningf(ctx, "open export file err=%v path=%s", err, full)
		r.Response.WriteStatus(404, "文件不存在")
		return
	}
	defer f.Close()
	r.Response.Header().Set("Content-Type", "application/gzip")
	r.Response.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, rel))
	if _, err := io.Copy(r.Response.RawWriter(), f); err != nil {
		g.Log().Warningf(ctx, "stream export file err=%v", err)
	}
}

func exportStatusText(s int) string {
	switch s {
	case 0:
		return "排队中"
	case 1:
		return "导出中"
	case 2:
		return "已就绪"
	case 3:
		return "失败"
	}
	return ""
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
