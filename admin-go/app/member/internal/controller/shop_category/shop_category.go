package shop_category

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

func csvSafeShopCategory(s string) string {
	if len(s) == 0 {
		return s
	}
	if s[0] == 0 || strings.ContainsAny(s[:1], "=+-@\t\r") {
		return "'" + s
	}
	return s
}

var ShopCategory = cShopCategory{}

type cShopCategory struct{}

// Create 创建商城商品分类
func (c *cShopCategory) Create(ctx context.Context, req *v1.ShopCategoryCreateReq) (res *v1.ShopCategoryCreateRes, err error) {
	err = service.ShopCategory().Create(ctx, &model.ShopCategoryCreateInput{
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

// Update 更新商城商品分类
func (c *cShopCategory) Update(ctx context.Context, req *v1.ShopCategoryUpdateReq) (res *v1.ShopCategoryUpdateRes, err error) {
	err = service.ShopCategory().Update(ctx, &model.ShopCategoryUpdateInput{
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

// Delete 删除商城商品分类
func (c *cShopCategory) Delete(ctx context.Context, req *v1.ShopCategoryDeleteReq) (res *v1.ShopCategoryDeleteRes, err error) {
	err = service.ShopCategory().Delete(ctx, req.ID)
	return
}

// BatchDelete 批量删除商城商品分类
func (c *cShopCategory) BatchDelete(ctx context.Context, req *v1.ShopCategoryBatchDeleteReq) (res *v1.ShopCategoryBatchDeleteRes, err error) {
	err = service.ShopCategory().BatchDelete(ctx, req.IDs)
	return
}

// BatchUpdate 批量编辑商城商品分类
func (c *cShopCategory) BatchUpdate(ctx context.Context, req *v1.ShopCategoryBatchUpdateReq) (res *v1.ShopCategoryBatchUpdateRes, err error) {
	err = service.ShopCategory().BatchUpdate(ctx, &model.ShopCategoryBatchUpdateInput{
		IDs: req.IDs,
		Status: req.Status,
	})
	return
}

// Detail 获取商城商品分类详情
func (c *cShopCategory) Detail(ctx context.Context, req *v1.ShopCategoryDetailReq) (res *v1.ShopCategoryDetailRes, err error) {
	res = &v1.ShopCategoryDetailRes{}
	res.ShopCategoryDetailOutput, err = service.ShopCategory().Detail(ctx, req.ID)
	return
}

// List 获取商城商品分类列表
func (c *cShopCategory) List(ctx context.Context, req *v1.ShopCategoryListReq) (res *v1.ShopCategoryListRes, err error) {
	res = &v1.ShopCategoryListRes{}
	res.List, res.Total, err = service.ShopCategory().List(ctx, &model.ShopCategoryListInput{
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
// Export 导出商城商品分类
func (c *cShopCategory) Export(ctx context.Context, req *v1.ShopCategoryExportReq) (res *v1.ShopCategoryExportRes, err error) {
	list, err := service.ShopCategory().Export(ctx, &model.ShopCategoryListInput{
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
	r.Response.Header().Set("Content-Disposition", `attachment; filename="shop_category.csv"`)
	r.Response.Write("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(r.Response.Writer)
	// 表头（与导入模板列对齐，末尾追加只读列）
	_ = w.Write([]string{"上级分类", "分类名称", "分类图标", "排序", "状态", "创建时间"})
	// 数据行
	for _, item := range list {
		_ = w.Write([]string{
			csvSafeShopCategory(item.ShopCategoryName),
			csvSafeShopCategory(item.Name),
			csvSafeShopCategory(item.Icon),
			fmt.Sprintf("%v", item.Sort),
			fmt.Sprintf("%v", item.Status),
			func() string { if item.CreatedAt != nil { return item.CreatedAt.String() }; return "" }(),
		})
	}
	w.Flush()
	return
}

// Tree 获取商城商品分类树形结构
func (c *cShopCategory) Tree(ctx context.Context, req *v1.ShopCategoryTreeReq) (res *v1.ShopCategoryTreeRes, err error) {
	res = &v1.ShopCategoryTreeRes{}
	res.List, err = service.ShopCategory().Tree(ctx, &model.ShopCategoryTreeInput{
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

