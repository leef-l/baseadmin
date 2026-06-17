package plan

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

func csvSafePlan(s string) string {
	if len(s) == 0 {
		return s
	}
	if s[0] == 0 || strings.ContainsAny(s[:1], "=+-@\t\r") {
		return "'" + s
	}
	return s
}

var Plan = cPlan{}

type cPlan struct{}

// Create 创建套餐表
func (c *cPlan) Create(ctx context.Context, req *v1.PlanCreateReq) (res *v1.PlanCreateRes, err error) {
	err = service.Plan().Create(ctx, &model.PlanCreateInput{
		Name: req.Name,
		Code: req.Code,
		Description: req.Description,
		UserLimit: req.UserLimit,
		StorageMb: req.StorageMb,
		Features: req.Features,
		Price: req.Price,
		Sort: req.Sort,
		Status: req.Status,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
	})
	return
}

// Update 更新套餐表
func (c *cPlan) Update(ctx context.Context, req *v1.PlanUpdateReq) (res *v1.PlanUpdateRes, err error) {
	err = service.Plan().Update(ctx, &model.PlanUpdateInput{
		ID: req.ID,
		Name: req.Name,
		Code: req.Code,
		Description: req.Description,
		UserLimit: req.UserLimit,
		StorageMb: req.StorageMb,
		Features: req.Features,
		Price: req.Price,
		Sort: req.Sort,
		Status: req.Status,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
	})
	return
}

// Delete 删除套餐表
func (c *cPlan) Delete(ctx context.Context, req *v1.PlanDeleteReq) (res *v1.PlanDeleteRes, err error) {
	err = service.Plan().Delete(ctx, req.ID)
	return
}

// BatchDelete 批量删除套餐表
func (c *cPlan) BatchDelete(ctx context.Context, req *v1.PlanBatchDeleteReq) (res *v1.PlanBatchDeleteRes, err error) {
	err = service.Plan().BatchDelete(ctx, req.IDs)
	return
}

// BatchUpdate 批量编辑套餐表
func (c *cPlan) BatchUpdate(ctx context.Context, req *v1.PlanBatchUpdateReq) (res *v1.PlanBatchUpdateRes, err error) {
	err = service.Plan().BatchUpdate(ctx, &model.PlanBatchUpdateInput{
		IDs: req.IDs,
		Status: req.Status,
	})
	return
}

// Detail 获取套餐表详情
func (c *cPlan) Detail(ctx context.Context, req *v1.PlanDetailReq) (res *v1.PlanDetailRes, err error) {
	res = &v1.PlanDetailRes{}
	res.PlanDetailOutput, err = service.Plan().Detail(ctx, req.ID)
	return
}

// List 获取套餐表列表
func (c *cPlan) List(ctx context.Context, req *v1.PlanListReq) (res *v1.PlanListRes, err error) {
	res = &v1.PlanListRes{}
	res.List, res.Total, err = service.Plan().List(ctx, &model.PlanListInput{
		PageNum:   req.PageNum,
		PageSize:  req.PageSize,
		OrderBy:   req.OrderBy,
		OrderDir:  req.OrderDir,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Keyword: req.Keyword,
		Code: req.Code,
		Name: req.Name,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
		Status: req.Status,
		Description: req.Description,
	})
	return
}
// Export 导出套餐表
func (c *cPlan) Export(ctx context.Context, req *v1.PlanExportReq) (res *v1.PlanExportRes, err error) {
	list, err := service.Plan().Export(ctx, &model.PlanListInput{
		OrderBy:   req.OrderBy,
		OrderDir:  req.OrderDir,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Keyword: req.Keyword,
		Code: req.Code,
		Name: req.Name,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
		Status: req.Status,
		Description: req.Description,
	})
	if err != nil {
		return
	}
	// CSV 导出（使用 csv.Writer 防止注入和格式问题）
	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Response.Header().Set("Content-Disposition", `attachment; filename="plan.csv"`)
	r.Response.Write("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(r.Response.Writer)
	// 表头（与导入模板列对齐，末尾追加只读列）
	_ = w.Write([]string{"套餐名称", "套餐编码", "套餐描述", "用户数量上限", "存储空间上限MB", "功能开关列表", "价格", "排序", "状态", "创建时间"})
	// 数据行
	for _, item := range list {
		_ = w.Write([]string{
			csvSafePlan(item.Name),
			csvSafePlan(item.Code),
			csvSafePlan(item.Description),
			fmt.Sprintf("%v", item.UserLimit),
			fmt.Sprintf("%v", item.StorageMb),
			csvSafePlan(item.Features),
			fmt.Sprintf("%.2f", float64(item.Price)/100),
			fmt.Sprintf("%v", item.Sort),
			fmt.Sprintf("%v", item.Status),
			func() string { if item.CreatedAt != nil { return item.CreatedAt.String() }; return "" }(),
		})
	}
	w.Flush()
	return
}

// Import 导入套餐表
func (c *cPlan) Import(ctx context.Context, req *v1.PlanImportReq) (res *v1.PlanImportRes, err error) {
	r := g.RequestFromCtx(ctx)
	file := r.GetUploadFile("file")
	if file == nil {
		return nil, fmt.Errorf("请上传文件")
	}
	success, fail, err := service.Plan().Import(ctx, file)
	if err != nil {
		return nil, err
	}
	res = &v1.PlanImportRes{Success: success, Fail: fail}
	return
}

// ImportTemplate 下载套餐表导入模板
func (c *cPlan) ImportTemplate(ctx context.Context, req *v1.PlanImportTemplateReq) (res *v1.PlanImportTemplateRes, err error) {
	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Response.Header().Set("Content-Disposition", `attachment; filename="plan_template.csv"`)
	r.Response.Write("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(r.Response.Writer)
	_ = w.Write([]string{"套餐名称", "套餐编码", "套餐描述", "用户数量上限", "存储空间上限MB", "功能开关列表", "价格", "排序", "状态"})
	w.Flush()
	return
}
