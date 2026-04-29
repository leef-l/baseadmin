package audit_log

import (
	"context"
	"encoding/csv"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"

	v1 "gbaseadmin/app/demo/api/demo/v1"
	"gbaseadmin/app/demo/internal/model"
	"gbaseadmin/app/demo/internal/service"
)

var AuditLog = cAuditLog{}

type cAuditLog struct{}

// Create 创建体验审计日志
func (c *cAuditLog) Create(ctx context.Context, req *v1.AuditLogCreateReq) (res *v1.AuditLogCreateRes, err error) {
	err = service.AuditLog().Create(ctx, &model.AuditLogCreateInput{
		LogNo: req.LogNo,
		OperatorID: req.OperatorID,
		Action: req.Action,
		TargetType: req.TargetType,
		TargetCode: req.TargetCode,
		RequestJSON: req.RequestJSON,
		Result: req.Result,
		ClientIP: req.ClientIP,
		OccurredAt: req.OccurredAt,
		Remark: req.Remark,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
	})
	return
}

// Update 更新体验审计日志
func (c *cAuditLog) Update(ctx context.Context, req *v1.AuditLogUpdateReq) (res *v1.AuditLogUpdateRes, err error) {
	err = service.AuditLog().Update(ctx, &model.AuditLogUpdateInput{
		ID: req.ID,
		LogNo: req.LogNo,
		OperatorID: req.OperatorID,
		Action: req.Action,
		TargetType: req.TargetType,
		TargetCode: req.TargetCode,
		RequestJSON: req.RequestJSON,
		Result: req.Result,
		ClientIP: req.ClientIP,
		OccurredAt: req.OccurredAt,
		Remark: req.Remark,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
	})
	return
}

// Delete 删除体验审计日志
func (c *cAuditLog) Delete(ctx context.Context, req *v1.AuditLogDeleteReq) (res *v1.AuditLogDeleteRes, err error) {
	err = service.AuditLog().Delete(ctx, req.ID)
	return
}

// BatchDelete 批量删除体验审计日志
func (c *cAuditLog) BatchDelete(ctx context.Context, req *v1.AuditLogBatchDeleteReq) (res *v1.AuditLogBatchDeleteRes, err error) {
	err = service.AuditLog().BatchDelete(ctx, req.IDs)
	return
}

// Detail 获取体验审计日志详情
func (c *cAuditLog) Detail(ctx context.Context, req *v1.AuditLogDetailReq) (res *v1.AuditLogDetailRes, err error) {
	res = &v1.AuditLogDetailRes{}
	res.AuditLogDetailOutput, err = service.AuditLog().Detail(ctx, req.ID)
	return
}

// List 获取体验审计日志列表
func (c *cAuditLog) List(ctx context.Context, req *v1.AuditLogListReq) (res *v1.AuditLogListRes, err error) {
	res = &v1.AuditLogListRes{}
	res.List, res.Total, err = service.AuditLog().List(ctx, &model.AuditLogListInput{
		PageNum:   req.PageNum,
		PageSize:  req.PageSize,
		OrderBy:   req.OrderBy,
		OrderDir:  req.OrderDir,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Keyword: req.Keyword,
		LogNo: req.LogNo,
		TargetCode: req.TargetCode,
		OperatorID: req.OperatorID,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
		ClientIP: req.ClientIP,
		Action: req.Action,
		TargetType: req.TargetType,
		Result: req.Result,
		OccurredAtStart: req.OccurredAtStart,
		OccurredAtEnd: req.OccurredAtEnd,
	})
	return
}
// Export 导出体验审计日志
func (c *cAuditLog) Export(ctx context.Context, req *v1.AuditLogExportReq) (res *v1.AuditLogExportRes, err error) {
	list, err := service.AuditLog().Export(ctx, &model.AuditLogListInput{
		OrderBy:   req.OrderBy,
		OrderDir:  req.OrderDir,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Keyword: req.Keyword,
		LogNo: req.LogNo,
		TargetCode: req.TargetCode,
		OperatorID: req.OperatorID,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
		ClientIP: req.ClientIP,
		Action: req.Action,
		TargetType: req.TargetType,
		Result: req.Result,
		OccurredAtStart: req.OccurredAtStart,
		OccurredAtEnd: req.OccurredAtEnd,
	})
	if err != nil {
		return
	}
	// CSV 导出（使用 csv.Writer 防止注入和格式问题）
	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Response.Header().Set("Content-Disposition", `attachment; filename="audit_log.csv"`)
	r.Response.Write("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(r.Response.Writer)
	// 表头
	_ = w.Write([]string{"日志编号", "操作人", "动作", "对象类型", "对象编号", "请求JSON", "结果", "客户端IP", "发生时间", "备注", "租户", "商户", "创建时间"})
	// 数据行
	for _, item := range list {
		_ = w.Write([]string{
			item.LogNo,
			item.UsersUsername,
			fmt.Sprintf("%v", item.Action),
			fmt.Sprintf("%v", item.TargetType),
			item.TargetCode,
			item.RequestJSON,
			fmt.Sprintf("%v", item.Result),
			item.ClientIP,
			func() string { if item.OccurredAt != nil { return item.OccurredAt.String() }; return "" }(),
			item.Remark,
			item.TenantName,
			item.MerchantName,
			func() string { if item.CreatedAt != nil { return item.CreatedAt.String() }; return "" }(),
		})
	}
	w.Flush()
	return
}

// Import 导入体验审计日志
func (c *cAuditLog) Import(ctx context.Context, req *v1.AuditLogImportReq) (res *v1.AuditLogImportRes, err error) {
	r := g.RequestFromCtx(ctx)
	file := r.GetUploadFile("file")
	if file == nil {
		return nil, fmt.Errorf("请上传文件")
	}
	success, fail, err := service.AuditLog().Import(ctx, file)
	if err != nil {
		return nil, err
	}
	res = &v1.AuditLogImportRes{Success: success, Fail: fail}
	return
}

// ImportTemplate 下载体验审计日志导入模板
func (c *cAuditLog) ImportTemplate(ctx context.Context, req *v1.AuditLogImportTemplateReq) (res *v1.AuditLogImportTemplateRes, err error) {
	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Response.Header().Set("Content-Disposition", `attachment; filename="audit_log_template.csv"`)
	r.Response.Write("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(r.Response.Writer)
	_ = w.Write([]string{"日志编号", "操作人", "动作", "对象类型", "对象编号", "请求JSON", "结果", "客户端IP", "备注", "租户", "商户"})
	w.Flush()
	return
}
