package product

import (
	"context"
	"encoding/csv"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"

	v1 "gbaseadmin/app/demo/api/demo/v1"
	"gbaseadmin/app/demo/internal/model"
	"gbaseadmin/app/demo/internal/service"
)

var Product = cProduct{}

type cProduct struct{}

// Create 创建体验商品
func (c *cProduct) Create(ctx context.Context, req *v1.ProductCreateReq) (res *v1.ProductCreateRes, err error) {
	err = service.Product().Create(ctx, &model.ProductCreateInput{
		CategoryID: req.CategoryID,
		SkuNo: req.SkuNo,
		Name: req.Name,
		Cover: req.Cover,
		ManualFile: req.ManualFile,
		DetailContent: req.DetailContent,
		SpecJSON: req.SpecJSON,
		WebsiteURL: req.WebsiteURL,
		Type: req.Type,
		IsRecommend: req.IsRecommend,
		SalePrice: req.SalePrice,
		StockNum: req.StockNum,
		WeightNum: req.WeightNum,
		Sort: req.Sort,
		Icon: req.Icon,
		Status: req.Status,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
	})
	return
}

// Update 更新体验商品
func (c *cProduct) Update(ctx context.Context, req *v1.ProductUpdateReq) (res *v1.ProductUpdateRes, err error) {
	err = service.Product().Update(ctx, &model.ProductUpdateInput{
		ID: req.ID,
		CategoryID: req.CategoryID,
		SkuNo: req.SkuNo,
		Name: req.Name,
		Cover: req.Cover,
		ManualFile: req.ManualFile,
		DetailContent: req.DetailContent,
		SpecJSON: req.SpecJSON,
		WebsiteURL: req.WebsiteURL,
		Type: req.Type,
		IsRecommend: req.IsRecommend,
		SalePrice: req.SalePrice,
		StockNum: req.StockNum,
		WeightNum: req.WeightNum,
		Sort: req.Sort,
		Icon: req.Icon,
		Status: req.Status,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
	})
	return
}

// Delete 删除体验商品
func (c *cProduct) Delete(ctx context.Context, req *v1.ProductDeleteReq) (res *v1.ProductDeleteRes, err error) {
	err = service.Product().Delete(ctx, req.ID)
	return
}

// BatchDelete 批量删除体验商品
func (c *cProduct) BatchDelete(ctx context.Context, req *v1.ProductBatchDeleteReq) (res *v1.ProductBatchDeleteRes, err error) {
	err = service.Product().BatchDelete(ctx, req.IDs)
	return
}

// BatchUpdate 批量编辑体验商品
func (c *cProduct) BatchUpdate(ctx context.Context, req *v1.ProductBatchUpdateReq) (res *v1.ProductBatchUpdateRes, err error) {
	err = service.Product().BatchUpdate(ctx, &model.ProductBatchUpdateInput{
		IDs: req.IDs,
		Type: req.Type,
		IsRecommend: req.IsRecommend,
		Status: req.Status,
	})
	return
}

// Detail 获取体验商品详情
func (c *cProduct) Detail(ctx context.Context, req *v1.ProductDetailReq) (res *v1.ProductDetailRes, err error) {
	res = &v1.ProductDetailRes{}
	res.ProductDetailOutput, err = service.Product().Detail(ctx, req.ID)
	return
}

// List 获取体验商品列表
func (c *cProduct) List(ctx context.Context, req *v1.ProductListReq) (res *v1.ProductListRes, err error) {
	res = &v1.ProductListRes{}
	res.List, res.Total, err = service.Product().List(ctx, &model.ProductListInput{
		PageNum:   req.PageNum,
		PageSize:  req.PageSize,
		OrderBy:   req.OrderBy,
		OrderDir:  req.OrderDir,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		SkuNo: req.SkuNo,
		Name: req.Name,
		CategoryID: req.CategoryID,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
		Type: req.Type,
		IsRecommend: req.IsRecommend,
		Status: req.Status,
	})
	return
}
// Export 导出体验商品
func (c *cProduct) Export(ctx context.Context, req *v1.ProductExportReq) (res *v1.ProductExportRes, err error) {
	list, err := service.Product().Export(ctx, &model.ProductListInput{
		OrderBy:   req.OrderBy,
		OrderDir:  req.OrderDir,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		SkuNo: req.SkuNo,
		Name: req.Name,
		CategoryID: req.CategoryID,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
		Type: req.Type,
		IsRecommend: req.IsRecommend,
		Status: req.Status,
	})
	if err != nil {
		return
	}
	// CSV 导出（使用 csv.Writer 防止注入和格式问题）
	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Response.Header().Set("Content-Disposition", `attachment; filename="product.csv"`)
	r.Response.Write("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(r.Response.Writer)
	// 表头
	_ = w.Write([]string{"商品分类", "SKU编号", "商品名称", "封面", "说明书文件", "详情内容", "规格JSON", "官网URL", "类型", "是否推荐", "销售价", "库存数量", "重量", "排序", "图标", "状态", "租户", "商户", "创建时间"})
	// 数据行
	for _, item := range list {
		_ = w.Write([]string{
			item.CategoryName,
			item.SkuNo,
			item.Name,
			item.Cover,
			item.ManualFile,
			item.DetailContent,
			item.SpecJSON,
			item.WebsiteURL,
			fmt.Sprintf("%v", item.Type),
			fmt.Sprintf("%v", item.IsRecommend),
			fmt.Sprintf("%v", item.SalePrice),
			fmt.Sprintf("%v", item.StockNum),
			fmt.Sprintf("%v", item.WeightNum),
			fmt.Sprintf("%v", item.Sort),
			item.Icon,
			fmt.Sprintf("%v", item.Status),
			item.TenantName,
			item.MerchantName,
			func() string { if item.CreatedAt != nil { return item.CreatedAt.String() }; return "" }(),
		})
	}
	w.Flush()
	return
}

// Import 导入体验商品
func (c *cProduct) Import(ctx context.Context, req *v1.ProductImportReq) (res *v1.ProductImportRes, err error) {
	r := g.RequestFromCtx(ctx)
	file := r.GetUploadFile("file")
	if file == nil {
		return nil, fmt.Errorf("请上传文件")
	}
	success, fail, err := service.Product().Import(ctx, file)
	if err != nil {
		return nil, err
	}
	res = &v1.ProductImportRes{Success: success, Fail: fail}
	return
}

// ImportTemplate 下载体验商品导入模板
func (c *cProduct) ImportTemplate(ctx context.Context, req *v1.ProductImportTemplateReq) (res *v1.ProductImportTemplateRes, err error) {
	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Response.Header().Set("Content-Disposition", `attachment; filename="product_template.csv"`)
	r.Response.Write("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(r.Response.Writer)
	_ = w.Write([]string{"商品分类", "SKU编号", "商品名称", "封面", "说明书文件", "详情内容", "规格JSON", "官网URL", "类型", "是否推荐", "销售价", "库存数量", "重量", "排序", "图标", "状态", "租户", "商户"})
	w.Flush()
	return
}
