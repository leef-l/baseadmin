package category

import (
	"context"
	"encoding/csv"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"

	v1 "gbaseadmin/app/demo/api/demo/v1"
	"gbaseadmin/app/demo/internal/model"
	"gbaseadmin/app/demo/internal/service"
)

var Category = cCategory{}

type cCategory struct{}

// Create 创建体验分类
func (c *cCategory) Create(ctx context.Context, req *v1.CategoryCreateReq) (res *v1.CategoryCreateRes, err error) {
	err = service.Category().Create(ctx, &model.CategoryCreateInput{
		ParentID: req.ParentID,
		Name: req.Name,
		Icon: req.Icon,
		Sort: req.Sort,
		Status: req.Status,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
	})
	return
}

// Update 更新体验分类
func (c *cCategory) Update(ctx context.Context, req *v1.CategoryUpdateReq) (res *v1.CategoryUpdateRes, err error) {
	err = service.Category().Update(ctx, &model.CategoryUpdateInput{
		ID: req.ID,
		ParentID: req.ParentID,
		Name: req.Name,
		Icon: req.Icon,
		Sort: req.Sort,
		Status: req.Status,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
	})
	return
}

// Delete 删除体验分类
func (c *cCategory) Delete(ctx context.Context, req *v1.CategoryDeleteReq) (res *v1.CategoryDeleteRes, err error) {
	err = service.Category().Delete(ctx, req.ID)
	return
}

// BatchDelete 批量删除体验分类
func (c *cCategory) BatchDelete(ctx context.Context, req *v1.CategoryBatchDeleteReq) (res *v1.CategoryBatchDeleteRes, err error) {
	err = service.Category().BatchDelete(ctx, req.IDs)
	return
}

// BatchUpdate 批量编辑体验分类
func (c *cCategory) BatchUpdate(ctx context.Context, req *v1.CategoryBatchUpdateReq) (res *v1.CategoryBatchUpdateRes, err error) {
	err = service.Category().BatchUpdate(ctx, &model.CategoryBatchUpdateInput{
		IDs: req.IDs,
		Status: req.Status,
	})
	return
}

// Detail 获取体验分类详情
func (c *cCategory) Detail(ctx context.Context, req *v1.CategoryDetailReq) (res *v1.CategoryDetailRes, err error) {
	res = &v1.CategoryDetailRes{}
	res.CategoryDetailOutput, err = service.Category().Detail(ctx, req.ID)
	return
}

// List 获取体验分类列表
func (c *cCategory) List(ctx context.Context, req *v1.CategoryListReq) (res *v1.CategoryListRes, err error) {
	res = &v1.CategoryListRes{}
	res.List, res.Total, err = service.Category().List(ctx, &model.CategoryListInput{
		PageNum:   req.PageNum,
		PageSize:  req.PageSize,
		OrderBy:   req.OrderBy,
		OrderDir:  req.OrderDir,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Name: req.Name,
		ParentID: req.ParentID,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
		Status: req.Status,
	})
	return
}
// Export 导出体验分类
func (c *cCategory) Export(ctx context.Context, req *v1.CategoryExportReq) (res *v1.CategoryExportRes, err error) {
	list, err := service.Category().Export(ctx, &model.CategoryListInput{
		OrderBy:   req.OrderBy,
		OrderDir:  req.OrderDir,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Name: req.Name,
		ParentID: req.ParentID,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
		Status: req.Status,
	})
	if err != nil {
		return
	}
	// CSV 导出（使用 csv.Writer 防止注入和格式问题）
	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Response.Header().Set("Content-Disposition", `attachment; filename="category.csv"`)
	r.Response.Write("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(r.Response.Writer)
	// 表头
	_ = w.Write([]string{"父分类", "分类名称", "图标", "排序", "状态", "租户", "商户", "创建时间"})
	// 数据行
	for _, item := range list {
		_ = w.Write([]string{
			item.CategoryName,
			item.Name,
			item.Icon,
			fmt.Sprintf("%v", item.Sort),
			fmt.Sprintf("%v", item.Status),
			item.TenantName,
			item.MerchantName,
			func() string { if item.CreatedAt != nil { return item.CreatedAt.String() }; return "" }(),
		})
	}
	w.Flush()
	return
}

// Tree 获取体验分类树形结构
func (c *cCategory) Tree(ctx context.Context, req *v1.CategoryTreeReq) (res *v1.CategoryTreeRes, err error) {
	res = &v1.CategoryTreeRes{}
	res.List, err = service.Category().Tree(ctx, &model.CategoryTreeInput{
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Name: req.Name,
		ParentID: req.ParentID,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
		Status: req.Status,
	})
	return
}

