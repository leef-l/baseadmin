package cron

import (
	"context"
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/frame/g"

	v1 "gbaseadmin/app/system/api/system/v1"
	"gbaseadmin/app/system/internal/model"
	"gbaseadmin/app/system/internal/service"
)

func csvSafeCron(s string) string {
	if len(s) == 0 {
		return s
	}
	if s[0] == 0 || strings.ContainsAny(s[:1], "=+-@\t\r") {
		return "'" + s
	}
	return s
}

var Cron = cCron{}

type cCron struct{}

// Create 创建定时任务表
func (c *cCron) Create(ctx context.Context, req *v1.CronCreateReq) (res *v1.CronCreateRes, err error) {
	err = service.Cron().Create(ctx, &model.CronCreateInput{
		Name: req.Name,
		Expression: req.Expression,
		Handler: req.Handler,
		Params: req.Params,
		Remark: req.Remark,
		Status: req.Status,
		LastRunAt: req.LastRunAt,
		LastResult: req.LastResult,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
	})
	return
}

// Update 更新定时任务表
func (c *cCron) Update(ctx context.Context, req *v1.CronUpdateReq) (res *v1.CronUpdateRes, err error) {
	err = service.Cron().Update(ctx, &model.CronUpdateInput{
		ID: req.ID,
		Name: req.Name,
		Expression: req.Expression,
		Handler: req.Handler,
		Params: req.Params,
		Remark: req.Remark,
		Status: req.Status,
		LastRunAt: req.LastRunAt,
		LastResult: req.LastResult,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
	})
	return
}

// Delete 删除定时任务表
func (c *cCron) Delete(ctx context.Context, req *v1.CronDeleteReq) (res *v1.CronDeleteRes, err error) {
	err = service.Cron().Delete(ctx, req.ID)
	return
}

// BatchDelete 批量删除定时任务表
func (c *cCron) BatchDelete(ctx context.Context, req *v1.CronBatchDeleteReq) (res *v1.CronBatchDeleteRes, err error) {
	err = service.Cron().BatchDelete(ctx, req.IDs)
	return
}

// BatchUpdate 批量编辑定时任务表
func (c *cCron) BatchUpdate(ctx context.Context, req *v1.CronBatchUpdateReq) (res *v1.CronBatchUpdateRes, err error) {
	err = service.Cron().BatchUpdate(ctx, &model.CronBatchUpdateInput{
		IDs: req.IDs,
		Status: req.Status,
	})
	return
}

// Detail 获取定时任务表详情
func (c *cCron) Detail(ctx context.Context, req *v1.CronDetailReq) (res *v1.CronDetailRes, err error) {
	res = &v1.CronDetailRes{}
	res.CronDetailOutput, err = service.Cron().Detail(ctx, req.ID)
	return
}

// List 获取定时任务表列表
func (c *cCron) List(ctx context.Context, req *v1.CronListReq) (res *v1.CronListRes, err error) {
	res = &v1.CronListRes{}
	res.List, res.Total, err = service.Cron().List(ctx, &model.CronListInput{
		PageNum:   req.PageNum,
		PageSize:  req.PageSize,
		OrderBy:   req.OrderBy,
		OrderDir:  req.OrderDir,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Keyword: req.Keyword,
		Name: req.Name,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
		Status: req.Status,
		Remark: req.Remark,
		LastRunAtStart: req.LastRunAtStart,
		LastRunAtEnd: req.LastRunAtEnd,
	})
	return
}
// Export 导出定时任务表
func (c *cCron) Export(ctx context.Context, req *v1.CronExportReq) (res *v1.CronExportRes, err error) {
	list, err := service.Cron().Export(ctx, &model.CronListInput{
		OrderBy:   req.OrderBy,
		OrderDir:  req.OrderDir,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Keyword: req.Keyword,
		Name: req.Name,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
		Status: req.Status,
		Remark: req.Remark,
		LastRunAtStart: req.LastRunAtStart,
		LastRunAtEnd: req.LastRunAtEnd,
	})
	if err != nil {
		return
	}
	// CSV 导出（使用 csv.Writer 防止注入和格式问题）
	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Response.Header().Set("Content-Disposition", `attachment; filename="cron.csv"`)
	r.Response.Write("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(r.Response.Writer)
	// 表头（与导入模板列对齐，末尾追加只读列）
	_ = w.Write([]string{"任务名称", "Cron表达式", "处理器名称", "参数", "备注", "状态", "上次执行结果", "创建时间"})
	// 数据行
	for _, item := range list {
		_ = w.Write([]string{
			csvSafeCron(item.Name),
			csvSafeCron(item.Expression),
			csvSafeCron(item.Handler),
			csvSafeCron(item.Params),
			csvSafeCron(item.Remark),
			fmt.Sprintf("%v", item.Status),
			csvSafeCron(item.LastResult),
			func() string { if item.CreatedAt != nil { return item.CreatedAt.String() }; return "" }(),
		})
	}
	w.Flush()
	return
}

// Import 导入定时任务表
func (c *cCron) Import(ctx context.Context, req *v1.CronImportReq) (res *v1.CronImportRes, err error) {
	r := g.RequestFromCtx(ctx)
	file := r.GetUploadFile("file")
	if file == nil {
		return nil, fmt.Errorf("请上传文件")
	}
	success, fail, err := service.Cron().Import(ctx, file)
	if err != nil {
		return nil, err
	}
	res = &v1.CronImportRes{Success: success, Fail: fail}
	return
}

// ImportTemplate 下载定时任务表导入模板
func (c *cCron) ImportTemplate(ctx context.Context, req *v1.CronImportTemplateReq) (res *v1.CronImportTemplateRes, err error) {
	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Response.Header().Set("Content-Disposition", `attachment; filename="cron_template.csv"`)
	r.Response.Write("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(r.Response.Writer)
	_ = w.Write([]string{"任务名称", "Cron表达式", "处理器名称", "参数", "备注", "状态", "上次执行结果"})
	w.Flush()
	return
}
