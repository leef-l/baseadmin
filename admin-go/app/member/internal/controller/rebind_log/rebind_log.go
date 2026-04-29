package rebind_log

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

func csvSafeRebindLog(s string) string {
	if len(s) == 0 {
		return s
	}
	if s[0] == 0 || strings.ContainsAny(s[:1], "=+-@\t\r") {
		return "'" + s
	}
	return s
}

var RebindLog = cRebindLog{}

type cRebindLog struct{}

// Create 创建换绑上级日志
func (c *cRebindLog) Create(ctx context.Context, req *v1.RebindLogCreateReq) (res *v1.RebindLogCreateRes, err error) {
	err = service.RebindLog().Create(ctx, &model.RebindLogCreateInput{
		UserID: req.UserID,
		OldParentID: req.OldParentID,
		NewParentID: req.NewParentID,
		Reason: req.Reason,
		OperatorID: req.OperatorID,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
	})
	return
}

// Update 更新换绑上级日志
func (c *cRebindLog) Update(ctx context.Context, req *v1.RebindLogUpdateReq) (res *v1.RebindLogUpdateRes, err error) {
	err = service.RebindLog().Update(ctx, &model.RebindLogUpdateInput{
		ID: req.ID,
		UserID: req.UserID,
		OldParentID: req.OldParentID,
		NewParentID: req.NewParentID,
		Reason: req.Reason,
		OperatorID: req.OperatorID,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
	})
	return
}

// Delete 删除换绑上级日志
func (c *cRebindLog) Delete(ctx context.Context, req *v1.RebindLogDeleteReq) (res *v1.RebindLogDeleteRes, err error) {
	err = service.RebindLog().Delete(ctx, req.ID)
	return
}

// BatchDelete 批量删除换绑上级日志
func (c *cRebindLog) BatchDelete(ctx context.Context, req *v1.RebindLogBatchDeleteReq) (res *v1.RebindLogBatchDeleteRes, err error) {
	err = service.RebindLog().BatchDelete(ctx, req.IDs)
	return
}

// Detail 获取换绑上级日志详情
func (c *cRebindLog) Detail(ctx context.Context, req *v1.RebindLogDetailReq) (res *v1.RebindLogDetailRes, err error) {
	res = &v1.RebindLogDetailRes{}
	res.RebindLogDetailOutput, err = service.RebindLog().Detail(ctx, req.ID)
	return
}

// List 获取换绑上级日志列表
func (c *cRebindLog) List(ctx context.Context, req *v1.RebindLogListReq) (res *v1.RebindLogListRes, err error) {
	res = &v1.RebindLogListRes{}
	res.List, res.Total, err = service.RebindLog().List(ctx, &model.RebindLogListInput{
		PageNum:   req.PageNum,
		PageSize:  req.PageSize,
		OrderBy:   req.OrderBy,
		OrderDir:  req.OrderDir,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		UserID: req.UserID,
		OldParentID: req.OldParentID,
		NewParentID: req.NewParentID,
		OperatorID: req.OperatorID,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
	})
	return
}
// Export 导出换绑上级日志
func (c *cRebindLog) Export(ctx context.Context, req *v1.RebindLogExportReq) (res *v1.RebindLogExportRes, err error) {
	list, err := service.RebindLog().Export(ctx, &model.RebindLogListInput{
		OrderBy:   req.OrderBy,
		OrderDir:  req.OrderDir,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		UserID: req.UserID,
		OldParentID: req.OldParentID,
		NewParentID: req.NewParentID,
		OperatorID: req.OperatorID,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
	})
	if err != nil {
		return
	}
	// CSV 导出（使用 csv.Writer 防止注入和格式问题）
	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Response.Header().Set("Content-Disposition", `attachment; filename="rebind_log.csv"`)
	r.Response.Write("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(r.Response.Writer)
	// 表头（与导入模板列对齐，末尾追加只读列）
	_ = w.Write([]string{"会员", "原上级", "新上级", "换绑原因", "操作人", "创建时间"})
	// 数据行
	for _, item := range list {
		_ = w.Write([]string{
			csvSafeRebindLog(item.UserNickname),
			csvSafeRebindLog(item.OldParentNickname),
			csvSafeRebindLog(item.NewParentNickname),
			csvSafeRebindLog(item.Reason),
			csvSafeRebindLog(item.UsersUsername),
			func() string { if item.CreatedAt != nil { return item.CreatedAt.String() }; return "" }(),
		})
	}
	w.Flush()
	return
}

// Import 导入换绑上级日志
func (c *cRebindLog) Import(ctx context.Context, req *v1.RebindLogImportReq) (res *v1.RebindLogImportRes, err error) {
	r := g.RequestFromCtx(ctx)
	file := r.GetUploadFile("file")
	if file == nil {
		return nil, fmt.Errorf("请上传文件")
	}
	success, fail, err := service.RebindLog().Import(ctx, file)
	if err != nil {
		return nil, err
	}
	res = &v1.RebindLogImportRes{Success: success, Fail: fail}
	return
}

// ImportTemplate 下载换绑上级日志导入模板
func (c *cRebindLog) ImportTemplate(ctx context.Context, req *v1.RebindLogImportTemplateReq) (res *v1.RebindLogImportTemplateRes, err error) {
	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Response.Header().Set("Content-Disposition", `attachment; filename="rebind_log_template.csv"`)
	r.Response.Write("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(r.Response.Writer)
	_ = w.Write([]string{"会员", "原上级", "新上级", "换绑原因", "操作人"})
	w.Flush()
	return
}
