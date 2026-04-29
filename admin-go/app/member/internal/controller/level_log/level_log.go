package level_log

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

func csvSafeLevelLog(s string) string {
	if len(s) == 0 {
		return s
	}
	if s[0] == 0 || strings.ContainsAny(s[:1], "=+-@\t\r") {
		return "'" + s
	}
	return s
}

var LevelLog = cLevelLog{}

type cLevelLog struct{}

// Create 创建等级变更日志
func (c *cLevelLog) Create(ctx context.Context, req *v1.LevelLogCreateReq) (res *v1.LevelLogCreateRes, err error) {
	err = service.LevelLog().Create(ctx, &model.LevelLogCreateInput{
		UserID: req.UserID,
		OldLevelID: req.OldLevelID,
		NewLevelID: req.NewLevelID,
		ChangeType: req.ChangeType,
		ExpireAt: req.ExpireAt,
		Remark: req.Remark,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
	})
	return
}

// Update 更新等级变更日志
func (c *cLevelLog) Update(ctx context.Context, req *v1.LevelLogUpdateReq) (res *v1.LevelLogUpdateRes, err error) {
	err = service.LevelLog().Update(ctx, &model.LevelLogUpdateInput{
		ID: req.ID,
		UserID: req.UserID,
		OldLevelID: req.OldLevelID,
		NewLevelID: req.NewLevelID,
		ChangeType: req.ChangeType,
		ExpireAt: req.ExpireAt,
		Remark: req.Remark,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
	})
	return
}

// Delete 删除等级变更日志
func (c *cLevelLog) Delete(ctx context.Context, req *v1.LevelLogDeleteReq) (res *v1.LevelLogDeleteRes, err error) {
	err = service.LevelLog().Delete(ctx, req.ID)
	return
}

// BatchDelete 批量删除等级变更日志
func (c *cLevelLog) BatchDelete(ctx context.Context, req *v1.LevelLogBatchDeleteReq) (res *v1.LevelLogBatchDeleteRes, err error) {
	err = service.LevelLog().BatchDelete(ctx, req.IDs)
	return
}

// Detail 获取等级变更日志详情
func (c *cLevelLog) Detail(ctx context.Context, req *v1.LevelLogDetailReq) (res *v1.LevelLogDetailRes, err error) {
	res = &v1.LevelLogDetailRes{}
	res.LevelLogDetailOutput, err = service.LevelLog().Detail(ctx, req.ID)
	return
}

// List 获取等级变更日志列表
func (c *cLevelLog) List(ctx context.Context, req *v1.LevelLogListReq) (res *v1.LevelLogListRes, err error) {
	res = &v1.LevelLogListRes{}
	res.List, res.Total, err = service.LevelLog().List(ctx, &model.LevelLogListInput{
		PageNum:   req.PageNum,
		PageSize:  req.PageSize,
		OrderBy:   req.OrderBy,
		OrderDir:  req.OrderDir,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		UserID: req.UserID,
		OldLevelID: req.OldLevelID,
		NewLevelID: req.NewLevelID,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
		ChangeType: req.ChangeType,
		ExpireAtStart: req.ExpireAtStart,
		ExpireAtEnd: req.ExpireAtEnd,
	})
	return
}
// Export 导出等级变更日志
func (c *cLevelLog) Export(ctx context.Context, req *v1.LevelLogExportReq) (res *v1.LevelLogExportRes, err error) {
	list, err := service.LevelLog().Export(ctx, &model.LevelLogListInput{
		OrderBy:   req.OrderBy,
		OrderDir:  req.OrderDir,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		UserID: req.UserID,
		OldLevelID: req.OldLevelID,
		NewLevelID: req.NewLevelID,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
		ChangeType: req.ChangeType,
		ExpireAtStart: req.ExpireAtStart,
		ExpireAtEnd: req.ExpireAtEnd,
	})
	if err != nil {
		return
	}
	// CSV 导出（使用 csv.Writer 防止注入和格式问题）
	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Response.Header().Set("Content-Disposition", `attachment; filename="level_log.csv"`)
	r.Response.Write("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(r.Response.Writer)
	// 表头（与导入模板列对齐，末尾追加只读列）
	_ = w.Write([]string{"会员", "变更前等级", "变更后等级", "变更类型", "变更说明", "创建时间"})
	// 数据行
	for _, item := range list {
		_ = w.Write([]string{
			csvSafeLevelLog(item.UserNickname),
			csvSafeLevelLog(item.LevelName),
			csvSafeLevelLog(item.NewLevelName),
			fmt.Sprintf("%v", item.ChangeType),
			csvSafeLevelLog(item.Remark),
			func() string { if item.CreatedAt != nil { return item.CreatedAt.String() }; return "" }(),
		})
	}
	w.Flush()
	return
}

// Import 导入等级变更日志
func (c *cLevelLog) Import(ctx context.Context, req *v1.LevelLogImportReq) (res *v1.LevelLogImportRes, err error) {
	r := g.RequestFromCtx(ctx)
	file := r.GetUploadFile("file")
	if file == nil {
		return nil, fmt.Errorf("请上传文件")
	}
	success, fail, err := service.LevelLog().Import(ctx, file)
	if err != nil {
		return nil, err
	}
	res = &v1.LevelLogImportRes{Success: success, Fail: fail}
	return
}

// ImportTemplate 下载等级变更日志导入模板
func (c *cLevelLog) ImportTemplate(ctx context.Context, req *v1.LevelLogImportTemplateReq) (res *v1.LevelLogImportTemplateRes, err error) {
	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Response.Header().Set("Content-Disposition", `attachment; filename="level_log_template.csv"`)
	r.Response.Write("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(r.Response.Writer)
	_ = w.Write([]string{"会员", "变更前等级", "变更后等级", "变更类型", "变更说明"})
	w.Flush()
	return
}
