package team_export

import (
	"context"
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/frame/g"

	v1 "gbaseadmin/app/member/api/member/v1"
	"gbaseadmin/app/member/internal/model"
	"gbaseadmin/app/member/internal/service"
)

func csvSafeTeamExport(s string) string {
	if len(s) == 0 {
		return s
	}
	if s[0] == 0 || strings.ContainsAny(s[:1], "=+-@\t\r") {
		return "'" + s
	}
	return s
}

var TeamExport = cTeamExport{}

type cTeamExport struct{}

// Create 创建团队数据导出
func (c *cTeamExport) Create(ctx context.Context, req *v1.TeamExportCreateReq) (res *v1.TeamExportCreateRes, err error) {
	err = service.TeamExport().Create(ctx, &model.TeamExportCreateInput{
		UserID: req.UserID,
		TeamMemberCount: req.TeamMemberCount,
		ExportType: req.ExportType,
		FileURL: req.FileURL,
		FileSize: req.FileSize,
		DeployStatus: req.DeployStatus,
		DeployDomain: req.DeployDomain,
		DeployedAt: req.DeployedAt,
		Remark: req.Remark,
		Status: req.Status,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
	})
	return
}

// Update 更新团队数据导出
func (c *cTeamExport) Update(ctx context.Context, req *v1.TeamExportUpdateReq) (res *v1.TeamExportUpdateRes, err error) {
	err = service.TeamExport().Update(ctx, &model.TeamExportUpdateInput{
		ID: req.ID,
		UserID: req.UserID,
		TeamMemberCount: req.TeamMemberCount,
		ExportType: req.ExportType,
		FileURL: req.FileURL,
		FileSize: req.FileSize,
		DeployStatus: req.DeployStatus,
		DeployDomain: req.DeployDomain,
		DeployedAt: req.DeployedAt,
		Remark: req.Remark,
		Status: req.Status,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
	})
	return
}

// Delete 删除团队数据导出
func (c *cTeamExport) Delete(ctx context.Context, req *v1.TeamExportDeleteReq) (res *v1.TeamExportDeleteRes, err error) {
	err = service.TeamExport().Delete(ctx, req.ID)
	return
}

// BatchDelete 批量删除团队数据导出
func (c *cTeamExport) BatchDelete(ctx context.Context, req *v1.TeamExportBatchDeleteReq) (res *v1.TeamExportBatchDeleteRes, err error) {
	err = service.TeamExport().BatchDelete(ctx, req.IDs)
	return
}

// BatchUpdate 批量编辑团队数据导出
func (c *cTeamExport) BatchUpdate(ctx context.Context, req *v1.TeamExportBatchUpdateReq) (res *v1.TeamExportBatchUpdateRes, err error) {
	err = service.TeamExport().BatchUpdate(ctx, &model.TeamExportBatchUpdateInput{
		IDs: req.IDs,
		ExportType: req.ExportType,
		DeployStatus: req.DeployStatus,
		Status: req.Status,
	})
	return
}

// Detail 获取团队数据导出详情
func (c *cTeamExport) Detail(ctx context.Context, req *v1.TeamExportDetailReq) (res *v1.TeamExportDetailRes, err error) {
	res = &v1.TeamExportDetailRes{}
	res.TeamExportDetailOutput, err = service.TeamExport().Detail(ctx, req.ID)
	return
}

// List 获取团队数据导出列表
func (c *cTeamExport) List(ctx context.Context, req *v1.TeamExportListReq) (res *v1.TeamExportListRes, err error) {
	res = &v1.TeamExportListRes{}
	res.List, res.Total, err = service.TeamExport().List(ctx, &model.TeamExportListInput{
		PageNum:   req.PageNum,
		PageSize:  req.PageSize,
		OrderBy:   req.OrderBy,
		OrderDir:  req.OrderDir,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		UserID: req.UserID,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
		ExportType: req.ExportType,
		DeployStatus: req.DeployStatus,
		Status: req.Status,
		DeployDomain: req.DeployDomain,
		DeployedAtStart: req.DeployedAtStart,
		DeployedAtEnd: req.DeployedAtEnd,
	})
	return
}
// Export 导出团队数据导出
func (c *cTeamExport) Export(ctx context.Context, req *v1.TeamExportExportReq) (res *v1.TeamExportExportRes, err error) {
	list, err := service.TeamExport().Export(ctx, &model.TeamExportListInput{
		OrderBy:   req.OrderBy,
		OrderDir:  req.OrderDir,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		UserID: req.UserID,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
		ExportType: req.ExportType,
		DeployStatus: req.DeployStatus,
		Status: req.Status,
		DeployDomain: req.DeployDomain,
		DeployedAtStart: req.DeployedAtStart,
		DeployedAtEnd: req.DeployedAtEnd,
	})
	if err != nil {
		return
	}
	// CSV 导出（使用 csv.Writer 防止注入和格式问题）
	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Response.Header().Set("Content-Disposition", `attachment; filename="team_export.csv"`)
	r.Response.Write("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(r.Response.Writer)
	// 表头（与导入模板列对齐，末尾追加只读列）
	_ = w.Write([]string{"目标会员", "团队成员数", "导出类型", "导出文件地址", "文件大小", "部署状态", "部署域名", "备注", "状态", "创建时间"})
	// 数据行
	for _, item := range list {
		_ = w.Write([]string{
			csvSafeTeamExport(item.UserNickname),
			fmt.Sprintf("%v", item.TeamMemberCount),
			fmt.Sprintf("%v", item.ExportType),
			csvSafeTeamExport(item.FileURL),
			fmt.Sprintf("%v", item.FileSize),
			fmt.Sprintf("%v", item.DeployStatus),
			csvSafeTeamExport(item.DeployDomain),
			csvSafeTeamExport(item.Remark),
			fmt.Sprintf("%v", item.Status),
			func() string { if item.CreatedAt != nil { return item.CreatedAt.String() }; return "" }(),
		})
	}
	w.Flush()
	return
}

// Import 导入团队数据导出
func (c *cTeamExport) Import(ctx context.Context, req *v1.TeamExportImportReq) (res *v1.TeamExportImportRes, err error) {
	r := g.RequestFromCtx(ctx)
	file := r.GetUploadFile("file")
	if file == nil {
		return nil, fmt.Errorf("请上传文件")
	}
	success, fail, err := service.TeamExport().Import(ctx, file)
	if err != nil {
		return nil, err
	}
	res = &v1.TeamExportImportRes{Success: success, Fail: fail}
	return
}

// ImportTemplate 下载团队数据导出导入模板
func (c *cTeamExport) ImportTemplate(ctx context.Context, req *v1.TeamExportImportTemplateReq) (res *v1.TeamExportImportTemplateRes, err error) {
	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Response.Header().Set("Content-Disposition", `attachment; filename="team_export_template.csv"`)
	r.Response.Write("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(r.Response.Writer)
	_ = w.Write([]string{"目标会员", "团队成员数", "导出类型", "导出文件地址", "文件大小", "部署状态", "部署域名", "备注", "状态"})
	w.Flush()
	return
}
