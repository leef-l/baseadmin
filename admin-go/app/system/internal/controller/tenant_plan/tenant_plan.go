package tenant_plan

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

func csvSafeTenantPlan(s string) string {
	if len(s) == 0 {
		return s
	}
	if s[0] == 0 || strings.ContainsAny(s[:1], "=+-@\t\r") {
		return "'" + s
	}
	return s
}

var TenantPlan = cTenantPlan{}

type cTenantPlan struct{}

// Create 创建租户套餐订阅表
func (c *cTenantPlan) Create(ctx context.Context, req *v1.TenantPlanCreateReq) (res *v1.TenantPlanCreateRes, err error) {
	err = service.TenantPlan().Create(ctx, &model.TenantPlanCreateInput{
		TenantID: req.TenantID,
		PlanID: req.PlanID,
		StartAt: req.StartAt,
		ExpireAt: req.ExpireAt,
		Status: req.Status,
		Remark: req.Remark,
		MerchantID: req.MerchantID,
	})
	return
}

// Update 更新租户套餐订阅表
func (c *cTenantPlan) Update(ctx context.Context, req *v1.TenantPlanUpdateReq) (res *v1.TenantPlanUpdateRes, err error) {
	err = service.TenantPlan().Update(ctx, &model.TenantPlanUpdateInput{
		ID: req.ID,
		TenantID: req.TenantID,
		PlanID: req.PlanID,
		StartAt: req.StartAt,
		ExpireAt: req.ExpireAt,
		Status: req.Status,
		Remark: req.Remark,
		MerchantID: req.MerchantID,
	})
	return
}

// Delete 删除租户套餐订阅表
func (c *cTenantPlan) Delete(ctx context.Context, req *v1.TenantPlanDeleteReq) (res *v1.TenantPlanDeleteRes, err error) {
	err = service.TenantPlan().Delete(ctx, req.ID)
	return
}

// BatchDelete 批量删除租户套餐订阅表
func (c *cTenantPlan) BatchDelete(ctx context.Context, req *v1.TenantPlanBatchDeleteReq) (res *v1.TenantPlanBatchDeleteRes, err error) {
	err = service.TenantPlan().BatchDelete(ctx, req.IDs)
	return
}

// BatchUpdate 批量编辑租户套餐订阅表
func (c *cTenantPlan) BatchUpdate(ctx context.Context, req *v1.TenantPlanBatchUpdateReq) (res *v1.TenantPlanBatchUpdateRes, err error) {
	err = service.TenantPlan().BatchUpdate(ctx, &model.TenantPlanBatchUpdateInput{
		IDs: req.IDs,
		Status: req.Status,
	})
	return
}

// Detail 获取租户套餐订阅表详情
func (c *cTenantPlan) Detail(ctx context.Context, req *v1.TenantPlanDetailReq) (res *v1.TenantPlanDetailRes, err error) {
	res = &v1.TenantPlanDetailRes{}
	res.TenantPlanDetailOutput, err = service.TenantPlan().Detail(ctx, req.ID)
	return
}

// List 获取租户套餐订阅表列表
func (c *cTenantPlan) List(ctx context.Context, req *v1.TenantPlanListReq) (res *v1.TenantPlanListRes, err error) {
	res = &v1.TenantPlanListRes{}
	res.List, res.Total, err = service.TenantPlan().List(ctx, &model.TenantPlanListInput{
		PageNum:   req.PageNum,
		PageSize:  req.PageSize,
		OrderBy:   req.OrderBy,
		OrderDir:  req.OrderDir,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		TenantID: req.TenantID,
		PlanID: req.PlanID,
		MerchantID: req.MerchantID,
		Status: req.Status,
		Remark: req.Remark,
		StartAtStart: req.StartAtStart,
		StartAtEnd: req.StartAtEnd,
		ExpireAtStart: req.ExpireAtStart,
		ExpireAtEnd: req.ExpireAtEnd,
	})
	return
}
// Export 导出租户套餐订阅表
func (c *cTenantPlan) Export(ctx context.Context, req *v1.TenantPlanExportReq) (res *v1.TenantPlanExportRes, err error) {
	list, err := service.TenantPlan().Export(ctx, &model.TenantPlanListInput{
		OrderBy:   req.OrderBy,
		OrderDir:  req.OrderDir,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		TenantID: req.TenantID,
		PlanID: req.PlanID,
		MerchantID: req.MerchantID,
		Status: req.Status,
		Remark: req.Remark,
		StartAtStart: req.StartAtStart,
		StartAtEnd: req.StartAtEnd,
		ExpireAtStart: req.ExpireAtStart,
		ExpireAtEnd: req.ExpireAtEnd,
	})
	if err != nil {
		return
	}
	// CSV 导出（使用 csv.Writer 防止注入和格式问题）
	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Response.Header().Set("Content-Disposition", `attachment; filename="tenant_plan.csv"`)
	r.Response.Write("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(r.Response.Writer)
	// 表头（与导入模板列对齐，末尾追加只读列）
	_ = w.Write([]string{"套餐ID", "状态", "备注", "创建时间"})
	// 数据行
	for _, item := range list {
		_ = w.Write([]string{
			csvSafeTenantPlan(item.PlanName),
			fmt.Sprintf("%v", item.Status),
			csvSafeTenantPlan(item.Remark),
			func() string { if item.CreatedAt != nil { return item.CreatedAt.String() }; return "" }(),
		})
	}
	w.Flush()
	return
}

// Import 导入租户套餐订阅表
func (c *cTenantPlan) Import(ctx context.Context, req *v1.TenantPlanImportReq) (res *v1.TenantPlanImportRes, err error) {
	r := g.RequestFromCtx(ctx)
	file := r.GetUploadFile("file")
	if file == nil {
		return nil, fmt.Errorf("请上传文件")
	}
	success, fail, err := service.TenantPlan().Import(ctx, file)
	if err != nil {
		return nil, err
	}
	res = &v1.TenantPlanImportRes{Success: success, Fail: fail}
	return
}

// ImportTemplate 下载租户套餐订阅表导入模板
func (c *cTenantPlan) ImportTemplate(ctx context.Context, req *v1.TenantPlanImportTemplateReq) (res *v1.TenantPlanImportTemplateRes, err error) {
	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Response.Header().Set("Content-Disposition", `attachment; filename="tenant_plan_template.csv"`)
	r.Response.Write("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(r.Response.Writer)
	_ = w.Write([]string{"套餐ID", "状态", "备注"})
	w.Flush()
	return
}
